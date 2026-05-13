package executor

import (
	"fly-go/database"
	"fly-go/fly"
)

// Register 注册所有执行器
func Register(db *database.MongoDB) {
	// 注册示例执行器
	fly.RegisterExecutor(&ExampleExecutor{})

	// 注册股票执行器
	fly.RegisterExecutor(NewStockExecutor(db))
}
