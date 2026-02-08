// Package main provides the entry point for the application.
package main

import (
	"flag"
	"fly-go/fly"
	"fly-go/server"
	"fmt"
)

// main function

var mode = flag.String("mode", "server", "Select spider|server to run it.")
var port = flag.Int("port", 8000, "Set the port when run server.")

// main 是应用程序的入口点，根据命令行参数启动不同的服务模式
// 支持两种模式：
//   - "spider": 启动爬虫服务
//   - "server": 启动服务器服务
//
// 默认情况下同时启动爬虫和服务器
func main() {

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
	flag.Parse()
	switch *mode {
	case "spider":
		{
			fly.Start()
			fmt.Println("Spider is starting...")
		}
	case "server":
		{

			fmt.Println("Server is starting...")
			fmt.Printf("Server is running with port %d ...\n", *port)
			server.Start(*port)
		}
	default:
		{
			fmt.Println("Only `spider` or `server` can be selected.")
			fly.Start()
			// server.Start(*port)
			fmt.Printf("Server is running with port %d ...\n", *port)
		}
	}
}
