package fly

import (
	"context"
	"fly-go/database"
	log "fly-go/logger"
	"sync"
	"time"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks      map[string]*Task
	mu         sync.RWMutex
	db         *database.MongoDB
	logger     *log.ILogger
	collection string
	running    bool
	stopCh     chan struct{}
}

// NewTaskManager 创建任务管理器
func NewTaskManager(db *database.MongoDB, logger *log.ILogger) *TaskManager {
	return &TaskManager{
		tasks:      make(map[string]*Task),
		db:         db,
		logger:     logger,
		collection: "tasks",
		stopCh:     make(chan struct{}),
	}
}

// RegisterTask 注册任务
func (tm *TaskManager) RegisterTask(task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks[task.ID] = task
	tm.logger.Info("Task registered",
		log.String("id", task.ID),
		log.String("name", task.Name),
		log.String("executor", task.ExecutorName),
	)
}

// UnregisterTask 注销任务
func (tm *TaskManager) UnregisterTask(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if task, ok := tm.tasks[taskID]; ok {
		task.Trigger.Enabled = false
		delete(tm.tasks, taskID)
		tm.logger.Info("Task unregistered", log.String("id", taskID))
	}
}

// GetTask 获取任务
func (tm *TaskManager) GetTask(taskID string) *Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.tasks[taskID]
}

// ListTasks 列出所有任务
func (tm *TaskManager) ListTasks() []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// LoadFromDB 从数据库加载任务
func (tm *TaskManager) LoadFromDB() error {
	tasks, err := LoadTask(tm.db, tm.logger)
	if err != nil {
		return err
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tasks {
		if task.executor != nil {
			tm.tasks[task.ID] = task
			tm.logger.Info("Task loaded from DB",
				log.String("id", task.ID),
				log.String("name", task.Name),
			)
		} else {
			tm.logger.Warn("No executor found for task",
				log.String("id", task.ID),
				log.String("name", task.Name),
			)
		}
	}

	tm.logger.Info("Loaded tasks", log.Int("count", len(tm.tasks)))
	return nil
}

// InitDefaultTask 初始化默认任务
func (tm *TaskManager) InitDefaultTask() error {
	task := CreateDefaultTask()
	task.db = tm.db
	task.logger = tm.logger
	task.Trigger.Refresh()
	task.Save()

	tm.RegisterTask(task)
	return nil
}

// SaveAll 保存所有任务到数据库
func (tm *TaskManager) SaveAll() {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	for _, task := range tm.tasks {
		if err := task.Save(); err != nil {
			tm.logger.Error("Failed to save task",
				log.String("id", task.ID),
				log.Error(err.Error()),
			)
		}
	}
}

// Start 启动任务调度
func (tm *TaskManager) Start() {
	tm.running = true
	tm.logger.Info("Task manager started")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-tm.stopCh:
			tm.logger.Info("Task manager stopping")
			tm.SaveAll()
			tm.running = false
			return
		case <-ticker.C:
			tm.runPendingTasks()
		}
	}
}

// Stop 停止任务调度
func (tm *TaskManager) Stop() {
	if tm.running {
		close(tm.stopCh)
		tm.logger.Info("Task manager stopped")
	}
}

// runPendingTasks 运行待执行的任务
func (tm *TaskManager) runPendingTasks() {
	tm.mu.RLock()
	tasks := tm.tasks
	tm.mu.RUnlock()

	for _, task := range tasks {
		if !task.Trigger.CheckAndUpdate() {
			continue
		}

		tm.logger.Info("Running task",
			log.String("id", task.ID),
			log.String("name", task.Name),
		)

		go tm.executeTask(task)
	}
}

// executeTask 执行单个任务
func (tm *TaskManager) executeTask(task *Task) {
	task.UpdateStatus(StatusRunning, "Task is running")

	result := NewTaskResult()
	result.StartTime = time.Now()

	if task.executor == nil {
		result.Status = StatusError
		result.Message = "No executor registered"
		task.UpdateStatus(StatusError, result.Message)
		return
	}

	res, err := task.executor.Execute(context.Background(), task)
	result = res
	if err != nil {
		result.Status = StatusError
		result.Message = err.Error()
		tm.logger.Error("Task execution failed",
			log.String("id", task.ID),
			log.Error(err.Error()),
		)
	} else {
		result.Status = StatusSuccess
		result.Message = "Task completed successfully"
		tm.logger.Info("Task executed successfully",
			log.String("id", task.ID),
		)
	}

	result.EndTime = time.Now()
	result.SpendTime = result.EndTime.Sub(result.StartTime).Seconds()

	task.AddResult(result)
	task.LastRunTime = result.StartTime
	task.UpdateStatus(result.Status, result.Message)
}
