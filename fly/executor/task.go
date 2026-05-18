package executor

import (
	"context"
	"fly-go/database"
	log "fly-go/logger"
	"time"
)

type Runner interface {
	Name() string
	Run() (TaskResult, error)
	Stop() error
}

type TaskResult struct {
	Status    string        `json:"status"`
	Data      any           `json:"data"`
	Error     error         `json:"error,omitempty"`
	Message   string        `json:"message,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

type Task struct {
	Name        string            `json:"name"`
	Runner      Runner            `json:"-"`
	RunnerName  string            `json:"runner_name"`
	Params      TaskParam         `json:"params"`
	Result      []TaskResult      `json:"result"`
	DB          *database.MongoDB `json:"-"`
	Logger      *log.ILogger      `json:"-"`
	LastEndTime time.Time         `json:"last_end_time"`
	LastRunTime time.Time         `json:"last_run_time"`
	LastStatus  string            `json:"last_status"`
	LastMessage string            `json:"last_message"`
}

func (t *Task) Run() (TaskResult, error) {
	return t.Runner.Run()
}

func (t *Task) Stop() error {
	return t.Runner.Stop()
}
func (t *Task) Save(result TaskResult) {
	t.Result = append(t.Result, result)
	// 只保留最近10次结果
	if len(t.Result) > 10 {
		t.Result = t.Result[len(t.Result)-10:]
	}
	t.LastEndTime = result.EndTime
	t.LastRunTime = time.Now()
	t.LastStatus = result.Status
	t.LastMessage = result.Message

	// 同步更新到数据库
	if t.DB != nil {
		t.updateDB()
	}
}

// updateDB 将任务状态同步到数据库
func (t *Task) updateDB() {
	ctx := context.Background()
	filter := map[string]any{"name": t.Name}
	update := map[string]any{
		"$set": map[string]any{
			"last_end_time": t.LastEndTime.Format(time.RFC3339),
			"last_run_time": t.LastRunTime.Format(time.RFC3339),
			"last_status":   t.LastStatus,
			"last_message":  t.LastMessage,
		},
	}
	t.DB.UpdateOne(ctx, TASK_COLLECTION, filter, update)
}

// CanRun 检查当前时间是否可以执行任务
// 是否可以执行依据有以下几点：
//
//	1、TaskParam.Active() - 当前时间是否在有效日期/时间范围内
//	2、固定间隔任务：LastEndTime 与 Now 的间隔是否大于 interval，
//	   且 LastEndTime 是否在 StartTime ~ EndTime 内
//	3、一次执行任务：检查是否执行过，且当前时间是否在预期执行时间之后的10秒内
func (t *Task) CanRun() bool {
	now := time.Now()

	// 1. 检查参数定义的有效范围
	if !t.Params.Active() {
		return false
	}

	// 2. 固定间隔任务 (interval)
	if t.Params.Type == "interval" && t.Params.Interval > 0 {
		// 已执行过，需要检查间隔
		if !t.LastEndTime.IsZero() {
			elapsed := now.Sub(t.LastEndTime)
			// 间隔是否足够（elapsed 必须大于等于 interval）
			if elapsed < t.Params.Interval {
				return false
			}
			// 上次执行时间是否在有效时间范围内
			if !t.isTimeInRange(t.LastEndTime) {
				return false
			}
		}
		return true
	}

	// 3. 一次执行任务 (once)
	if t.Params.Type == "once" {
		// 已执行过，不再执行
		if !t.LastEndTime.IsZero() {
			return false
		}
		// 检查当前时间是否在预期执行时间之后的10秒内
		if !t.LastRunTime.IsZero() {
			elapsed := now.Sub(t.LastRunTime)
			if elapsed < 0 || elapsed > 10*time.Second {
				return false
			}
		}
		return true
	}

	// 默认：每天定时任务
	return true
}

// isTimeInRange 检查指定时间是否在 StartTime ~ EndTime 范围内
func (t *Task) isTimeInRange(tm time.Time) bool {
	currentSec := tm.Hour()*3600 + tm.Minute()*60 + tm.Second()

	if t.Params.StartTime > 0 && currentSec < t.Params.StartTime {
		return false
	}
	if t.Params.EndTime > 0 && currentSec > t.Params.EndTime {
		return false
	}
	return true
}
