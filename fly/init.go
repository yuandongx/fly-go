package fly

import (
	"fmt"

	"fly-go/config"
	"fly-go/database"
	log "fly-go/logger"
	"go.uber.org/zap"
)

// Init 初始化任务系统
func Init() (*TaskManager, error) {
	logger := log.DefaultLogger()
	logger.Info("Initializing fly task system...")

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 连接数据库
	db, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 创建任务管理器
	tm := NewTaskManager(db, logger)

	// 初始化默认任务
	if err := tm.InitDefaultTask(); err != nil {
		logger.Error("Failed to init default task", zap.String("error", err.Error()))
	}

	// 从数据库加载任务
	if err := tm.LoadFromDB(); err != nil {
		logger.Error("Failed to load tasks", zap.String("error", err.Error()))
	}

	logger.Info("Fly task system initialized",
		log.Int("executors", len(ListExecutors())),
	)
	return tm, nil
}
