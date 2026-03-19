package model

type SubjectID int

// Subject ────> "what topic"
//     │
//     └──> Requirement ────> "who teaches it, to whom, how often"
//             │
//             ├──> SchoolClass  ────> "which group of students"
//             ├──> Teacher      ────> "which teacher"
//             └── >SessionsPerWeek  ────> "how many times/week"

type Subject struct {
	ID          SubjectID
	Name        string // "Math"
	RequiresLab bool
	IsHeavy     bool
}

func NewSubject(entity Subject) *Subject {
	return &entity
}
