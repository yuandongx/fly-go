// Package main provides the entry point for the application.
package main

import (
	"flag"
	"fly-go/fly"
	"fly-go/server"
	"fmt"
)

var mode = flag.String("mode", "server", "Select spider|server to run it.")
var port = flag.Int("port", 8080, "Set the port when run server.")

// ModeHandler 处理不同运行模式的接口
type ModeHandler interface {
	StartSpider()
	StartServer(port int)
}

// DefaultHandler 默认的 ModeHandler 实现
type DefaultHandler struct{}

func (h *DefaultHandler) StartSpider() {
	fly.Start()
}

func (h *DefaultHandler) StartServer(port int) {
	server.Start(port)
}

// RunMode 根据模式执行相应的服务
// mode: 运行模式 ("spider", "server", 或其他)
// port: 服务器端口号
// handler: ModeHandler 接口，用于解耦测试
func RunMode(mode string, port int, handler ModeHandler) {
	switch mode {
	case "spider":
		handler.StartSpider()
		fmt.Println("Spider is starting...")
	case "server":
		fmt.Println("Server is starting...")
		fmt.Printf("Server is running with port %d ...\n", port)
		handler.StartServer(port)
	default:
		fmt.Println("Only `spider` or `server` can be selected.")
		handler.StartSpider()
		fmt.Printf("Server is running with port %d ...\n", port)
	}
}

// main 是应用程序的入口点，根据命令行参数启动不同的服务模式
func main() {
	printBanner()
	flag.Parse()
	handler := &DefaultHandler{}
	RunMode(*mode, *port, handler)
}

func printBanner() {
	println(`
		________  __    __      __         ______  ________   ______   _______  ________ 
		|        \|  \  |  \    /  \       /      \|        \ /      \ |       \|        \
		| $$$$$$$$| $$   \$$\  /  $$      |  $$$$$$\\$$$$$$$$|  $$$$$$\| $$$$$$$\\$$$$$$$$
		| $$__    | $$    \$$\/  $$       | $$___\$$  | $$   | $$__| $$| $$__| $$  | $$   
		| $$  \   | $$     \$$  $$         \$$    \   | $$   | $$    $$| $$    $$  | $$   
		| $$$$$   | $$      \$$$$          _\$$$$$$\  | $$   | $$$$$$$$| $$$$$$$\  | $$   
		| $$      | $$_____ | $$          |  \__| $$  | $$   | $$  | $$| $$  | $$  | $$   
		| $$      | $$     \| $$           \$$    $$  | $$   | $$  | $$| $$  | $$  | $$   
		\$$       \$$$$$$$$ \$$            \$$$$$$    \$$    \$$   \$$ \$$   \$$   \$$
	`)
}
