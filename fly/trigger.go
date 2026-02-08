package fly

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	Interval = 1 << iota
	Once
	Everyday
	Everyweek
	Month
	Cron
)

type Trigger struct {
	Type        int            `json:"type" bson:"type,omitempty" form:"type" validate:"required"`
	Period      int64          `json:"period" bson:"period,omitempty" form:"period"` // 仅在Type为Interval时有效，单位为秒
	Enabled     bool           `json:"enabled" bson:"enabled,omitempty" form:"enabled"`
	StartAtDate string         `json:"start_at" bson:"start_at,omitempty" form:"start_at"`       // 任务开始日期，格式为 "YYYY-MM-DD"
	EndAtDate   string         `json:"end_at" bson:"end_at,omitempty" form:"end_at"`             // 任务结束日期，格式为 "YYYY-MM-DD"
	StartTime   string         `json:"start_time" bson:"start_time,omitempty" form:"start_time"` // 每天的开始时间，格式为 "HH:MM"
	EndTime     string         `json:"end_time" bson:"end_time,omitempty" form:"end_time"`
	Weekdays    []time.Weekday `json:"weekdays" bson:"weekdays,omitempty" form:"weekdays"`
	SkipDays    []string       `json:"skip_days" bson:"skip_days,omitempty" form:"skip_days"`
	NextRunTime time.Time      `json:"next_run_time" bson:"next_run_time,omitempty"`
	LastRunTime time.Time      `json:"last_run_time" bson:"last_run_time,omitempty"`
	RangeTime   [][]string     `json:"range_time" bson:"range_time,omitempty" form:"range_time"`

	rangeTime   [][]int
	start       time.Time
	end         time.Time
	startSecond int
	endSecond   int
	interval    time.Duration
	local       *time.Location
}

func parseTime(t string) (min int, ok bool) {
	h := 0
	m := 0
	ok = true
	if regx, err := regexp.Compile(`(\d{2}):(\d{2})`); err == nil {
		tmp := regx.FindStringSubmatch(t)
		if len(tmp) > 2 {
			if _h, err := strconv.Atoi(tmp[1]); err == nil {
				h = _h
			} else {
				ok = false
			}
			if _m, err := strconv.Atoi(tmp[2]); err == nil {
				m = _m
			} else {
				ok = false
			}
		}
	}
	min = 3600*h + 60*m
	return
}
func parseDate(date string) (y, m, d int, ok bool) {
	y = 2026
	m = 1
	d = 1
	ok = true
	if rgex, err := regexp.Compile(`(\d{4})-(\d{2})-(\d{2})`); err == nil {
		res := rgex.FindStringSubmatch(date)
		if _y, err := strconv.Atoi(res[1]); err == nil {
			y = int(_y)
		} else {
			ok = false
		}
		if _d, err := strconv.Atoi(res[2]); err == nil {
			d = int(_d)
		} else {
			ok = false
		}
		if _m, err := strconv.Atoi(res[3]); err == nil {
			m = int(_m)
		} else {
			ok = false
		}
	}
	return
}
func NewTrigger(t int) *Trigger {
	local := time.FixedZone("utc+8", 8*3600)
	tg := &Trigger{Type: t,
		Enabled:     true,
		StartTime:   "00:00:00",
		EndTime:     "23:59:59",
		LastRunTime: time.Time{},
		local:       local,
		startSecond: 0,
		endSecond:   24 * 3600,
		rangeTime:   make([][]int, 0),
		Weekdays:    make([]time.Weekday, 0),
		SkipDays:    make([]string, 0)}
	now := time.Now()
	tg.SetEndDate("2099-01-01")
	tg.SetSatrtDate(now.Format("2006-01-02"))
	return tg
}

func (tg *Trigger) Refresh() {
	tg.SetSatrtDate(tg.StartAtDate)

	tg.SetEndDate(tg.EndAtDate)

	tg.SetSatrtTime(tg.StartTime)

	tg.SetEndTime(tg.EndTime)
	// tg.SetWeekDays(tg.Weekdays)

	// tg.SetSkipDays(tg.SkipDays)

	// 更新RangeTime
	tg.RangeTime = make([][]string, 0)
	tg.rangeTime = make([][]int, 0)
	for _, rt := range tg.RangeTime {
		tg.SetRangeTime(rt[0], rt[1])
	}

	tg.SetPeriod(int(tg.Period))

}
func (tg *Trigger) SetRangeTime(start, end string) {
	s := 0
	e := 24 * 3600
	if m, ok := parseTime(start); ok {
		s = m
	}
	if m, ok := parseTime(end); ok {
		e = m
	}
	tg.rangeTime = append(tg.rangeTime, []int{s, e})
	tg.RangeTime = append(tg.RangeTime, []string{start, end})
}
func (tg *Trigger) SetSatrtDate(date string) {
	tg.StartAtDate = date
	year2, month2, day2, ok3 := parseDate(date)
	if !ok3 {
		year2 = 2026
		month2 = 1
		day2 = 1
	}
	tg.start = time.Date(year2, time.Month(month2), day2, 0, 0, 0, 0, tg.local)
}

