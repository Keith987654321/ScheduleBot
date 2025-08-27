package main

import (
	"log"
	"os"

	"github.com/Keith987654321/schedule-tg-bot/bot"
	"github.com/Keith987654321/schedule-tg-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	ssl := os.Getenv("SSL_MODE")
	token := os.Getenv("TOKEN")

	db.Connect(username, host, port, pass, name, ssl)

	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		bot.HandleUpdate(botAPI, update)
	}
}
