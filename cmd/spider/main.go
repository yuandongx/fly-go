package main

import (
	"fmt"
	"fly-go/fly"
	"os"
)

func main() {
	os.Setenv("FLY_CONFIG", "D:\\code\\web-app\\fly-go\\config.yaml")

	// 初始化并启动
	tm, err := fly.Init()
	if err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return
	}

	// 启动调度
	go tm.Start()

	fmt.Println("Spider service started. Press Ctrl+C to stop.")

	// 阻塞
	select {}
}
