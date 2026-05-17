package executor

import (
	"fmt"
	"time"
)

type DemoRunner1 struct {
	Alias string
}
type DemoRunner2 struct {
	Alias string
}

func (d *DemoRunner1) Name() string {
	return d.Alias
}
func (d *DemoRunner1) Run() (TaskResult, error) {
	fmt.Println("=====>>>> DemoRunner1 is running ....")
	start := time.Now()
	return TaskResult{
		Status:    StatusSuccess,
		StartTime: start,
		EndTime:   time.Now(),
	}, nil
}

func (d *DemoRunner1) Stop() error {
	return nil
}

func (d *DemoRunner2) Name() string {
	return d.Alias
}
func (d *DemoRunner2) Run() (TaskResult, error) {
	fmt.Println("=====>>>> DemoRunner2 is running ....")
	start := time.Now()
	return TaskResult{
		Status:    StatusSuccess,
		StartTime: start,
		EndTime:   time.Now(),
	}, nil
}

func (d *DemoRunner2) Stop() error {
	return nil
}
