package models

import "github.com/robfig/cron/v3"

type User struct {
	ID         int    `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	Role       string `db:"role"`
	FirstName  string `db:"first_name,omitempty"`
	UserName   string `db:"username,omitempty"`
	Subgroup   int    `db:"subgroup"`
}

type ScheduleItem struct {
	ID         int    `db:"id"`
	DayOfWeek  int    `db:"day_of_week"`
	PairNumber int    `db:"pair_number"`
	Classroom  int    `db:"classroom"`
	Subject    string `db:"subject"`
	Subgroup   int    `db:"subgroup"`
}

type Suggestion struct {
	ID         int    `db:"id"`
	UserID     int    `db:"user_id"`
	DayOfWeek  int    `db:"day_of_week"`
	PairNumber int    `db:"pair_number"`
	Classroom  int    `db:"classroom"`
	NewSubject string `db:"new_subject"`
	Subgroup   int    `db:"subgroup"`
	Status     string `db:"status"`
}

type Teacher struct {
	ID         int    `db:"id"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	SecondName string `db:"second_name"`
	Subject    string `db:"subject"`
	Subgroup   int    `db:"subgroup"`
}

type ScheduledMessage struct {
	ID          int          `db:"id"`
	CronSpec    string       `db:"cron_spec"`
	CronEntryID cron.EntryID `db:"cron_entry_id"`
	Message     string       `db:"message"`
	Subgroup    int          `db:"subgroup"`
}
