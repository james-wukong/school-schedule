package model

// Subject ────> "what topic"
//     │
//     └──> Requirement ────> "who teaches it, to whom, how often"
//             │
//             ├──> SchoolClass  ────> "which group of students"
//             ├──> Teacher      ────> "which teacher"
//             └── >SessionsPerWeek  ────> "how many times/week"

type ClassID int

// SchoolClass Constraint: Teacher can't teach 2 classes at same time
// Audience
type SchoolClass struct {
	ID           ClassID
	TenantID     int
	StudentCount int
	Duration     int    // minutes
	Grade        string // "Year-10"
	Class        string // "A"
}

type ClassDaySubjectKey struct {
	ClassID ClassID
	Day     DayOfWeek
}

func NewSchoolClass(entity SchoolClass) *SchoolClass {
	s := &entity

	if entity.Duration == 0 {
		s.Duration = 45
	}
	return s
}
