package fly

import (
	"fly-go/database"
	"fly-go/fly/executor"
	log "fly-go/logger"
)

type TaskManager struct {
	DB     *database.MongoDB
	Logger *log.ILogger
	Queue  *executor.TaskQueue
}

// NewTaskManager 创建任务管理器
func NewTaskManager(db *database.MongoDB, logger *log.ILogger) *TaskManager {
	return &TaskManager{
		DB:     db,
		Logger: logger,
		Queue:  executor.New(db),
	}
}

// set taskManager 的 Queue 属性runner
func (tm *TaskManager) SetRunners(runners map[string]executor.Runner) {
	tm.Queue.Runners = runners
}

// InitDefaultTask 初始化默认任务
func (tm *TaskManager) InitDefaultTask() error {
	// TODO: 从配置文件或代码中加载默认任务
	return nil
}

// LoadFromDB 从数据库加载任务
func (tm *TaskManager) LoadFromDB() error {
	tm.Queue.Init()
	return nil
}

// Start 启动任务队列
func (tm *TaskManager) Start() {
	tm.Queue.Start()
}
