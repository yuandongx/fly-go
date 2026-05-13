package fly

import (
	"regexp"
	"strconv"
	"time"
)

// parseTimeOfDay 解析 HH:MM 或 HH:MM:SS 格式时间，返回当天秒数
func parseTimeOfDay(t string) int {
	re := regexp.MustCompile(`(\d{2}):(\d{2}):?(\d{2})?`)
	matches := re.FindStringSubmatch(t)
	if len(matches) < 3 {
		return 0
	}

	h, _ := strconv.Atoi(matches[1])
	m, _ := strconv.Atoi(matches[2])
	s := 0
	if len(matches) > 3 && matches[3] != "" {
		s, _ = strconv.Atoi(matches[3])
	}

	return h*3600 + m*60 + s
}

// formatDuration 格式化时长为可读字符串
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Round(time.Second).String()
	}
	if d < time.Hour {
		return d.Round(time.Minute).String()
	}
	return d.Round(time.Hour).String()
}

// weekdayNames 星期名称映射
var weekdayNames = map[time.Weekday]string{
	time.Sunday:    "周日",
	time.Monday:    "周一",
	time.Tuesday:   "周二",
	time.Wednesday: "周三",
	time.Thursday:  "周四",
	time.Friday:    "周五",
	time.Saturday:  "周六",
}

// WeekdayName 获取星期名称
func WeekdayName(w time.Weekday) string {
	if name, ok := weekdayNames[w]; ok {
		return name
	}
	return "未知"
}
