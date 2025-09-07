package keyboards

import (
	"fmt"

	"github.com/Keith987654321/schedule-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	subgroups = [2]int{1, 2}
	days      = [6]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}
	Commands  = map[string]string{
		"Сегодня": "/today", "Завтра": "/tomorrow", "Расписание": "/schedule",
		"Преподы": "/teachers", "Предложить": "/suggest", "Сменить подгруппу": "/change_group",
		"Одобрить": "/approve", "Отклонить": "/reject", "Редактировать": "/edit", "Удалить": "/delete",
	}
)

func CreateMainMenu(role string) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	if role == "admin" {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Сегодня"), tgbotapi.NewKeyboardButton("Завтра")},
			{tgbotapi.NewKeyboardButton("Преподы"), tgbotapi.NewKeyboardButton("Расписание"), tgbotapi.NewKeyboardButton("Сменить подгруппу")},
			{tgbotapi.NewKeyboardButton("Одобрить"), tgbotapi.NewKeyboardButton("Отклонить")},
			{tgbotapi.NewKeyboardButton("Редактировать"), tgbotapi.NewKeyboardButton("Удалить")},
		}
	} else {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Сегодня"), tgbotapi.NewKeyboardButton("Завтра")}, {tgbotapi.NewKeyboardButton("Расписание"), tgbotapi.NewKeyboardButton("Преподы")},
			{tgbotapi.NewKeyboardButton("Предложить"), tgbotapi.NewKeyboardButton("Сменить подгруппу")},
		}
	}
	return tgbotapi.NewReplyKeyboard(buttons...)
}

func CreateDaySelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for i, day := range days {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(day, fmt.Sprintf("schedule_%d", i+1)),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func CreateSubgroupSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, subgroup := range subgroups {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", subgroup), fmt.Sprintf("subgroup_%d", subgroup)),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func CreateChangeSubgroupSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, subgroup := range subgroups {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", subgroup), fmt.Sprintf("change_group_%d", subgroup)),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func CreateSuggestionsKeyboard(suggs []models.Suggestion, action string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, sugg := range suggs {
		buttonText := fmt.Sprintf("ID %d: День %d, Пара %d -> %s, подгруппа %d", sugg.ID, sugg.DayOfWeek, sugg.PairNumber, sugg.NewSubject, sugg.Subgroup)
		callbackData := fmt.Sprintf("%s_%d", action, sugg.ID)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
