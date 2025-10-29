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
	keyboards "github.com/Keith987654321/schedule-tg-bot/pkg"
	"github.com/Keith987654321/schedule-tg-bot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
		if user.Subgroup == 0 {
			msg.Text = "Выберите свою подгруппу:"
			msg.ReplyMarkup = keyboards.CreateSubgroupSelectionKeyboard()
			bot.Send(msg)
			return
		}
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
		if strings.HasPrefix(callbackData, "subgroup_") {
			subgroup, err := strconv.Atoi(strings.TrimPrefix(callbackData, "subgroup_"))
			if err != nil {
				log.Println("Can't convert subgroup from string to int")
				msg.Text = "Не удалось вас присвоить к одной из подгрупп."
				bot.Send(msg)
				return
			}
			db.AddUserToSubgroup(user.TelegramID, subgroup)
			msg.Text = fmt.Sprintf("Вы были добавлены в подгруппу %d", subgroup)
			msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
			bot.Send(msg)
			// Deleting inline keayboard
			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			bot.Send(editMsg)
			return
		} else if strings.HasPrefix(callbackData, "schedule_") {
			day, _ := strconv.Atoi(strings.TrimPrefix(callbackData, "schedule_"))
			items, err := db.GetScheduleForDay(day, user.Subgroup)
			if err != nil {
				msg.Text = "Не получилось получить расписание."
			} else {
				msg.Text = sprintSchedule(day, items)
			}
			msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
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
			if err != nil && err.Error() == "Already exist at this pair time" {
				msg.Text = "Уже есть пара запланированная на это время."
			} else if err != nil {
				msg.Text = "Ошибка одобрения."
			} else {
				msg.Text = "Предложение одобрено и применено."
			}

			msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
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
			msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
			bot.Send(msg)

			editMsg := tgbotapi.NewEditMessageReplyMarkup(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			bot.Send(editMsg)
			return
		} else if strings.HasPrefix(callbackData, "change_group_") {
			newSubgroup, err := strconv.Atoi(strings.TrimPrefix(callbackData, "change_group_"))
			if err != nil {
				msg.Text = "Неправильный номер подгруппы."
				msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
				bot.Send(msg)
				return
			}
			err = db.ChangeSubgroup(user.TelegramID, newSubgroup)
			if err != nil {
				msg.Text = "Не удалось поменять подгруппу."
			} else {
				msg.Text = fmt.Sprintf("Подгруппа успешно изменена на %d.", newSubgroup)
			}
			msg.ReplyMarkup = keyboards.CreateMainMenu(user.Role)
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
	for ruText, enText := range keyboards.Commands {
		if text == strings.ToLower(ruText) {
			text = enText
			break
		}
	}

	today := int(time.Now().Weekday()) // 0=Sunday, 1=Monday
	if today == 0 {
		today = 7
	}

	switch {
	case strings.HasPrefix(text, "/today"):
		// Format /today
		items, err := db.GetScheduleForDay(today, user.Subgroup)
		if err != nil {
			msg.Text = "Не получилось получить расписание."
		} else {
			msg.Text = sprintSchedule(today, items)
		}

	case strings.HasPrefix(text, "/schedule"):
		// Format: /schedule [day]
		parts := strings.SplitN(text, " ", 2)
		if len(parts) == 2 {
			parsedDay, err := strconv.Atoi(parts[1])
			if err != nil || parsedDay < 1 || parsedDay > 7 {
				msg.Text = "Укажите день недели (1-7, где 1 — понедельник, 7 — воскресенье)."
				bot.Send(msg)
				return
			}
			day := parsedDay
			items, err := db.GetScheduleForDay(day, user.Subgroup)
			if err != nil {
				msg.Text = "Не получилось получить расписание."
			} else {
				msg.Text = sprintSchedule(day, items)
			}
		} else {
			msg.Text = "Выберите день недели:"
			msg.ReplyMarkup = keyboards.CreateDaySelectionKeyboard()
		}

	case strings.HasPrefix(text, "/tomorrow"):
		// Format /tomorrow
		day := today + 1
		if day > 7 {
			day = 1
		}
		items, err := db.GetScheduleForDay(day, user.Subgroup)
		if err != nil {
			msg.Text = "Не получилось получить расписание."
		} else {
			msg.Text = sprintSchedule(day, items)
		}

	case strings.HasPrefix(text, "/change_group"):
		// Format: /change_group <id>
		parts := strings.SplitN(strings.TrimPrefix(text, "/change_group"), " ", 2)
		if len(parts) == 2 {
			newSubgroup, err := strconv.Atoi(parts[1])
			if err != nil {
				msg.Text = "Неправильный номер подгруппы."
			} else {
				err := db.ChangeSubgroup(user.TelegramID, newSubgroup)
				if err != nil {
					msg.Text = "Не удалось поменять подгруппу."
				} else {
					msg.Text = fmt.Sprintf("Подгруппа успешно изменена на %d.", user.Subgroup)
				}
			}
		} else {
			msg.Text = "Выберите номер подгруппы:"
			msg.ReplyMarkup = keyboards.CreateChangeSubgroupSelectionKeyboard()
		}

	case strings.HasPrefix(text, "/teachers"):
		teachers, err := db.GetTeachers()
		if err != nil {
			msg.Text = "Не удалось получить список преподавателей."
		} else {
			msg.Text = sprintTeachers(teachers)
		}

	case strings.HasPrefix(text, "/suggest") && (user.Role == "user" || user.Role == "admin"):
		item, err := parseSuggestion(text)
		if err != nil {
			msg.Text = handleParseErrors(err, "suggest")
			bot.Send(msg)
			return
		}
		err = db.SuggestChange(user.ID, item.DayOfWeek, item.PairNumber, item.Subject, item.Classroom, item.Subgroup)
		if err != nil {
			log.Println(err)
			msg.Text = "Не удалось отправить предложение"
		} else {
			msg.Text = "Предложение отправлено старосте."
		}

	case strings.HasPrefix(text, "/edit") && user.Role == "admin":
		item, err := parseSuggestion(text)
		if err != nil {
			msg.Text = handleParseErrors(err, "edit")
			bot.Send(msg)
			return
		}
		err = db.EditSchedule(item.DayOfWeek, item.PairNumber, item.Subject, item.Classroom, item.Subgroup)
		if err != nil && err.Error() == "Already exist at this pair time" {
			msg.Text = "Уже есть пара запланированная на это время"
		} else if err != nil {
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
					msg.ReplyMarkup = keyboards.CreateSuggestionsKeyboard(suggs, "approve")
				}
				msg.Text = sb.String()
			}
			bot.Send(msg)
			return
		}
		suggID, _ := strconv.Atoi(parts[1])
		err := db.ApproveSuggestion(suggID)
		if err.Error() == "Already exist at this pair time" {
			msg.Text = "Уже запланирована пара на это время"
		} else if err != nil {
			msg.Text = fmt.Sprintln("Ошибка одобрения.\n", err)
		} else {
			msg.Text = "Предложение одобрено и применено."
		}

	case strings.HasPrefix(text, "/delete") && user.Role == "admin":
		// Format: /delete <day> <pair> <subgroup>
		var subgroup int
		var err error
		parts := strings.SplitN(text, " ", 4)

		if len(parts) == 4 {
			subgroup, err = strconv.Atoi(parts[3])
			if err != nil {
				msg.Text = "Неподходящий номер подгруппы."
				bot.Send(msg)
				return
			}
		}
		if len(parts) < 3 {
			msg.Text = "Неправильно введены данные.\nФормат: /delete <день> <пара> <подгруппа>"
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
		err = db.DeleteSubject(day, pair, subgroup)
		if err != nil {
			msg.Text = "Не получилось удалить предмет из расписания."
			bot.Send(msg)
			return
		}
		msg.Text = fmt.Sprintf("Предмет удален в %d дне %d парой для подгруппы %d.", day, pair, subgroup)

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
					msg.ReplyMarkup = keyboards.CreateSuggestionsKeyboard(suggs, "reject")
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

	case strings.HasPrefix(text, "/schedmes") && user.Role == "admin":
		// Format: /schedmes \<message>\<day>\<pair number>\<subgroup>
		var err error
		var message string
		var schedDay, pairNumber, subgroup int
		parts := strings.SplitN(text, "\\", 5)
		if len(parts) >= 4 {
			message = parts[1]
			if schedDay, err = strconv.Atoi(parts[2]); err != nil || schedDay < 1 || schedDay > 7 {
				msg.Text = "Номер дня должен быть от 1 до 7"
				bot.Send(msg)
				return
			}
			if pairNumber, err = strconv.Atoi(parts[3]); err != nil || pairNumber < 1 || pairNumber > 7 {
				msg.Text = "Номер пары должен быть от 1 до 7"
				bot.Send(msg)
				return
			}
			if len(parts) == 5 {
				subgroup, err = strconv.Atoi(parts[4])
				if err != nil || !(keyboards.Subgroups[0] <= subgroup && keyboards.Subgroups[len(keyboards.Subgroups)-1] <= subgroup) {
					msg.Text = "Введена некорректная подгруппа"
					bot.Send(msg)
					return
				}
			}
			crocSpec := convertTimeToCronSpec(schedDay, pairNumber)
			if err = addScheduledMessage(bot, crocSpec, message, subgroup); err != nil {
				msg.Text = fmt.Sprintf("Не удалось запланировать отложенное сообщение: %v", err)
			} else {
				msg.Text = "Отложенное сообщение успешно добавлено"
			}
			bot.Send(msg)
			return

		} else {
			msg.Text = "Неправильно введённые данные\nFormat: /schedmes \\<message>\\<day>\\<pair number>\\<subgroup=0>"
			bot.Send(msg)
			return
		}

	case strings.HasPrefix(text, "/showschedmes") && user.Role == "admin":
		items, err := db.GetScheduledMessages()
		if err != nil {
			msg.Text = fmt.Sprintf("Can't get scheduled messages: %v\n", err)
			bot.Send(msg)
			return
		}

		msg.Text = sprintSchedmes(items)
		bot.Send(msg)
		return
	case strings.HasPrefix(text, "/delschedmes_") && user.Role == "admin":
		schedmesID, err := strconv.Atoi(strings.TrimPrefix(text, "/delschedmes_"))
		if err != nil {
			msg.Text = "Invalid scheduled message id"
			bot.Send(msg)
			return
		}
		err = deleteScheduledMessage(schedmesID)
		if err != nil {
			msg.Text = fmt.Sprint(err)
		} else {
			msg.Text = fmt.Sprintf("Сообщение с id: %d успешно удаленно", schedmesID)
		}

	default:
		if user.Role == "user" {
			msg.Text = "Команды:\nПосмотреть расписание на определенный день недели /schedule\n\n" +
				"Посмотреть расписание на сегодня /today\n\n" +
				"Поменять подгруппу /change_group\n\nПосмотреть список преподавателей /teachers\n\n" +
				"Предложить старосте изменить расписание\n/suggest день(1-7)  пара  предмет(1 словом)  аудитория  подгруппа(если пара только для одной из подгрупп)\n"
		} else {
			msg.Text = "Команды:\n/schedule <день 1-7>\n\n/suggest <день 1-7  пара 1-8  премдет (1 словом)  аудитория  подгруппа>\n" +
				"Поменять подгруппу /change_group\n\nПосмотреть список преподавателей /teachers\n\n" +
				"/edit <день 1-7  пара 1-8  предмет  аудитория  подгруппа>\n/approve <id>\n/reject <id>\n/clear_suggestions <status_to_clear>\n" +
				"/delete <day  pair subgroup>\n/schedmes \\<message>\\<day>\\<pair number>\\<subgroup=0>\n" +
				"Посмотреть все запланированные сообщения /showschedmes\nУдалить запланированное сообщение /delschedmes_<id>\n"
		}

	}

	bot.Send(msg)
}

