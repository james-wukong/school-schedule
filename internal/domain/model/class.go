package model

type Class struct {
	ID           int
	TenantID     int
	TeacherID    int
	StudentCount int
	Duration     int // minutes
	RequiredDays int // e.g., 3 means schedule 3 times per week
	Name         string
	// Constraints
	MinDayGap     int         // minimum days between classes (e.g., 1 = no consecutive days)
	PreferredDays []DayOfWeek // preferred days ("Monday", "Tuesday", etc.)

}
