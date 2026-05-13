package fly

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// 任务状态常量
const (
	StatusIdle     = "idle"
	StatusRunning  = "running"
	StatusStopped  = "stopped"
	StatusError    = "error"
	StatusSuccess  = "success"
	StatusTimeout  = "timeout"
	StatusPending  = "pending"
	StatusCanceled = "canceled"
)

// 触发器类型
type TriggerType int

const (
	TriggerInterval TriggerType = iota // 固定间隔执行
	TriggerOnce                       // 执行一次
	TriggerDaily                      // 每天定时执行
	TriggerWeekly                     // 每周定时执行
	TriggerCron                       // Cron表达式执行
)

func (t TriggerType) String() string {
	switch t {
	case TriggerInterval:
		return "interval"
	case TriggerOnce:
		return "once"
	case TriggerDaily:
		return "daily"
	case TriggerWeekly:
		return "weekly"
	case TriggerCron:
		return "cron"
	default:
		return "unknown"
	}
}

// Trigger 触发器配置
type Trigger struct {
	Type        TriggerType   `json:"type" bson:"type" form:"type"`
	Enabled     bool          `json:"enabled" bson:"enabled" form:"enabled"`
	Period      int64         `json:"period" bson:"period" form:"period"`           // 周期(秒), 仅Interval类型
	StartAtDate string        `json:"start_at" bson:"start_at" form:"start_at"`     // 开始日期 YYYY-MM-DD
	EndAtDate   string        `json:"end_at" bson:"end_at" form:"end_at"`           // 结束日期 YYYY-MM-DD
	StartTime   string        `json:"start_time" bson:"start_time" form:"start_time"` // 每日开始时间 HH:MM
	EndTime     string        `json:"end_time" bson:"end_time" form:"end_time"`     // 每日结束时间 HH:MM
	FixedTime   string        `json:"fixed_time" bson:"fixed_time" form:"fixed_time"` // 定时执行时间 HH:MM:SS
	Weekdays    []int         `json:"weekdays" bson:"weekdays" form:"weekdays"`     // 允许执行的周几(0-6)
	SkipDays    []string      `json:"skip_days" bson:"skip_days" form:"skip_days"`  // 跳过日期列表
	NextRunTime time.Time     `json:"next_run_time" bson:"next_run_time"`
	LastRunTime time.Time     `json:"last_run_time" bson:"last_run_time"`

	// 内部状态
	interval    time.Duration
	location    *time.Location
	startOfDay  time.Time
	endOfDay    time.Time
}

// NewTrigger 创建触发器
func NewTrigger(t TriggerType) *Trigger {
	loc := time.FixedZone("CST", 8*3600)
	now := time.Now().In(loc)
	return &Trigger{
		Type:        t,
		Enabled:     true,
		Period:      60,
		StartAtDate: now.Format("2006-01-02"),
		EndAtDate:   "2099-12-31",
		StartTime:   "00:00",
		EndTime:     "23:59",
		Weekdays:    []int{0, 1, 2, 3, 4, 5, 6},
		SkipDays:    []string{},
		location:    loc,
	}
}

// Refresh 刷新触发器内部状态
func (t *Trigger) Refresh() {
	loc := time.FixedZone("CST", 8*3600)
	t.location = loc

	t.startOfDay = parseDateInLocation(t.StartAtDate, loc)
	t.endOfDay = parseDateInLocation(t.EndAtDate, loc)
	t.interval = time.Duration(t.Period) * time.Second
}

// parseDateInLocation 解析日期字符串为指定时区的时间
func parseDateInLocation(date string, loc *time.Location) time.Time {
	t, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
	}
	return t
}

// ShouldRun 判断当前是否应该执行任务
func (t *Trigger) ShouldRun() bool {
	if !t.Enabled {
		return false
	}

	now := time.Now().In(t.location)

	// 检查日期范围
	if now.Before(t.startOfDay) || now.After(t.endOfDay) {
		return false
	}

	// 检查跳过日期
	for _, skip := range t.SkipDays {
		if now.Format("2006-01-02") == skip {
			return false
		}
	}

	// 检查周几
	if len(t.Weekdays) > 0 {
		weekday := int(now.Weekday())
		valid := false
		for _, w := range t.Weekdays {
			if w == weekday {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}

	// 检查日内时间范围
	startSec := parseTimeOfDay(t.StartTime)
	endSec := parseTimeOfDay(t.EndTime)
	currentSec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	if currentSec < startSec || currentSec > endSec {
		return false
	}

	return true
}

// CheckAndUpdate 检查是否可以执行并更新触发时间
func (t *Trigger) CheckAndUpdate() bool {
	if !t.ShouldRun() {
		return false
	}

	now := time.Now().In(t.location)

	switch t.Type {
	case TriggerOnce:
		if !t.LastRunTime.IsZero() {
			return false // 已经执行过
		}
		t.LastRunTime = now
		return true

	case TriggerInterval:
		if t.LastRunTime.IsZero() {
			t.LastRunTime = now
			t.NextRunTime = now.Add(t.interval)
			return true
		}
		if now.Sub(t.LastRunTime) >= t.interval {
			t.LastRunTime = now
			t.NextRunTime = now.Add(t.interval)
			return true
		}
		return false

	case TriggerDaily, TriggerWeekly:
		if !t.LastRunTime.IsZero() {
			lastRunDay := t.LastRunTime.Format("2006-01-02")
			if lastRunDay == now.Format("2006-01-02") {
				return false // 今天已执行
			}
		}
		if t.FixedTime != "" {
			fixedSec := parseTimeOfDay(t.FixedTime)
			currentSec := now.Hour()*3600 + now.Minute()*60 + now.Second()
			if currentSec >= fixedSec && currentSec <= fixedSec+10 {
				t.LastRunTime = now
				return true
			}
		}
		return false
	}

	return false
}

// TaskResult 任务执行结果
type TaskResult struct {
	StartTime time.Time      `json:"start_time" bson:"start_time"`
	EndTime   time.Time      `json:"end_time" bson:"end_time"`
	Status    string         `json:"status" bson:"status"`
	Message   string         `json:"message" bson:"message"`
	SpendTime float64        `json:"spend_time" bson:"spend_time"`
	Data      []bson.M       `json:"data" bson:"data"`
}

// NewTaskResult 创建任务结果
func NewTaskResult() TaskResult {
	return TaskResult{
		StartTime: time.Now(),
		Status:    StatusIdle,
		Data:      []bson.M{},
	}
}

// TaskInterface 任务执行接口
type TaskInterface interface {
	Execute(ctx context.Context) (TaskResult, error)
	Name() string
}

// BM 类型别名
type BM = bson.M
