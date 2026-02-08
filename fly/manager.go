package fly

import (
	"context"
	"fly-go/database"
	log "fly-go/logger"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type TaskManager struct {
	TM         map[string]Runner
	Count      int
	Names      []string
	DB         *database.MongoDB
	Logger     *log.ILogger
	Collection string
}

func NewTaskManager(db *database.MongoDB, logger *log.ILogger) *TaskManager {
	return &TaskManager{
		Count:      0,
		DB:         db,
		Logger:     logger,
		TM:         make(map[string]Runner),
		Collection: "tasks",
	}
}

// AddTask adds a new task to the TaskManager and assigns it a unique ID based on the last task's ID in the TaskManager slice
func (tm *TaskManager) AddTask(task Runner) {
	tm.Count += 1
	key := fmt.Sprintf("id.%s.%s", task.ID, task.Name)
	tm.TM[key] = task
	tm.Names = append(tm.Names, task.Name)
}

// RemoveTask sets the status of the task to stopped based on the task ID
// and disables the trigger, but does not remove the task from the TaskManager slice
func (tm *TaskManager) RemoveTask(taskID string) {
	delete(tm.TM, taskID)
	tm.Count = tm.Count - 1
}

// LoadTask 从数据库加载任务信息并创建对应的 Runner 实例添加到 TaskManager 中
func (tm *TaskManager) LoadTask() error {

	taskInf := taskInfMap()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 从数据库查询所有任务记录
	collection, err := tm.DB.Collection(tm.Collection)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []Runner
	if err := cursor.All(ctx, &tasks); err != nil {
		return fmt.Errorf("failed to decode tasks: %w", err)
	}

	// 为每个任务创建 Runner 并添加到 TaskManager
	for _, taskRunner := range tasks {
		taskID := taskRunner.ID
		if taskID == "" {
			tm.Logger.Warn("Skipping task with empty ID", zap.String("Name", taskRunner.Name))
			continue
		}

		tif, ok := taskInf[taskRunner.Name]
		if !ok {
			tm.Logger.Warn("No task implementation found for task", zap.String("Name", taskRunner.Name))
			continue
		} else {
			// 添加到任务管理器
			taskRunner.Task = tif
			taskRunner.DB = tm.DB
			taskRunner.Logger = tm.Logger
			tm.AddTask(taskRunner)
			tm.Logger.Info("Task loaded", zap.String("id", taskID), zap.String("name", taskRunner.Name))
		}
	}

	tm.Logger.Info(fmt.Sprintf("Loaded %d tasks from database", tm.Count))
	return nil
}

func (tm *TaskManager) DumpTask() error {

	if tm.Count == 0 {
		tm.Logger.Info("No tasks to dump")
		return nil
	}

	// TODO: 实现将 TaskManager 中的任务信息保存到数据库的逻辑
	// 需要遍历 tm.TM 并更新每个任务的状态到数据库
	for _, task := range tm.TM {

		task.Update()

	}
	tm.Logger.Info(fmt.Sprintf("Dumped %d tasks to database", tm.Count))
	return nil
}

// DumpDefaultTask 将 TaskManager 中的默认任务信息保存到数据库
func (tm *TaskManager) DumpDefaultTask() error {
	runner := Runner{
		ID:            "000",
		Name:          "default_task",
		Status:        StatusUnknown,
		Msg:           "This is a default task",
		Collection:    tm.Collection,
		InterfaceName: "default_task",
		Trigger: Trigger{
			Period:      1,
			StartTime:   "00:00",
			EndTime:     "23:59",
			Weekdays:    []time.Weekday{0, 1, 2, 3, 4, 5, 6},
			SkipDays:    []string{"2026-01-01", "2026-02-11", "2026-02-12", "2026-02-13", "2026-02-14", "2026-02-15", "2026-02-16", "2026-02-17", "2026-02-18", "2026-02-19", "2026-02-20", "2026-02-21"},
			Type:        Interval,
			Enabled:     true,
			StartAtDate: "2026-01-01",
			EndAtDate:   "2099-01-01",
			LastRunTime: time.Time{},
			RangeTime:   [][]string{{"00:00", "23:59"}},
		},
		Task:   nil,
		DB:     tm.DB,
		Logger: tm.Logger,
	}

	runner.Update()
	return nil
}

func (tm *TaskManager) RunAllTask() {
	for _, r := range tm.TM {
		go func(runner Runner) {

			// Check if the task can run before executing
			tm.Logger.Info(fmt.Sprintf("Checking task: %s", runner.Name))
			ok := runner.TimeIsUp()
			if !ok {
				tm.Logger.Info(fmt.Sprintf("Task %s is not ready to run", runner.Name))
				return
			}

			// Run the task and update its status
			runner.Status = StatusRunning
			runner.Msg = "Task is running"
			runner.Update()
			result := NewTaskResult()
			result.StartTime = time.Now()

			// Log the task execution
			tm.Logger.Info(fmt.Sprintf("Running task: %s", runner.Name))
			if err := runner.Run(); err != nil {
				runner.Status = StatusError
				runner.Msg = err.Error()
				tm.Logger.Error(fmt.Sprintf("Error running task %s: %s", runner.Name, err.Error()))
				result.Msg = err.Error()
				result.Status = StatusError
			} else {
				runner.Status = StatusSuccess
				tm.Logger.Info(fmt.Sprintf("Task %s completed successfully", runner.Name))
				result.Status = StatusSuccess
				result.Msg = "Task completed successfully"
			}

			// Update the task's last runtime and next runtime after execution
			// Update the task in the database
			result.EndTime = time.Now()
			result.SpendTime = result.EndTime.Sub(result.StartTime).Seconds()
			runner.Results = append(runner.Results, result)
			runner.Update()
		}(r)
	}
}
