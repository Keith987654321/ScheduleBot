package keyboards

import (
	"fmt"

	"github.com/Keith987654321/schedule-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	subgroups = [2]int{1, 2}
	days      = [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
)

func CreateMainMenu(role string) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	if role == "admin" {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("/schedule"), tgbotapi.NewKeyboardButton("/today")},
			{tgbotapi.NewKeyboardButton("/teachers"), tgbotapi.NewKeyboardButton("/change_group")},
			{tgbotapi.NewKeyboardButton("/approve"), tgbotapi.NewKeyboardButton("/reject")},
			{tgbotapi.NewKeyboardButton("/edit"), tgbotapi.NewKeyboardButton("/delete")},
		}
	} else {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("/today")}, {tgbotapi.NewKeyboardButton("/schedule"), tgbotapi.NewKeyboardButton("/teachers")},
			{tgbotapi.NewKeyboardButton("/suggest"), tgbotapi.NewKeyboardButton("/change_group")},
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
