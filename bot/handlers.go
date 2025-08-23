package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Keith987654321/schedule-tg-bot/db"
	"github.com/Keith987654321/schedule-tg-bot/models"
	"github.com/Keith987654321/schedule-tg-bot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createMainMenu(role string) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	if role == "admin" {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("/schedule"), tgbotapi.NewKeyboardButton("/today")},
			{tgbotapi.NewKeyboardButton("/approve"), tgbotapi.NewKeyboardButton("/reject")},
			{tgbotapi.NewKeyboardButton("/edit"), tgbotapi.NewKeyboardButton("/delete")},
		}
	} else {
		buttons = [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("/today")}, {tgbotapi.NewKeyboardButton("/schedule"), tgbotapi.NewKeyboardButton("/suggest")},
		}
	}
	return tgbotapi.NewReplyKeyboard(buttons...)
}

func createDaySelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	days := []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for i, day := range days {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(day, fmt.Sprintf("schedule_%d", i+1)),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func createSuggestionsKeyboard(suggs []models.Suggestion, action string) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, sugg := range suggs {
		buttonText := fmt.Sprintf("ID %d: День %d, Пара %d -> %s", sugg.ID, sugg.DayOfWeek, sugg.PairNumber, sugg.NewSubject)
		callbackData := fmt.Sprintf("%s_%d", action, sugg.ID)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		})
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig
	var user *models.User
	var err error
	var chatID int64

	if update.Message != nil {
		chatID = update.Message.Chat.ID
		msg = tgbotapi.NewMessage(chatID, "")
		user, err = db.GetUserByTelegramID(update.Message.From.ID)
		if err != nil {
			msg.Text = "Ошибка: 66" + err.Error()
			bot.Send(msg)
			return
		}
		db.CheckUserInfo(update.Message.From.ID, update.Message.From.FirstName, update.Message.From.UserName)
		msg.ReplyMarkup = createMainMenu(user.Role)
	} else if update.CallbackQuery != nil {
		if update.CallbackQuery.Message == nil {
			log.Printf("CallbackQuery has no Meesage: %v", update.CallbackQuery)
			return
		}
		chatID = update.CallbackQuery.Message.Chat.ID
		msg = tgbotapi.NewMessage(chatID, "")
		user, err = db.GetUserByTelegramID(update.CallbackQuery.From.ID)
		if err != nil {
			msg.Text = "Ошибка" + err.Error()
			bot.Send(msg)
			return
		}
		// Handler for callbacks from inline keyboard
		callbackData := update.CallbackQuery.Data
		if strings.HasPrefix(callbackData, "schedule_") {
			day, _ := strconv.Atoi(strings.TrimPrefix(callbackData, "schedule_"))
			items, err := db.GetScheduleForDay(day)
			if err != nil {
				msg.Text = "Не получилось получить расписание."
			} else {
				msg.Text = sprintSchedule(day, items)
			}
			msg.ReplyMarkup = createMainMenu(user.Role)
			bot.Send(msg)
			// Deleting inline keayboard
			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			bot.Send(editMsg)
			return
		} else if strings.HasPrefix(callbackData, "approve_") && user.Role == "admin" {
			suggID, _ := strconv.Atoi(strings.TrimPrefix(callbackData, "approve_"))
			err := db.ApproveSuggestion(suggID)
			if err != nil {
				msg.Text = "Ошибка одобрения."
			} else {
				msg.Text = "Предложение одобрено и применено."
			}
			msg.ReplyMarkup = createMainMenu(user.Role)
			bot.Send(msg)

			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			bot.Send(editMsg)
			return
		} else if strings.HasPrefix(callbackData, "reject_") && user.Role == "admin" {
			suggID, _ := strconv.Atoi(strings.TrimPrefix(callbackData, "reject_"))
			err := db.RejectSuggestion(suggID)
			if err != nil {
				msg.Text = "Ошибка отклонения."
			} else {
				msg.Text = "Предложение отклонено."
			}
			msg.ReplyMarkup = createMainMenu(user.Role)
			bot.Send(msg)

			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			bot.Send(editMsg)
			return
		}
		return
	}
	if chatID == 0 {
		log.Printf("Error: chat_id is empty for update: %v", update)
	}

	text := strings.ToLower(update.Message.Text)
	today := int(time.Now().Weekday()) // 0=Sunday, 1=Monday
	if today == 0 {
		today = 7
	}

	switch {
	case strings.HasPrefix(text, "/today"):
		// Format /today
		items, err := db.GetScheduleForDay(today)
		if err != nil {
			msg.Text = "Не получилось получить расписание."
		} else {
			msg.Text = sprintSchedule(today, items)
		}

	case strings.HasPrefix(text, "/schedule"):
		// Format: /schedule [day]
		parts := strings.SplitN(text, " ", 2)
		day := today
		if len(parts) == 2 {
			parsedDay, err := strconv.Atoi(parts[1])
			if err != nil || parsedDay < 1 || parsedDay > 7 {
				msg.Text = "Укажите день недели (1-7, где 1 — понедельник, 7 — воскресенье)."
				bot.Send(msg)
				return
			}
			day = parsedDay
			items, err := db.GetScheduleForDay(day)
			if err != nil {
				msg.Text = "Не получилось получить расписание."
			} else {
				msg.Text = sprintSchedule(day, items)
			}
		} else {
			msg.Text = "Выберите день недели:"
			msg.ReplyMarkup = createDaySelectionKeyboard()
		}

	case strings.HasPrefix(text, "/suggest") && (user.Role == "user" || user.Role == "admin"):
		err, item := parseSuggestion(text)
		if err != nil {
			msg.Text = handleParseErrors(err)
			bot.Send(msg)
			return
		}
		err = db.SuggestChange(user.ID, item.DayOfWeek, item.PairNumber, item.Subject, item.Classroom)
		if err != nil {
			msg.Text = "Не удалось отправить предложение"
		} else {
			msg.Text = "Предложение отправлено админу."
		}

	case strings.HasPrefix(text, "/edit") && user.Role == "admin":
		err, item := parseSuggestion(text)
		if err != nil {
			msg.Text = handleParseErrors(err)
			bot.Send(msg)
			return
		}
		err = db.EditSchedule(item.DayOfWeek, item.PairNumber, item.Subject, item.Classroom)
		if err != nil {
			msg.Text = fmt.Sprintln("Не удалось обновить расписание\n", err)
		} else {
			msg.Text = "Расписание обновлено."
		}

	case strings.HasPrefix(text, "/approve") && user.Role == "admin":
		// Format: /approve <sugg_id>
		parts := strings.SplitN(text, " ", 2)
		if len(parts) < 2 {
			// Show pending
			suggs, err := db.GetPendingSuggestions()
			if err != nil {
				msg.Text = "Ошибка."
			} else {
				var sb strings.Builder
				sb.WriteString("Pending предложения:\n")
				if len(suggs) == 0 {
					sb.WriteString("Нет предложений.\n")
				} else {
					msg.ReplyMarkup = createSuggestionsKeyboard(suggs, "approve")
				}
				msg.Text = sb.String()
			}
			bot.Send(msg)
			return
		}
		suggID, _ := strconv.Atoi(parts[1])
		err := db.ApproveSuggestion(suggID)
		if err != nil {
			msg.Text = fmt.Sprintln("Ошибка одобрения.\n", err)
		} else {
			msg.Text = "Предложение одобрено и применено."
		}

	case strings.HasPrefix(text, "/delete") && user.Role == "admin":
		// Format: /delete <day> <pair>
		parts := strings.SplitN(text, " ", 3)
		if len(parts) < 3 {
			msg.Text = "Неправильно введены данные.\nФормат: /delete <день> <пара>"
			bot.Send(msg)
			return
		}
		day, err := strconv.Atoi(parts[1])
		if err != nil || day < 1 || day > 7 {
			msg.Text = "День должен быть от 1 до 7."
			bot.Send(msg)
			return
		}
		pair, err := strconv.Atoi(parts[2])
		if err != nil || pair < 1 || pair > 8 {
			msg.Text = "Номер пары должен быть от 1 до 8."
			bot.Send(msg)
			return
		}
		err = db.DeleteSubject(day, pair)
		if err != nil {
			msg.Text = "Не получилось удалить предмет из расписания."
			bot.Send(msg)
			return
		}
		msg.Text = fmt.Sprintf("Предмет удален в %d дне %d парой.", day, pair)

	case strings.HasPrefix(text, "/clear_suggestions") && user.Role == "admin":
		// Format: /clearSuggestions <status_to_clear>
		parts := strings.SplitN(text, " ", 2)
		if len(parts) < 2 {
			db.ClearSuggestions("approved")
			db.ClearSuggestions("rejected")
			db.ClearSuggestions("pending")
			msg.Text = "Успешно почистилось"
			bot.Send(msg)
			return
		}
		err = db.ClearSuggestions(parts[1])
		if err != nil {
			msg.Text = fmt.Sprintln("Неправильный sql запрос\n", err)
		} else {
			msg.Text = "Успешно почистилось"
		}

	case strings.HasPrefix(text, "/reject") && user.Role == "admin":
		// Format: /reject <sugg_id>
		parts := strings.SplitN(text, " ", 2)
		if len(parts) < 2 {
			// Show pending
			suggs, err := db.GetPendingSuggestions()
			if err != nil {
				msg.Text = "Ошибка."
			} else {
				var sb strings.Builder
				sb.WriteString("Pending предложения:\n")
				if len(suggs) == 0 {
					sb.WriteString("Нет предложений.\n")
				} else {
					msg.ReplyMarkup = createSuggestionsKeyboard(suggs, "reject")
				}
				msg.Text = sb.String()
			}
			bot.Send(msg)
			return
		}
		suggID, _ := strconv.Atoi(parts[1])
		err := db.RejectSuggestion(suggID)
		if err != nil {
			msg.Text = "Ошибка отклонения."
		} else {
			msg.Text = "Предложение отклонено."
		}

	default:
		if user.Role == "user" {
			msg.Text = "Команды:\nПосмотреть расписание на определенный день недели /schedule\n\n" +
				"Посмотреть расписание на сегодня /today\n\n" +
				"Предложить старосте изменить расписание\n/suggest <день  пара  предмет(1 словом)  аудитория>\n"
		} else {
			msg.Text = "Команды:\n/schedule <день 1-7>\n\n/suggest <день 1-7  пара 1-8  премдет (1 словом)  аудитория>\n" +
				"/edit <день 1-7  пара 1-8  предмет  аудитория>\n/approve <id>\n/reject <id>\n/clear_suggestions <status_to_clear>\n" +
				"/delete <day  pair>\n"
		}

	}

	bot.Send(msg)
}

