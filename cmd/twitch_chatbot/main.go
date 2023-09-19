package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vikpe/qw-demobot/internal/app/twitch_chatbot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("unable to load environment variables", err)
		return
	}

	chatbot := twitch_chatbot.New(
		os.Getenv("TWITCH_CHANNEL_USERNAME"),
		os.Getenv("TWITCH_CHANNEL_ACCESS_TOKEN"),
		os.Getenv("ZMQ_SUBSCRIBER_ADDRESS"),
		os.Getenv("ZMQ_PUBLISHER_ADDRESS"),
	)
	chatbot.Start()
}
