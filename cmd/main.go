package main

import (
	"os"
	"who-counters-bot/internal/bot"
)

func main() {

	var err error

	token := os.Getenv("TOKEN")

	for err == nil {
		err = bot.TelegramBot(token) //TODO: replace with config
	}

	if err != nil {
		panic(err)
	}
}
