package db

import (
	"errors"
	"fmt"

	"github.com/Keith987654321/schedule-tg-bot/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect(user, host, port, pass, db, sslmode string) {
	dsn := fmt.Sprintf("user=%s host=%s port=%s password=%s dbname=%s sslmode=%s", user, host, port, pass, db, sslmode)
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

func AddUserToSubgroup(telegramID int64, subgroup int) error {
	_, err := DB.Exec("UPDATE users SET subgroup = $1 WHERE telegram_id = $2", subgroup, telegramID)
	return err
}

func ChangeSubgroup(telegramID int64, newSubgroup int) error {
	_, err := DB.Exec("UPDATE users SET subgroup = $1 where telegram_id = $2", newSubgroup, telegramID)
	return err
}

func GetScheduleForDay(day, subgroup int) ([]models.ScheduleItem, error) {
	var items []models.ScheduleItem
	err := DB.Select(&items, "SELECT * FROM schedule WHERE day_of_week = $1 AND (subgroup = $2 OR subgroup = 0) ORDER BY pair_number", day, subgroup)
	return items, err
}

func SuggestChange(userID int, day, pair int, newSubject string, classroom int, subgroup int) error {
	_, err := DB.Exec("INSERT INTO suggestions (user_id, day_of_week, pair_number, new_subject, classroom, subgroup) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, day, pair, newSubject, classroom, subgroup)
	return err
}

func GetTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := DB.Select(&teachers, "SELECT * FROM teachers ORDER BY id")
	return teachers, err
}

func GetPendingSuggestions() ([]models.Suggestion, error) {
	var suggs []models.Suggestion
	err := DB.Select(&suggs, "SELECT * FROM suggestions WHERE status = 'pending'")
	return suggs, err
}

func UpdateSchedule(sugg models.Suggestion) error {
	var err error
	var pairs []models.ScheduleItem
	if err := DB.Select(&pairs, "SELECT * FROM schedule WHERE day_of_week = $1 AND pair_number = $2 AND (subgroup = 0 OR subgroup = 1 OR subgroup = 2) ORDER BY subgroup",
		sugg.DayOfWeek, sugg.PairNumber); err != nil {
		return err
	}

	if len(pairs) == 0 {
		_, err = DB.Exec("INSERT INTO schedule (day_of_week, pair_number, subject, classroom, subgroup) values ($1, $2, $3, $4, $5)", sugg.DayOfWeek, sugg.PairNumber,
			sugg.NewSubject, sugg.Classroom, sugg.Subgroup)
	} else if len(pairs) >= 1 && sugg.Subgroup == pairs[0].Subgroup {
		_, err = DB.Exec("UPDATE schedule SET subject = $1, classroom = $2 WHERE day_of_week = $3 AND pair_number = $4 AND subgroup = $5",
			sugg.NewSubject, sugg.Classroom, sugg.DayOfWeek, sugg.PairNumber, sugg.Subgroup)
	} else if len(pairs) == 1 && sugg.Subgroup != pairs[0].Subgroup && pairs[0].Subgroup != 0 && sugg.Subgroup != 0 {
		_, err = DB.Exec("INSERT INTO schedule (day_of_week, pair_number, subject, classroom, subgroup) values ($1, $2, $3, $4, $5)", sugg.DayOfWeek, sugg.PairNumber,
			sugg.NewSubject, sugg.Classroom, sugg.Subgroup)
	} else {
		return errors.New("Already exist at this pair time")
	}

	return err
}

func ApproveSuggestion(suggID int) error {
	var sugg models.Suggestion
	err := DB.Get(&sugg, "SELECT * FROM suggestions WHERE id = $1", suggID)
	if err != nil {
		return err
	}

	err = UpdateSchedule(sugg)

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

func EditSchedule(day, pair int, newSubject string, classroom, subgroup int) error {
	var pairs []models.ScheduleItem
	if err := DB.Select(&pairs, "SELECT * FROM schedule WHERE day_of_week = $1 AND pair_number = $2 AND (subgroup = $3 OR subgroup = 0 OR subgroup = 1 OR subgroup = 2)",
		day, pair, subgroup); err != nil {
		return err
	}

	err := UpdateSchedule(models.Suggestion{
		DayOfWeek:  day,
		PairNumber: pair,
		NewSubject: newSubject,
		Classroom:  classroom,
		Subgroup:   subgroup,
	})

	return err
}

func DeleteSubject(day, pair, subgroup int) error {
	_, err := DB.Exec("DELETE FROM schedule WHERE day_of_week = $1 AND pair_number = $2 AND subgroup = $3", day, pair, subgroup)
	return err
}

func ClearSuggestions(status string) error {
	_, err := DB.Exec("DELETE FROM suggestions WHERE status = $1", status)
	return err
}