func sprintSchedule(day int, items []models.ScheduleItem) string {
	var sb strings.Builder
	days := []string{"", "Понедельник", "Вторник", "Среду", "Четверг", "Пятницу", "Субботу", "Воскресенье"}
	sb.WriteString(fmt.Sprintf("Расписание на %s:\n", days[day]))
	if len(items) == 0 {
		sb.WriteString("Расписание пусто.\n")
	} else {
		pairTime := []string{"", "8:30 - 10:00", "10:10 - 11:40", "11:50 - 13:20", "13:50 - 15:20", "15:30 - 17:00", "17:10 - 18:40", "18:50 - 20:20", "ggwp"}
		var lastPairNumber int
		for _, item := range items {
			if item.PairNumber > lastPairNumber+1 {
				for i := lastPairNumber + 1; i < item.PairNumber; i++ {
					sb.WriteString(fmt.Sprintf("%d: Окно (%s)\n", i, pairTime[i]))
				}
			}
			if item.Classroom != 0 {
				sb.WriteString(fmt.Sprintf("%d: %s (%s) | Каб. %d\n", item.PairNumber, item.Subject, pairTime[item.PairNumber], item.Classroom))
			} else {
				sb.WriteString(fmt.Sprintf("%d: %s (%s)\n", item.PairNumber, item.Subject, pairTime[item.PairNumber]))
			}

			lastPairNumber = item.PairNumber
		}
	}
	return sb.String()
}

