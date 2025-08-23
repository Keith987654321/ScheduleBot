package db

import (
	"fmt"

	"github.com/Keith987654321/schedule-tg-bot/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect(user, pass, db, sslmode string) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, pass, db, sslmode)
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("DB connection error: %v", err))
	}
}

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	user := &models.User{}
	err := DB.Get(user, "SELECT * FROM users WHERE telegram_id = $1", telegramID)
	if err != nil {
		// If there is no user we create a new one
		_, err = DB.Exec("INSERT INTO users (telegram_id, role) VALUES ($1, 'user')", telegramID)
		if err != nil {
			return nil, err
		}
		return GetUserByTelegramID(telegramID) // Getting user by recursion
	}
	return user, nil
}

func CheckUserInfo(telegramID int64, firstName, userName string) error {
	user := &models.User{}
	if err := DB.Get(user, "SELECT * FROM users WHERE telegram_id = $1", telegramID); err != nil {
		return err
	}
	if user.FirstName != firstName || user.UserName != userName {
		_, err := DB.Exec("UPDATE users SET first_name = $1, username = $2 WHERE telegram_id = $3", firstName, userName, telegramID)
		return err
	}
	return nil
}

func GetScheduleForDay(day int) ([]models.ScheduleItem, error) {
	var items []models.ScheduleItem
	err := DB.Select(&items, "SELECT * FROM schedule WHERE day_of_week = $1 ORDER BY pair_number", day)
	return items, err
}

func GetLastPairNumber(pairs []models.ScheduleItem) int {
	maxNumber := 0
	for _, pair := range pairs {
		if pair.PairNumber > maxNumber {
			maxNumber = pair.PairNumber
		}
	}
	return maxNumber
}

func SuggestChange(userID int, day, pair int, newSubject string, classroom int) error {
	_, err := DB.Exec("INSERT INTO suggestions (user_id, day_of_week, pair_number, new_subject, classroom) VALUES ($1, $2, $3, $4, $5)", userID, day, pair, newSubject, classroom)
	return err
}

func GetPendingSuggestions() ([]models.Suggestion, error) {
	var suggs []models.Suggestion
	err := DB.Select(&suggs, "SELECT * FROM suggestions WHERE status = 'pending'")
	return suggs, err
}

func ApproveSuggestion(suggID int) error {
	var sugg models.Suggestion
	err := DB.Get(&sugg, "SELECT * FROM suggestions WHERE id = $1", suggID)
	if err != nil {
		return err
	}
	// Updating schedule
	var pairs []models.ScheduleItem
	if err := DB.Select(&pairs, "SELECT * FROM schedule WHERE day_of_week = $1 AND pair_number = $2", sugg.DayOfWeek, sugg.PairNumber); err != nil {
		return err
	}

	if len(pairs) == 0 {
		_, err = DB.Exec("INSERT INTO schedule (day_of_week, pair_number, subject, classroom) values ($1, $2, $3, $4)", sugg.DayOfWeek, sugg.PairNumber,
			sugg.NewSubject, sugg.Classroom)
	} else {
		_, err = DB.Exec("UPDATE schedule SET subject = $1 WHERE day_of_week = $2 AND pair_number = $3", sugg.NewSubject, sugg.DayOfWeek, sugg.PairNumber)
	}

	if err != nil {
		return err
	}

	// Cleaning schedule
	var count int
	DB.Get(&count, "SELECT COUNT(*) FROM suggestions WHERE status = 'approved'")
	if count > 10 {
		ClearSuggestions("approved")
	}

	// Updating status
	_, err = DB.Exec("UPDATE suggestions SET status = 'approved' WHERE id = $1", sugg.ID)
	return err
}

func RejectSuggestion(suggID int) error {
	var count int
	DB.Get(&count, "SELECT COUNT(*) FROM suggestions WHERE status = 'rejected'")
	if count > 10 {
		ClearSuggestions("rejected")
	}

	_, err := DB.Exec("UPDATE suggestions SET status = 'rejected' WHERE id = $1", suggID)
	return err
}

func EditSchedule(day, pair int, newSubject string, classroom int) error {
	var pairs []models.ScheduleItem
	if err := DB.Select(&pairs, "SELECT * FROM schedule WHERE day_of_week = $1 AND pair_number = $2", day, pair); err != nil {
		return err
	}

	if len(pairs) == 0 {
		_, err := DB.Exec("INSERT INTO schedule (day_of_week, pair_number, subject, classroom) values ($1, $2, $3, $4)", day, pair, newSubject, classroom)
		return err
	}

	_, err := DB.Exec("UPDATE schedule SET subject = $1 WHERE day_of_week = $2 AND pair_number = $3", newSubject, day, pair)
	return err
}

func DeleteSubject(day, pair int) error {
	_, err := DB.Exec("DELETE FROM schedule WHERE day_of_week = $1 AND pair_number = $2", day, pair)
	return err
}

func ClearSuggestions(status string) error {
	_, err := DB.Exec("DELETE FROM suggestions WHERE status = $1", status)
	return err
}
