package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/qw-demobot/internal/app/quake_manager"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	manager := quake_manager.New(
		os.Getenv("EZQUAKE_BIN_PATH"),
		os.Getenv("EZQUAKE_PROCESS_USERNAME"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
	)
	manager.Start()
}