func sprintSchedmes(items []models.ScheduledMessage) string {
	if len(items) == 0 {
		return "Нет запланированных сообщений"
	}
	var sb strings.Builder
	sb.WriteString("Запланированные сообщения:\n")
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("%d) msg: %s\nsubgroup: %d\nId: %d\ncron spec: %s\nentryID: %v\n\n", i+1, item.Message, item.Subgroup, item.ID, item.CronSpec, item.CronEntryID))
	}
	return sb.String()
}

func sprintTeachers(teachers []models.Teacher) string {
	var sb strings.Builder
	for i, teacher := range teachers {
		sb.WriteString(fmt.Sprintf("%d: ", i+1))
		if teacher.SecondName != "-" {
			sb.WriteString(fmt.Sprintf("%s ", teacher.SecondName))
		}
		if teacher.FirstName != "-" {
			sb.WriteString(fmt.Sprintf("%s ", teacher.FirstName))
		}
		if teacher.MiddleName != "-" {
			sb.WriteString(fmt.Sprintf("%s ", teacher.MiddleName))
		}
		sb.WriteString(fmt.Sprintf("| %s ", teacher.Subject))
		if teacher.Subgroup != 0 {
			sb.WriteString(fmt.Sprintf("| подгруппа %d", teacher.Subgroup))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func parseSuggestion(text string) (*models.ScheduleItem, error) {
	// Format: /edit <day> <pair> <new_subject> <classroom> <subgroup>
	parts := strings.SplitN(text, " ", 6)
	var subgroup int
	var err error

	if len(parts) == 6 {
		subgroup, err = strconv.Atoi(parts[5])
		if err != nil {
			return nil, errors.New("invalid subgroup")
		}
		fmt.Printf("subgroup == %d\n\n", subgroup)
	}
	if len(parts) < 5 {
		return nil, errors.New("invalid format")
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 7 {
		return nil, errors.New("invalid day")
	}
	pair, err := strconv.Atoi(parts[2])
	if err != nil || pair < 1 || pair > 8 {
		return nil, errors.New("invalid pair number")
	}
	newSub := utils.CapitalizeFirstLetter(parts[3])
	classroom, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, errors.New("invalid classroom")
	}

	return &models.ScheduleItem{
		ID:         -1,
		DayOfWeek:  day,
		PairNumber: pair,
		Classroom:  classroom,
		Subject:    newSub,
		Subgroup:   subgroup,
	}, nil
}

func handleParseErrors(err error, funcName string) string {
	switch err.Error() {
	case "invalid format":
		return fmt.Sprintf("Формат: /%s день(1-7)  пара  новый_предмет(одним словом)  аудитория  подгруппа(если для конкретной подгруппы)", funcName)

	case "invalid day":
		return "День должен быть от 1 до 7."

	case "invalid pair number":
		return "Номер пары должен быть от 1 до 8."

	case "invalid classroom":
		return "Введите валидный номер аудитории."

	case "invalid subgroup":
		return "Введите правильный номер подгруппы."

	default:
		return fmt.Sprintln("Ошибка предложения.\n", err)
	}
}
