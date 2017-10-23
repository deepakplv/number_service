package main

import (
	"flag"
	"fmt"
	"os"
	"config"
	"data-adapters"
	"logger"
	"server"
)

func main() {
	// Load config based on the server mode
	environment := flag.String("e", "DEV", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)

	logger.Init()
	adapters.Init()
	server.Init()
}
