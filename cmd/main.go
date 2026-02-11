// Package main provides the entry point for the application.
package main

import (
	"flag"
	"fly-go/fly"
	"fmt"
)

// main function

var app = flag.String("service", "spider", "Select spider|server to run it.")
var port = flag.Int("port", 8000, "Set the port when run server.")

func main() {
	flag.Parse()
	switch *app {
	case "spider":
		{
			fly.Start()
			fmt.Println("Spider is starting...")
		}
	case "server":
		{

			fmt.Println("Server is starting...")
			Server(*port)
			fmt.Printf("Server is running with port %d ...\n", *port)
		}
	default:
		fmt.Println("Only `spider` or `server` can be selected.")
	}
}
