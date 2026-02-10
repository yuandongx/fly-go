package fly

import (
	"context"

	"fly-go/fly/runner"
	"fly-go/internal/config"
	"fly-go/internal/database"
	log "fly-go/logger"
)

// Global task manager instance
var tm *TaskManager = NewTaskManager()
var cntxt context.Context = context.Background()

func Init() {
	logger := log.DefaultLogger()

	logger.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config")
	}

	mongoDB, err := database.NewMongoDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", log.Zap("Error", err.Error()))
	}
	defer mongoDB.Close()

	// 注册任务
	registerTask(mongoDB, logger)

}

func registerTask(db *database.MongoDB, logger *log.ILogger) {
	// sina 更新信息
	stock := runner.NewStock()
	stockRunner := NewRunner("001",
		"sina-get_all_stock_info",
		"新浪更新股票信息",
		"stock", stock)
	stockRunner.Trigger = *NewTrigger(Interval)
	stockRunner.Trigger.SetPeriod(10)
	stockRunner.Trigger.SetWeekDays([]int{1, 2, 3, 4, 5})
	stockRunner.Trigger.SetRangeTime("09:25", "11:30")
	stockRunner.Trigger.SetRangeTime("13:00", "13:00")
	stockTask := NewTask(db, stockRunner, logger)
	tm.AddTask(*stockTask)
}
