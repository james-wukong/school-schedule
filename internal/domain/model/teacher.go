// Package model contains the core data models for the school scheduling system.
package model

type Teacher struct {
	ID       int
	TenantID int
	Name     string
	// AvailableHours maps day to available time slots (e.g., "Monday" -> ["9:00", "10:00"])
	AvailableHours   map[DayOfWeek][]string
	MaxClassesPerDay int
	MaxHoursPerWeek  int
	Preferences      map[DayOfWeek]int // "Monday" -> 5 means prefer Monday (higher is better)
}

// type Teacher struct {
// 	ID             string
// 	TenantID       string
// 	Name           string
// 	MaxHoursPerWeek int
// }

func (t Teacher) CanTakeMoreHours(current int) bool {
	return current < t.MaxHoursPerWeek
}
