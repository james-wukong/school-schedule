package model

// Subject ────> "what topic"
//     │
//     └──> Requirement ────> "who teaches it, to whom, how often"
//             │
//             ├──> SchoolClass  ────> "which group of students"
//             ├──> Teacher      ────> "which teacher"
//             └── >SessionsPerWeek  ────> "how many times/week"

// Requirement:  Year 10A | Math | Alice | 4x/week
//                  ↓ scheduler solves ↓
// Assignment 1: Year 10A | Math | Alice | Room-101 | Mon-9:00

// Assignment represents one concrete, scheduled event
type Assignment struct {
	Requirement *Requirement
	Room        *Room
	Slot        TimeSlot
}

func NewAssignment(entity Assignment) *Assignment {
	a := &entity

	if entity.Requirement != nil {
		a.Requirement = NewRequirement(*entity.Requirement)
	}
	if entity.Room != nil {
		a.Room = NewRoom(*entity.Room)
	}
	return a
}
