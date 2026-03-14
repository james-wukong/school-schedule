package model

import "time"

type DayOfWeek int

const (
	Monday DayOfWeek = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// type TimeSlot struct {
// 	Day  string // "Monday", "Tuesday", etc.
// 	Time string // "9:00", "10:00", etc.
// }

type TimeSlot struct {
	StartTime string // "9:00AM", "1:00PM"
	EndTime   string
	Day       DayOfWeek
	ID        int
	TenantID  int
}

func (t TimeSlot) DurationMinutes() int {
	start, err := time.Parse("15:04", t.StartTime)
	if err != nil {
		return 0
	}
	end, err := time.Parse("15:04", t.EndTime)
	if err != nil {
		return 0
	}
	return int(end.Sub(start).Minutes())
}
