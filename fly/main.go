package fly

import (
	"fmt"

	"fly-go/config"
	"fly-go/database"
	log "fly-go/logger"

	"go.uber.org/zap"
)

// ListExecutors 返回所有已注册的 Executor
func ListExecutors() []string {
	runners := GetRunners()
	names := make([]string, 0, len(runners))
	for name := range runners {
		names = append(names, name)
	}
	return names
}

// Init 初始化任务系统
// 1. 加载配置文件
// 2. 连接数据库
// 3. 创建任务管理器
// 4. 初始化默认任务
// 5. 加载数据库任务
// 6. 注册内置 Executor
func Init() (*TaskManager, error) {
	logger := log.DefaultLogger()
	logger.Info("Initializing fly task system...")

	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. 连接数据库
	db, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 3. 创建任务管理器
	tm := NewTaskManager(db, logger)

	// 4. 初始化默认任务（失败不影响主流程）
	if err := tm.InitDefaultTask(); err != nil {
		logger.Warn("Failed to init default task", zap.String("error", err.Error()))
	}

	// 5. 注册内置 Executor
	runners := GetRunners()
	tm.SetRunners(runners)
	logger.Info("Fly task system initialized",
		log.Int("runners", len(runners)),
		log.Int("tasks", tm.Queue.WaitCount()),
	)

	// 6. 从数据库加载任务（失败不影响主流程）
	if err := tm.LoadFromDB(); err != nil {
		logger.Warn("Failed to load tasks", zap.String("error", err.Error()))
	}
	return tm, nil
}

func Start() {
	tm, err := Init()
	if err == nil {
		tm.Start()
	}
}
