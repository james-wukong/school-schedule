package model

import (
	"slices"
)

type DayOfWeek int
type TimeSlotID int64

const (
	Monday DayOfWeek = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type TimeSlot struct {
	ID        TimeSlotID
	StartTime string // "09:00", "13:00"
	Day       DayOfWeek
}

// AllTimeSlots returns 30 time slots: 5 days/week * 6 slots/day
func AllTimeSlots() []TimeSlot {
	slots := make([]TimeSlot, 0, 30)
	for _, day := range []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday} {
		for _, slot := range []string{"09:00", "10:00", "11:00", "13:00", "14:00", "15:00"} {
			slots = append(slots, TimeSlot{Day: day, StartTime: slot})
		}
	}
	return slots
}

func AvailableTimeSlots(days []DayOfWeek, slots []string) map[DayOfWeek][]string {
	if len(days) == 0 || len(slots) == 0 {
		return make(map[DayOfWeek][]string, 0)
	}
	s := make(map[DayOfWeek][]string)
	for _, day := range days {
		s[day] = slots
	}
	return s
}

func RemoveElement[T DayOfWeek | string](elements []T, element T) []T {
	if len(elements) == 0 {
		return make([]T, 0)
	}
	for k, v := range elements {
		if v == element {
			elements = slices.Delete(elements, k, k+1)
		}
	}
	return elements
}
