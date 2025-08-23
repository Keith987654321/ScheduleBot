package main

import (
	"flag"
	"log"

	"github.com/Keith987654321/schedule-tg-bot/bot"
	"github.com/Keith987654321/schedule-tg-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	flagU := flag.String("user", "postgres", "db username")
	flagP := flag.String("pass", "qwerty123", "db password")
	flagN := flag.String("name", "postgres", "db name")
	flagSSL := flag.String("ssl", "disable", "SSL mode")
	flagToken := flag.String("token", "", "Telegram bot token")

	flag.Parse()

	db.Connect(*flagU, *flagP, *flagN, *flagSSL)

	botAPI, err := tgbotapi.NewBotAPI(*flagToken)
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
