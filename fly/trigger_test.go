package fly

import (
	"testing"
	"time"
)

func baseActiveTrigger() *Trigger {
	tg := NewTrigger(Interval)
	now := time.Now()
	tg.start = now.Add(-24 * time.Hour)
	tg.end = now.Add(24 * time.Hour)
	tg.startSecond = 0
	tg.endSecond = 24 * 3600
	tg.Weekdays = nil
	tg.SkipDays = nil
	tg.rangeTime = nil
	tg.RangeTime = nil
	return tg
}

func secondsOfDay(t time.Time) int {
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func TestTriggerActive_DateRange(t *testing.T) {
	tg := baseActiveTrigger()
	if !tg.Active() {
		t.Fatalf("expected trigger to be active inside date range")
	}

	now := time.Now()
	tg.start = now.Add(24 * time.Hour)
	tg.end = now.Add(48 * time.Hour)
	if tg.Active() {
		t.Fatalf("expected trigger to be inactive when now is before start date")
	}
}

func TestTriggerActive_SkipDays(t *testing.T) {
	tg := baseActiveTrigger()
	tg.SkipDays = []string{time.Now().Format("2006-01-02")}
	if tg.Active() {
		t.Fatalf("expected trigger to be inactive when today is in skip_days")
	}
}

func TestTriggerActive_Weekdays(t *testing.T) {
	tg := baseActiveTrigger()
	tg.Weekdays = []time.Weekday{time.Now().Weekday()}
	if !tg.Active() {
		t.Fatalf("expected trigger to be active on allowed weekday")
	}

	notToday := (time.Now().Weekday() + 1) % 7
	tg.Weekdays = []time.Weekday{notToday}
	if tg.Active() {
		t.Fatalf("expected trigger to be inactive on disallowed weekday")
	}
}

func TestTriggerActive_DayTimeWindow(t *testing.T) {
	tg := baseActiveTrigger()
	nowSec := secondsOfDay(time.Now())

	tg.startSecond = nowSec - 10
	tg.endSecond = nowSec - 9
	if tg.Active() {
		t.Fatalf("expected trigger to be inactive outside start_time/end_time window")
	}

	tg.startSecond = 0
	tg.endSecond = 24 * 3600
	if !tg.Active() {
		t.Fatalf("expected trigger to be active inside start_time/end_time window")
	}
}

func TestTriggerActive_RangeTime(t *testing.T) {
	tg := baseActiveTrigger()
	nowSec := secondsOfDay(time.Now())

	tg.rangeTime = [][]int{{nowSec + 10, nowSec + 20}}
	if tg.Active() {
		t.Fatalf("expected trigger to be inactive when now is outside all range_time windows")
	}

	tg.rangeTime = [][]int{{nowSec - 10, nowSec + 10}}
	if !tg.Active() {
		t.Fatalf("expected trigger to be active when now is inside one range_time window")
	}
}

func TestTriggerRunInterval(t *testing.T) {
	tg := baseActiveTrigger()
	tg.SetPeriod(3)

	time.Sleep(3 * time.Second)
	if !tg.RunInterval() {
		t.Fatalf("expected first RunInterval call to run")
	}

	if tg.RunInterval() {
		t.Fatalf("expected second immediate RunInterval call not to run")
	}

	stale := time.Now().Add(-3 * time.Second)
	tg.LastRunTime = stale
	if !tg.RunInterval() {
		t.Fatalf("expected overdue RunInterval call to return true with current implementation")
	}
	if !tg.LastRunTime.After(stale) {
		t.Fatalf("expected LastRunTime to be refreshed when interval is overdue")
	}
}

func TestTriggerTimeIsUp_Switch(t *testing.T) {
	tg := baseActiveTrigger()
	tg.SetPeriod(1)
	tg.Type = Interval
	if !tg.TimeIsUp() {
		t.Fatalf("expected TimeIsUp to run for Interval type")
	}

	now := time.Now()
	inactiveStart := now.Add(24 * time.Hour)
	inactiveEnd := now.Add(48 * time.Hour)
	for _, tt := range []int{Once, Everyday, Everyweek, Month, Cron} {
		tgCase := baseActiveTrigger()
		tgCase.Type = tt
		tgCase.start = inactiveStart
		tgCase.end = inactiveEnd
		if tgCase.TimeIsUp() {
			t.Fatalf("expected TimeIsUp to be false for inactive trigger type=%d", tt)
		}
	}
}
