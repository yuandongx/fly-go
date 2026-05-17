package executor

import (
	"context"
	"fly-go/database"
	"time"
)

const TASK_COLLECTION = "task"

type TaskQueue struct {
	Tasks        map[int]*Task     // 任务队列，key为加入的顺序
	Runners      map[string]Runner // 已注册的runner
	DB           *database.MongoDB
	runningCount int
	waitCount    int
	finishCount  int
	errorCount   int
}

// New 创建一个新的 TaskQueue
func New(db *database.MongoDB) *TaskQueue {
	return &TaskQueue{
		Tasks:        make(map[int]*Task),
		Runners:      make(map[string]Runner),
		DB:           db,
		runningCount: 0,
		waitCount:    0,
		finishCount:  0,
		errorCount:   0,
	}
}

// Init 从数据库加载任务队列
// 1. 从 MongoDB 读取所有任务数据
// 2. 根据 RunnerName 关联 Runner
// 3. 将任务加入队列
func (tq *TaskQueue) Init() {
	// 从数据库读取任务数据
	tasks, err := tq.loadFromDB()
	if err != nil {
		return
	}

	// 将任务加入队列
	for i, task := range tasks {
		// 根据 RunnerName 关联 Runner
		if r, ok := tq.Runners[task.RunnerName]; ok {
			task.Runner = r
		}
		tq.Tasks[i] = task
		tq.waitCount++
	}
}

// loadFromDB 从 MongoDB 加载所有任务
func (tq *TaskQueue) loadFromDB() ([]*Task, error) {
	rows := database.Rows{}
	collection, err := tq.DB.Collection(TASK_COLLECTION)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	cursor, err := collection.Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	// 转换为 Task 列表
	tasks := make([]*Task, 0, len(rows))
	for _, row := range rows {
		task := tq.rowToTask(row)
		if task != nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// rowToTask 将数据库行转换为 Task 结构体
func (tq *TaskQueue) rowToTask(row database.Row) *Task {
	task := &Task{
		DB: tq.DB,
	}

	// 提取字符串字段
	if name, ok := row["name"].(string); ok {
		task.Name = name
	}
	if runnerName, ok := row["runner_name"].(string); ok {
		task.RunnerName = runnerName
	}
	if lastStatus, ok := row["last_status"].(string); ok {
		task.LastStatus = lastStatus
	}
	if lastMessage, ok := row["last_message"].(string); ok {
		task.LastMessage = lastMessage
	}

	// 提取 Params
	if params, ok := row["params"].(map[string]any); ok {
		task.Params.Init(params)
	}

	// 提取时间字段
	if lastRunTime, ok := row["last_run_time"].(string); ok {
		task.LastRunTime = parseTime(lastRunTime)
	}
	if lastEndTime, ok := row["last_end_time"].(string); ok {
		task.LastEndTime = parseTime(lastEndTime)
	}

	return task
}

// parseTime 解析时间字符串
func parseTime(s string) (t time.Time) {
	t, _ = time.Parse(time.RFC3339, s)
	return
}

// 开始任务
// 1、开始执行每一个
// WaitCount 返回等待执行的任务数
func (tq *TaskQueue) WaitCount() int {
	return tq.waitCount
}

// RunningCount 返回正在运行的任务数
func (tq *TaskQueue) RunningCount() int {
	return tq.runningCount
}

// FinishCount 返回已完成的任务数
func (tq *TaskQueue) FinishCount() int {
	return tq.finishCount
}

// ErrorCount 返回执行失败的任务数
func (tq *TaskQueue) ErrorCount() int {
	return tq.errorCount
}

// Start 启动任务队列，遍历所有任务并执行符合条件的任务
func (tq *TaskQueue) Execute() {
	for i, task := range tq.Tasks {
		// 检查任务是否可以执行
		if !task.CanRun() {
			continue
		}

		// 检查 Runner 是否存在
		if task.Runner == nil {
			tq.waitCount--
			continue
		}

		// 使用 goroutine 并发执行任务
		go func(idx int, t *Task) {
			tq.runningCount++
			tq.waitCount--

			// 执行任务
			result, err := t.Run()
			if err != nil {
				result.Error = err
				result.Status = StatusError
			}

			// 保存结果并更新状态
			t.Save(result)
			tq.finishCount++
			if result.Status == StatusError {
				tq.errorCount++
			}
			tq.runningCount--

			if t.Logger != nil {
				t.Logger.Infof("任务 [%s] 执行完成，状态: %s", t.Name, result.Status)
			}
		}(i, task)
	}
}

// 启动TaskQueue后，就可以循环执行Execute
func (tq *TaskQueue) Start() {
	for {
		tq.Execute()
		time.Sleep(time.Second)
	}
}

func (tq *TaskQueue) Register(runner Runner) {
	tq.Runners[runner.Name()] = runner
}
