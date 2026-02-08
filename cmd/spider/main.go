package main

import (
	"fly-go/fly"
	"os"
)

func main() {
	// 这里是 Spider 服务的入口点，可以根据需要添加 Spider 相关的初始化和启动逻辑
	os.Setenv("FLY_CONFIG", "D:\\code\\web-app\\fly-go\\config.yaml")
	fly.Start()
}
