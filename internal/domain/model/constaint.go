package model

type ConstraintType string

const (
	ConstraintTeacherAvailability ConstraintType = "teacher_availability"
	ConstraintRoomCapacity        ConstraintType = "room_capacity"
	ConstraintNoOverlap           ConstraintType = "no_overlap"
	ConstraintPreferredTime       ConstraintType = "preferred_time"
)

type Constraint struct {
	ID       int
	TenantID int
	// Identifies which object the constraint applies to. For example,
	// for teacher availability, this would be the teacher ID;
	// for room capacity, this would be the classroom ID.
	EntityID int
	Type     ConstraintType
	Value    string
	IsHard   bool
}

type ConstraintManager struct {
	ID             int
	TenantID       int
	HardViolations int
	SoftViolations int
	MaxIterations  int
	Temperature    float64 // for simulated annealing
	CoolingRate    float64
}

func (c Constraint) IsSoft() bool {
	return !c.IsHard
}
