package fly

import (
	"fmt"
)

// Start 启动任务调度器
func Start() {
	tm, err := Init()
	if err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}

	fmt.Println("Starting task scheduler...")
	tm.Start()
}
