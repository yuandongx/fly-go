package fly_test

import (
	"fly-go/fly"
	"fmt"
	"testing"
	"time"
)

func TestTrigger_TimeIsUp(t *testing.T) {
	tg := fly.NewTrigger(1)
	tg.SetPeriod(3)
	tg.SetSatrtTime("13:37")
	tg.SetEndTime("13:39")
	tg.SetWeekDays([]int{0, 1, 2})
	for i := 0; i < 10; i++ {
		ok := tg.TimeIsUp()
		next := tg.LastRunTime
		last := tg.LastRunTime
		fmt.Printf("ok=%v next=%v last=%v\n", ok, next, last)
		time.Sleep(6 * time.Second)
	}
}
