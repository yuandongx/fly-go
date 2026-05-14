package executor

import (
	"strconv"
	"time"
)

func toInt(s string) int {
	// 将 "09:00" 转换为秒数 (32400)
	hour, _ := strconv.Atoi(s[:2])
	minute, _ := strconv.Atoi(s[3:])
	return hour*3600 + minute*60
}

func toIntDay(slist []string) []int {
	// 将 ["2026-01-01"] 转换成 int 时间戳
	ret := make([]int, 0, len(slist))
	for _, s := range slist {
		if t, err := parseDate(s); err == nil {
			ret = append(ret, int(t.Unix()))
		}
	}
	return ret
}

// parseDate 将 "2026-01-01" 解析为 time.Time
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// 参数初始化
func (tp *TaskParam) Init(params map[string]any) {
	// 将 params 里的参数转换为 TaskParam
	if v, ok := params["start_time"].(string); ok {
		tp.StartTime = toInt(v)
	}
	if v, ok := params["end_time"].(string); ok {
		tp.EndTime = toInt(v)
	}
	if v, ok := params["skip_dates"].([]any); ok {
		tp.SkipDates = toIntDay(toAnyStrings(v))
	}
	if v, ok := params["skip_weekdays"].([]any); ok {
		tp.SkipWeekdays = toInts(v)
	}
	if v, ok := params["type"].(string); ok {
		tp.Type = v
	}
	if v, ok := params["interval"].(int); ok {
		tp.interval = time.Duration(v) * time.Second
	}
	// 开始日期，解析为 time.Time
	if v, ok := params["start_date"].(string); ok {
		tp.StartDate, _ = parseDate(v)
	}
	// 结束日期，解析为 time.Time
	if v, ok := params["end_date"].(string); ok {
		tp.EndDate, _ = parseDate(v)
	}
}

func toAnyStrings(in []any) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toInts(in []any) []int {
	out := make([]int, 0, len(in))
	for _, v := range in {
		switch val := v.(type) {
		case int:
			out = append(out, val)
		case float64:
			out = append(out, int(val))
		}
	}
	return out
}

// Active 检查当前时间是否在任务参数定义的有效范围内
// 按以下顺序检查：
//   1. 日期范围：当前日期是否在 StartDate ~ EndDate 内
//   2. 跳过日期：当前日期是否在 SkipDates 中
//   3. 星期范围：当前星期是否在 SkipWeekdays 中（SkipWeekdays 为要跳过的星期）
//   4. 时间范围：当前时间是否在 StartTime ~ EndTime 内
func (tp *TaskParam) Active() bool {
	now := time.Now()

	// 1. 检查日期范围
	if !tp.StartDate.IsZero() && now.Before(tp.StartDate) {
		return false
	}
	if !tp.EndDate.IsZero() && now.After(tp.EndDate) {
		return false
	}

	// 2. 检查跳过日期
	for _, skipTs := range tp.SkipDates {
		skipDate := time.Unix(int64(skipTs), 0)
		if now.Year() == skipDate.Year() &&
			now.Month() == skipDate.Month() &&
			now.Day() == skipDate.Day() {
			return false
		}
	}

	// 3. 检查跳过星期 (SkipWeekdays 存储的是要跳过的星期)
	if len(tp.SkipWeekdays) > 0 {
		weekday := int(now.Weekday())
		for _, skip := range tp.SkipWeekdays {
			if skip == weekday {
				return false
			}
		}
	}

	// 4. 检查日内时间范围
	if tp.StartTime > 0 || tp.EndTime > 0 {
		currentSec := now.Hour()*3600 + now.Minute()*60 + now.Second()
		if tp.StartTime > 0 && currentSec < tp.StartTime {
			return false
		}
		if tp.EndTime > 0 && currentSec > tp.EndTime {
			return false
		}
	}

	return true
}