func (tg *Trigger) SetSatrtTime(t string) {
	tg.StartTime = t
	if m, ok := parseTime(t); ok {
		tg.startSecond = m
	}
}

func (tg *Trigger) SetEndDate(date string) {
	tg.EndAtDate = date
	year3, month3, day3, ok4 := parseDate(date)
	if !ok4 {
		year3 = 2099
		month3 = 1
		day3 = 1
	}
	tg.end = time.Date(year3, time.Month(month3), day3, 0, 0, 0, 0, tg.local)
}

func (tg *Trigger) SetEndTime(t string) {
	tg.EndTime = t
	if m, ok := parseTime(t); ok {
		tg.endSecond = m
	}
}

func (tg *Trigger) SetWeekDays(days []int) {
	for i := range days {
		if i >= 0 && i <= 6 {
			tg.Weekdays = append(tg.Weekdays, time.Weekday(i))
		}
	}
}
func (tg *Trigger) SetSkipDays(days []string) {
	for _, day := range days {
		if day != "" {
			tg.SkipDays = append(tg.SkipDays, day)
		}
	}
}

func (tg *Trigger) SetPeriod(second int) {
	s := fmt.Sprintf("%ds", second)
	if d, err := time.ParseDuration(s); err == nil {
		tg.interval = d
	}
}

// Active 哪些天是要跳过，不执行的
func (tg *Trigger) Active() bool {
	now := time.Now()
	t := now.Hour()*3600 + now.Minute()*60 + now.Second()
	// 有效期内
	ok1 := now.Before(tg.end) && now.After(tg.start)

	// 不在跳过日期内
	ok2 := true
	if len(tg.SkipDays) > 0 {
		for _, d := range tg.SkipDays {
			if d == now.Format("2006-01-02") {
				ok2 = false
			}
		}
	}
	// 一天的有效时间段内
	ok3 := t >= tg.startSecond && t <= tg.endSecond

	// 有效周内
	ok4 := len(tg.Weekdays) == 0
	for _, week := range tg.Weekdays {
		if week == now.Weekday() {
			ok4 = true
			break
		}
	}
	// 一天内的有效时间
	ok5 := len(tg.rangeTime) == 0
	for _, r := range tg.rangeTime {
		if t >= r[0] && t <= r[1] {
			ok5 = true
			break
		}
	}
	return ok1 && ok2 && ok3 && ok4 && ok5
}

// RunInterval 固定时间间隔
func (tg *Trigger) RunInterval() bool {
	if tg.Active() {

		now := time.Now()
		if tg.LastRunTime.IsZero() {
			tg.LastRunTime = now
			tg.NextRunTime = now.Add(tg.interval)
			return true
		} else {
			sub1 := now.Sub(tg.LastRunTime)
			if sub1 < tg.interval {
				return false
			} else if sub1 >= tg.interval {
				tg.NextRunTime = now.Add(tg.interval)
				tg.LastRunTime = now
				return true
			}
		}

	}
	return false
}

// RunOnce 只执行一次
func (tg *Trigger) RunOnce() bool {
	if tg.LastRunTime.IsZero() {
		now := time.Now()
		sub := now.Sub(tg.start)
		if sub < 3 && sub > -3 {
			tg.LastRunTime = now
			return true
		}
	}

	return false
}

func (tg *Trigger) RunEveryday() bool {
	if !tg.Active() {
		return false
	}
	now := time.Now()
	sub := now.Sub(tg.start)
	if sub < -3 || sub > 3 {
		return false
	}

	return true
}

// RunEveryweek 每周的一些天执行
func (tg *Trigger) RunEveryweek() bool {
	now := time.Now()
	ok := false
	weekday := now.Weekday()
	for _, d := range tg.Weekdays {
		if d == weekday {
			ok = true
			break
		}
	}
	sub := now.Sub(tg.start)
	flag := sub < 3 && sub > -3
	return ok && tg.Active() && flag
}
func (tg *Trigger) TimeIsUp() bool {
	ok := false
	switch tg.Type {
	case Interval:
		return tg.RunInterval()
	case Once:
		return tg.RunOnce()
	case Everyday:
		return tg.RunEveryday()
	case Everyweek:
		return tg.RunEveryweek()
	case Month:
	}
	return ok
}
