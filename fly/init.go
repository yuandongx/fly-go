// Package fly is the main package for the task scheduling application.
//
//	It initializes the application, loads tasks from the database, and manages their execution.
package fly

import (
	"fmt"

	"fly-go/config"
	"fly-go/database"
	log "fly-go/logger"
	"go.uber.org/zap"
)

type ExampleTask struct{}

func (t *ExampleTask) Run() ([]BM, error) {
	// 模拟任务执行逻辑
	fmt.Println("Running ExampleTask...")
	return nil, nil
}
func (t *ExampleTask) Stop() error {
	// 模拟任务停止逻辑
	fmt.Println("Stopping ExampleTask...")
	return nil
}

func taskInfMap() map[string]TaskInterface {
	return map[string]TaskInterface{
		"default_task": &ExampleTask{},
	}
}

func Init() (*TaskManager, error) {
	logger := log.DefaultLogger()
	logger.Info("Initializing application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	mongoDB, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer mongoDB.Close()

	// 注册任务接口实现
	tm := NewTaskManager(mongoDB, logger)
	err = tm.DumpDefaultTask()
	if err != nil {
		logger.Error("Failed to dump default tasks", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to dump default tasks: %w", err)
	}
	err = tm.LoadTask()
	if err != nil {
		logger.Error("Failed to load tasks", zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	// e := tm.DumpTask()
	// if e != nil {
	// 	logger.Error("Failed to dump tasks", log.Zap("Error", e.Error()))
	// 	return nil, fmt.Errorf("failed to dump tasks: %w", e)
	// }

	logger.Info("Application initialized successfully")
	return tm, nil
}