func sprintSchedule(day int, items []models.ScheduleItem) string {
	var sb strings.Builder
	days := []string{"", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
	sb.WriteString(fmt.Sprintf("Расписание на %s:\n", days[day]))
	if len(items) == 0 {
		sb.WriteString("Расписание пусто.\n")
	} else {
		pairTime := []string{"", "8:30 - 10:00", "10:10 - 11:40", "11:50 - 13:20", "13:50 - 15:20", "15:30 - 17:00", "17:10 - 18:40", "18:50 - 20:20", "ggwp"}
		var lastPairNumber int
		for _, item := range items {
			if item.PairNumber > lastPairNumber+1 {
				for i := lastPairNumber + 1; i < item.PairNumber; i++ {
					sb.WriteString(fmt.Sprintf("%d: Окно (%s) | Каб. %d |\n", i, pairTime[i], item.Classroom))
				}
			}
			sb.WriteString(fmt.Sprintf("%d: %s (%s) | Каб. %d |\n", item.PairNumber, item.Subject, pairTime[item.PairNumber], item.Classroom))
			lastPairNumber = item.PairNumber
		}
	}
	return sb.String()
}

func parseSuggestion(text string) (error, *models.ScheduleItem) {
	// Format: /edit <day> <pair> <new_subject> <classroom>
	parts := strings.SplitN(text, " ", 5)
	if len(parts) < 5 {
		return errors.New("Invalid format"), nil
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 7 {
		return errors.New("Invalid day"), nil
	}
	pair, err := strconv.Atoi(parts[2])
	if err != nil || pair < 1 || pair > 8 {
		return errors.New("Invalid pair number"), nil
	}
	newSub := utils.CapitalizeFirstLetter(parts[3])
	classroom, err := strconv.Atoi(parts[4])
	if err != nil {
		return errors.New("Invalid classroom"), nil
	}

	return nil, &models.ScheduleItem{
		ID:         -1,
		DayOfWeek:  day,
		PairNumber: pair,
		Classroom:  classroom,
		Subject:    newSub,
	}
}

func handleParseErrors(err error) string {
	switch err.Error() {
	case "Invalid format":
		return "Формат: /suggest <день(1-7)  пара  новый предмет  аудитория>"

	case "Invalid day":
		return "День должен быть от 1 до 7."

	case "Invalid pair number":
		return "Номер пары должен быть от 1 до 8."

	case "Invalid classroom":
		return "Введите валидный номер аудитории"

	default:
		return fmt.Sprintln("Ошибка предложения.\n", err)
	}
}
