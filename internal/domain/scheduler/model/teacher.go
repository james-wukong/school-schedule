// Package model contains the core data models for the school scheduling system.
package model

// Subject ────> "what topic"
//     │
//     └──> Requirement ────> "who teaches it, to whom, how often"
//             │
//             ├──> SchoolClass  ────> "which group of students"
//             ├──> Teacher      ────> "which teacher"
//             └── >SessionsPerWeek  ────> "how many times/week"

import (
	"maps"
)

type TeacherID int

type Teacher struct {
	ID       TeacherID
	TenantID int
	Name     string
	Subjects []SubjectID
	// AvailableTimes maps day to available time slots (e.g., "Monday" -> ["9:00", "10:00"])
	AvailableTimes   map[DayOfWeek][]string
	MaxClassesPerDay int
	MaxHoursPerWeek  int
	Preferences      map[DayOfWeek]int // "Monday" -> 5 means prefer Monday (higher is better)
}

type TeacherDaySlotsKey struct {
	TeacherID TeacherID
	Day       DayOfWeek
}

func NewTeacher(entity Teacher) *Teacher {
	t := &Teacher{
		ID:               entity.ID,
		TenantID:         entity.TenantID,
		Name:             entity.Name,
		Subjects:         make([]SubjectID, len(entity.Subjects)),
		AvailableTimes:   make(map[DayOfWeek][]string),
		MaxClassesPerDay: entity.MaxClassesPerDay,
		MaxHoursPerWeek:  entity.MaxHoursPerWeek,
		Preferences:      make(map[DayOfWeek]int),
	}
	if entity.Subjects != nil {
		copy(t.Subjects, entity.Subjects)
	}
	if len(entity.AvailableTimes) > 0 {
		maps.Copy(t.AvailableTimes, entity.AvailableTimes)
	}
	if len(entity.Preferences) > 0 {
		maps.Copy(t.Preferences, entity.Preferences)
	}
	return t
}

func (t Teacher) CanTakeMoreHours(current int) bool {
	return current < t.MaxHoursPerWeek
}
