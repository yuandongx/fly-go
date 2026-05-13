package executor

import (
	"context"
	"fmt"

	"fly-go/fly"
)

// ExampleExecutor 示例执行器
type ExampleExecutor struct{}

func (e *ExampleExecutor) Name() string {
	return "example"
}

func (e *ExampleExecutor) Execute(ctx context.Context, task *fly.Task) (fly.TaskResult, error) {
	result := fly.NewTaskResult()

	// 模拟任务执行
	fmt.Printf("Executing task: %s\n", task.Name)

	// 这里可以添加具体的业务逻辑
	// 可以从 task.Params 中获取参数
	// 可以使用 task.db 查询数据库获取数据源

	return result, nil
}

func init() {
	fly.RegisterExecutor(&ExampleExecutor{})
}
