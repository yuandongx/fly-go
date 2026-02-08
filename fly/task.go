package fly

import (
	"context"
	"fly-go/internal/database"
	log "fly-go/logger"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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



type Trigger struct {
	Type        string `json:"type" bson:"type,omitempty"`
	Interval    int64  `json:"interval" bson:"interval,omitempty"`
	Enabled     bool   `json:"enabled" bson:"enabled,omitempty"`
	StartAtDate int64  `json:"start_at" bson:"start_at,omitempty"`
	EndAtDate   int64  `json:"end_at" bson:"end_at,omitempty"`
	StartTime   int64  `json:"start_time" bson:"start_time,omitempty"`
	EndTime     int64  `json:"end_time" bson:"end_time,omitempty"`
}

type Runner struct {
	ID          string  `json:"id" bson:"id,omitempty"`
	Name        string  `json:"name" bson:"name,omitempty"`
	Description string  `json:"description" bson:"description,omitempty"`
	LastRuntime int64   `json:"last_runtime" bson:"last_runtime,omitempty"`
	NextRuntime int64   `json:"next_runtime" bson:"next_runtime,omitempty"`
	Status      string  `json:"status" bson:"status,omitempty"`
	Msg         string  `json:"msg" bson:"msg,omitempty"`
	Trigger     Trigger `json:"trigger" bson:"trigger,omitempty"`
	Colllection string  `json:"collection" bson:"collection,omitempty"`
	Task        TaskInterface
}

type Task struct {
	DB             *database.MongoDB
	Runner         *Runner
	Logger         *log.ILogger
	TaskCollection string
}

type TaskManager struct {
	TM    map[string]Task
	Count int
	Names []string
}

func NewRunner(id, name, des, collection string, task TaskInterface) *Runner {
	return &Runner{
		ID:          id,
		Name:        name,
		Description: des,
		Colllection: collection,
		Task:        task,
	}

}

func NewTask(db *database.MongoDB, runner *Runner, logger *log.ILogger) *Task {
	return &Task{DB: db, Runner: runner, Logger: logger, TaskCollection: "task-record"}
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Count: 0,
		Names: make([]string, 0),
		TM:    make(map[string]Task),
	}
}

// CanRun returns true if the trigger is enabled and the current time is greater than or equal to the next runtime.
func (t *Task) CanRun() bool {
	// if t.Runner.Trigger.Enabled {
	// 	currentTime := time.Now().Unix()
	// 	diff := int64(math.Abs(float64(currentTime - t.Runner.NextRuntime)))
	// 	if diff*2 < t.Runner.Trigger.Interval {
	// 		// Update last runtime and next runtime before running the task
	// 		t.Runner.LastRuntime = currentTime
	// 		t.Runner.NextRuntime = currentTime + t.Runner.Trigger.Interval
	// 		return true
	// 	}
	// }
	// return false

	return true
}

// Run task
func (t *Task) Run() error {
	start := time.Now()
	res, err := t.Runner.Task.Run()
	var models []mongo.WriteModel
	// 结果检查
	if len(res) == 0 {
		t.Runner.Status = StatusSuccess
		t.Runner.Msg = "No result was got"
		t.Update()
		return err
	}

	// 集合
	coll := t.Runner.Colllection
	if coll == "" {
		coll = t.TaskCollection
	}

	// 保存
	for _, m := range res {
		um := mongo.NewUpdateOneModel().SetFilter(bson.D{{Key: "code", Value: m["code"]}}).SetUpsert(true).SetUpdate(bson.D{{Key: "$set", Value: m}})
		models = append(models, um)
	}
	opts := options.BulkWrite().SetOrdered(false)
	ctx := context.TODO()
	collect := t.DB.Collection(coll)
	if collect == nil {
		t.Logger.Info("MongoDB 连接异常！")
		return nil
	}

	// 批量写入
	msg := ""
	if _, err := collect.BulkWrite(ctx, models, opts); err != nil {
		t.Runner.Msg = fmt.Sprintf("%s 更新失败了 %v", t.Runner.Name, err)
		return err
	} else {
		n := len(models)
		msg = fmt.Sprintf("更新了 %d 条记录", n)
		t.Logger.Info(msg)

	}
	end := time.Now()
	spend := end.Sub(start).Seconds()
	msg += fmt.Sprintf(" 用时 %10.2fs", spend)
	t.Runner.Msg = msg
	return nil
}

// Update updates the task's last runtime, next runtime, status, and message in the database
func (t *Task) Update() {
	collection := t.DB.Collection(t.TaskCollection)
	filter := map[string]interface{}{"id": t.Runner.ID}

	// Update last runtime and next runtime
	t.Runner.LastRuntime = time.Now().Unix()
	t.Runner.NextRuntime = t.Runner.LastRuntime + t.Runner.Trigger.Interval
	update := bson.D{{Key: "$set", Value: t.Runner}}
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := collection.UpdateOne(cntxt, filter, update, opts)
	if err != nil {
		println("Error updating task:", err.Error())
	}
}

// AddTask adds a new task to the TaskManager and assigns it a unique ID based on the last task's ID in the TaskManager slice
func (tm *TaskManager) AddTask(task Task) {
	tm.Count += 1
	key := fmt.Sprintf("id.%s.%s", task.Runner.ID, task.Runner.Name)
	tm.TM[key] = task
	tm.Names = append(tm.Names, task.Runner.Name)
}

// RemoveTask sets the status of the task to stopped based on the task ID
// and disables the trigger, but does not remove the task from the TaskManager slice
func (tm *TaskManager) RemoveTask(taskID string) {
	delete(tm.TM, taskID)
	tm.Count = tm.Count - 1
}
