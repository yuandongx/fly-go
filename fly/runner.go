package fly

import (
	"context"
	"fly-go/database"
	log "fly-go/logger"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	StatusIdle     = "idle"
	StatusRunning  = "running"
	StatusStopped  = "stopped"
	StatusError    = "error"
	StatusSuccess  = "success"
	StatusUnknown  = "unknown"
	StatusTimeout  = "timeout"
	StatusRetry    = "retry"
	StatusCanceled = "canceled"
	StatusPending  = "pending"
)

type BM = bson.M
type BD = bson.D

type TaskInterface interface {
	Run() (result []BM, err error)
	Stop() error
}

type TaskResult struct {
	StartTime time.Time `json:"start_time" bson:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time" bson:"end_time,omitempty"`
	Status    string    `json:"status" bson:"status,omitempty"`
	Msg       string    `json:"msg" bson:"msg,omitempty"`
	SpendTime float64   `json:"spend_time" bson:"spend_time,omitempty"`
	Result    []BM      `json:"result" bson:"result,omitempty"`
}
type Runner struct {
	ID            string            `json:"id" bson:"id,omitempty" form:"id" validate:"required"`
	Name          string            `json:"name" bson:"name,omitempty" form:"name" validate:"required"`
	Description   string            `json:"description" bson:"description,omitempty" form:"description"`
	Status        string            `json:"status" bson:"status,omitempty"`
	Msg           string            `json:"msg" bson:"msg,omitempty"`
	Trigger       Trigger           `json:"trigger" bson:"trigger,omitempty"`
	Collection    string            `json:"collection" bson:"collection,omitempty"`
	InterfaceName string            `json:"interface_name" bson:"interface_name,omitempty" form:"interface_name" validate:"required"`
	Results       []TaskResult      `json:"results" bson:"results,omitempty"`
	Task          TaskInterface     `json:"-" bson:"-"`
	DB            *database.MongoDB `json:"-" bson:"-"`
	Logger        *log.ILogger      `json:"-" bson:"-"`
}

// NewRunner 创建一个新的 Runner 实例
// id: 运行器唯一标识符
// name: 运行器名称
// des: 运行器描述信息
// collection: 所属集合名称
// task: 要执行的任务接口实现
// 返回: 初始化后的 Runner 指针
func NewRunner(id, name, des, collection string, task TaskInterface) *Runner {
	return &Runner{
		ID:          id,
		Name:        name,
		Description: des,
		Collection:  collection,
		Task:        task,
		Results:     make([]TaskResult, 5), // 初始化结果切片，预设长度为5
		Status:      StatusIdle,
		Msg:         "Task is initialized",
	}

}

func NewTaskResult() TaskResult {
	return TaskResult{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Status:    StatusIdle,
		Result:    nil,
		Msg:       "",
		SpendTime: 0,
	}
}

// TimeIsUp returns true if the trigger is enabled and the current time is greater than or equal to the next runtime.
func (runner *Runner) TimeIsUp() bool {
	return runner.Trigger.TimeIsUp()
}

// Run task
func (runner *Runner) Run() error {
	start := time.Now()
	runner.Logger.Info("任务开始执行",
		zap.String("id", runner.ID),
		zap.String("name", runner.Name))
	res, err := runner.Task.Run()
	var models []mongo.WriteModel
	// 结果检查
	if len(res) == 0 {
		runner.Status = StatusSuccess
		runner.Msg = "No result was got"
		runner.Update()
		return err
	}

	// 集合
	coll := runner.Collection
	if coll == "" {
		coll = runner.Collection
	}

	// 保存
	for _, m := range res {
		um := mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "code", Value: m["code"]}}).SetUpsert(true).SetUpdate(bson.D{{Key: "$set", Value: m}})
		models = append(models, um)
	}
	opts := options.BulkWrite().SetOrdered(false)
	ctx := context.TODO()
	collect, err := runner.DB.Collection(coll)
	if err != nil {
		runner.Logger.Info("MongoDB 连接异常！")
		return nil
	}

	// 批量写入
	msg := ""
	if _, err := collect.BulkWrite(ctx, models, opts); err != nil {
		runner.Msg = fmt.Sprintf("%s 更新失败了 %v", runner.Name, err)
		return err
	} else {
		n := len(models)
		msg = fmt.Sprintf("更新了 %d 条记录", n)
		runner.Logger.Info(msg)

	}
	end := time.Now()
	spend := end.Sub(start).Seconds()
	msg += fmt.Sprintf(" 用时 %10.2fs", spend)
	runner.Msg = msg
	return nil
}

// Update updates the task's last runtime, next runtime, status, and message in the database
func (runner *Runner) Update() {
	collection, err := runner.DB.Collection(runner.Collection)
	if err != nil {
		runner.Logger.Error("Error getting collection", zap.String("Error", err.Error()))
		return
	}
	filter := map[string]interface{}{"id": runner.ID}
	ctx := context.TODO()
	// 限制结果集大小
	if len(runner.Results) > 5 {
		runner.Results = runner.Results[len(runner.Results)-5:]
	}

	// Update last runtime and next runtime
	update := bson.D{{Key: "$set", Value: runner}}
	opts := options.Update()
	opts.SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		runner.Logger.Error("Error updating task", zap.String("Error", err.Error()))
	}
	runner.Logger.Info("Task updated successfully", zap.String("id", runner.ID), zap.String("name", runner.Name))
}
