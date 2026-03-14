// Package scheduling define the core types for the scheduling problem, including VariableID,
// Assignment, Variable, and Solution. These types are used throughout the scheduling engine
// and constraint evaluators to represent the scheduling problem and its solutions.
package scheduling

// VariableID A unique identifier for a scheduling variable
// and represents one class session that must be scheduled.
// example: Math class for Grade 10 taught by Teacher A (session #1):
// VariableID("course123-session1")
type VariableID string

// Assignment A specific placement for a class session, it represents
// When and where will this class happen
type Assignment struct {
	RoomID     string
	TimeSlotID string
}

// Variable represents one class session that must be scheduled.
// The domain is the list of all possible placements for that session.
type Variable struct {
	ID     VariableID
	Domain []Assignment
}

// Solution represents the final schedule produced by the solver.
type Solution map[VariableID]Assignment
