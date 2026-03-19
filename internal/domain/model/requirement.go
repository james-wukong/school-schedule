package model

// Subject ────> "what topic"
//     │
//     └──> Requirement ────> "who teaches it, to whom, how often"
//             │
//             ├──> SchoolClass  ────> "which group of students"
//             ├──> Teacher      ────> "which teacher"
//             └── >SessionsPerWeek  ────> "how many times/week"

type RequirementID int

// Year 10A  ──┬── Requirement: Math    / Alice / 4x
//             ├── Requirement: English / Bob   / 4x
//             ├── Requirement: Science / Carol / 3x
//             ├── Requirement: History / Dave  / 2x
//             └── Requirement: PE      / Eve   / 2x

// Requirement Represents one row in the curriculum matrix
type Requirement struct {
	ID              RequirementID
	SchoolClass     *SchoolClass // WHO is being taught
	Subject         *Subject     // WHAT is being taught
	Teacher         *Teacher     // WHO is teaching
	SessionsPerWeek int          // HOW MANY times per week

	// Constraints
	MinDayGap     int         // minimum days between classes (e.g., 1 = no consecutive days)
	PreferredDays []DayOfWeek // preferred days ("Monday", "Tuesday", etc.)
}

func NewRequirement(entity Requirement) *Requirement {
	r := &entity

	if entity.Teacher != nil {
		r.Teacher = NewTeacher(*entity.Teacher)
	}
	if entity.Subject != nil {
		r.Subject = NewSubject(*entity.Subject)
	}
	if entity.SchoolClass != nil {
		r.SchoolClass = NewSchoolClass(*entity.SchoolClass)
	}
	if len(entity.PreferredDays) > 0 {
		copy(r.PreferredDays, entity.PreferredDays)
	}
	return r
}
