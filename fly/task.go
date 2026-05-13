package fly

import (
	"context"
	"fly-go/database"
	logger "fly-go/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task 任务定义
type Task struct {
	ID             string   `json:"id" bson:"id" form:"id"`
	Name           string   `json:"name" bson:"name" form:"name" validate:"required"`
	Description    string   `json:"description" bson:"description" form:"description"`
	ExecutorName   string   `json:"executor_name" bson:"executor_name" form:"executor_name" validate:"required"`
	Collection     string   `json:"collection" bson:"collection" form:"collection"`
	OutputColl     string   `json:"output_collection" bson:"output_collection" form:"output_collection"`
	Trigger        Trigger  `json:"trigger" bson:"trigger" form:"trigger"`
	Params         bson.M   `json:"params" bson:"params" form:"params"`
	Status         string   `json:"status" bson:"status"`
	LastMessage    string   `json:"last_message" bson:"last_message"`
	LastRunTime    time.Time `json:"last_run_time" bson:"last_run_time"`
	NextRunTime    time.Time `json:"next_run_time" bson:"next_run_time"`
	TotalRunCount  int64    `json:"total_run_count" bson:"total_run_count"`
	SuccessCount   int64    `json:"success_count" bson:"success_count"`
	FailureCount   int64    `json:"failure_count" bson:"failure_count"`
	RecentResults  []TaskResult `json:"recent_results" bson:"recent_results"`

	// 运行时依赖
	executor TaskExecutor
	db       *database.MongoDB
	logger   *logger.ILogger
}

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Name() string
	Execute(ctx context.Context, task *Task) (TaskResult, error)
}

// executorRegistry 任务执行器注册表
var executorRegistry = make(map[string]TaskExecutor)

// RegisterExecutor 注册任务执行器
func RegisterExecutor(executor TaskExecutor) {
	executorRegistry[executor.Name()] = executor
}

// GetExecutor 获取任务执行器
func GetExecutor(name string) TaskExecutor {
	return executorRegistry[name]
}

// ListExecutors 列出所有已注册的执行器
func ListExecutors() []string {
	names := make([]string, 0, len(executorRegistry))
	for name := range executorRegistry {
		names = append(names, name)
	}
	return names
}

// logInfo 日志记录函数
func logInfo(msg string, fields ...interface{}) {
	logger.DefaultLogger().Infof(msg, fields...)
}

// NewTask 创建任务
func NewTask(db *database.MongoDB, l *logger.ILogger) *Task {
	return &Task{
		Status:        StatusIdle,
		Params:        bson.M{},
		RecentResults: make([]TaskResult, 0, 5),
		db:           db,
		logger:       l,
	}
}

// LoadTask 从数据库加载任务
func LoadTask(db *database.MongoDB, l *logger.ILogger) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll, err := db.Collection("tasks")
	if err != nil {
		return nil, err
	}

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	// 初始化任务依赖
	for _, task := range tasks {
		task.db = db
		task.logger = l
		task.Trigger.Refresh()

		// 查找执行器
		if executor, ok := executorRegistry[task.ExecutorName]; ok {
			task.executor = executor
		}
	}

	return tasks, nil
}

// Save 保存任务到数据库
func (t *Task) Save() error {
	if t.db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll, err := t.db.Collection("tasks")
	if err != nil {
		return err
	}

	filter := bson.M{"id": t.ID}
	update := bson.M{"$set": t}
	opts := options.Update().SetUpsert(true)
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateStatus 更新任务状态
func (t *Task) UpdateStatus(status, message string) {
	t.Status = status
	t.LastMessage = message
	t.Save()
}

// AddResult 添加执行结果
func (t *Task) AddResult(result TaskResult) {
	t.RecentResults = append(t.RecentResults, result)
	if len(t.RecentResults) > 5 {
		t.RecentResults = t.RecentResults[len(t.RecentResults)-5:]
	}
	t.TotalRunCount++
	if result.Status == StatusSuccess {
		t.SuccessCount++
	} else {
		t.FailureCount++
	}
}

// CreateDefaultTask 创建默认任务示例
func CreateDefaultTask() *Task {
	return &Task{
		ID:           "default",
		Name:         "default_task",
		Description:  "默认任务示例",
		ExecutorName: "example",
		Collection:   "tasks",
		OutputColl:   "output",
		Status:       StatusIdle,
		Trigger: Trigger{
			Type:        TriggerInterval,
			Enabled:     true,
			Period:      60,
			StartAtDate: time.Now().Format("2006-01-02"),
			EndAtDate:   "2099-12-31",
			StartTime:   "00:00",
			EndTime:     "23:59",
			Weekdays:    []int{0, 1, 2, 3, 4, 5, 6},
			SkipDays:    []string{},
		},
		Params: bson.M{},
	}
}
