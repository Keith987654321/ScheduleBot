package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/Keith987654321/schedule-tg-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

var cronRef *cron.Cron

func InitCron(c *cron.Cron, bot *tgbotapi.BotAPI) error {
	cronRef = c
	return scheduleMessages(cronRef, bot)
}

func scheduleMessages(c *cron.Cron, bot *tgbotapi.BotAPI) error {
	messages, err := db.GetScheduledMessages()
	if err == nil {
		for _, m := range messages {
			cronSpec, msg, subgroup := m.CronSpec, m.Message, m.Subgroup
			c.AddFunc(cronSpec, func() { sendScheduledMessages(bot, msg, subgroup) })
		}
	}

	return err
}

func addScheduledMessage(bot *tgbotapi.BotAPI, spec string, message string, subgroup int) error {
	cronEntryID, err := cronRef.AddFunc(spec, func() { sendScheduledMessages(bot, message, subgroup) })
	log.Printf("cron entry id: %v\n", cronEntryID)
	if err != nil {
		return err
	}

	return db.AddScheduledMessage(spec, cronEntryID, message, subgroup)
}

func convertTimeToCronSpec(day, pairNumber int) string {
	pairTime := []string{"", "8:30", "10:10", "11:50", "13:50", "15:30", "17:10", "18:50"}
	arr := strings.Split(pairTime[pairNumber], ":")
	hours, minutes := arr[0], arr[1]
	weekday := day - 1

	return fmt.Sprintf("* %s %s * * %d", minutes, hours, weekday)
}

func sendScheduledMessages(bot *tgbotapi.BotAPI, message string, subgroup int) {
	users, err := db.GetUsersBySubgroup(subgroup)
	if err != nil {
		log.Printf("can't get users by subgroup: %v\n", err)
	}

	limiter := make(chan struct{}, 20) // Avoiding tgapi rate limit ~ 30 per sec
	for _, user := range users {
		limiter <- struct{}{}
		go func(userID int64) {
			defer func() { <-limiter }()
			msg := tgbotapi.NewMessage(userID, message)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("sending error %d: %v\n", userID, err)

			}
		}(user.TelegramID)
	}
}

func deleteScheduledMessage(schedmesID int) error {
	schedmes, err := db.GetScheduledMessage(schedmesID)
	if err != nil {
		return err
	}
	cronRef.Remove(schedmes.CronEntryID)
	return db.DeleteScheduledMessage(schedmesID)
}
