package models

type User struct {
	ID         int    `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	Role       string `db:"role"`
	FirstName  string `db:"first_name,omitempty"`
	UserName   string `db:"username,omitempty"`
}

type ScheduleItem struct {
	ID         int    `db:"id"`
	DayOfWeek  int    `db:"day_of_week"`
	PairNumber int    `db:"pair_number"`
	Classroom  int    `db:"classroom"`
	Subject    string `db:"subject"`
}

type Suggestion struct {
	ID         int    `db:"id"`
	UserID     int    `db:"user_id"`
	DayOfWeek  int    `db:"day_of_week"`
	PairNumber int    `db:"pair_number"`
	Classroom  int    `db:"classroom"`
	NewSubject string `db:"new_subject"`
	Status     string `db:"status"`
}
