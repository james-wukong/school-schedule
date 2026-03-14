// Package constraint defines the HardConstraintEvaluator interface and its default
// implementation. The HardConstraint Evaluator is responsible for checking if a given
// variable assignment is valid based on the hard constraints of the scheduling problem,
// such as teacher availability, room capacity, and student group conflicts. The default
// implementation uses maps to efficiently check these constraints against the current
// solution.
package constraint

import "github.com/james-wukong/school-schedule/internal/domain/scheduling"

type HardConstraintEvaluator interface {
	IsValid(
		varID scheduling.VariableID,
		assignment scheduling.Assignment,
		current scheduling.Solution,
	) bool
}

type DefaultHardEvaluator struct {
	TeacherMap      map[scheduling.VariableID]string
	GroupMap        map[scheduling.VariableID]string
	RoomCapacityMap map[string]int
	// GroupSizeMap    map[scheduling.VariableID]int
}

func (h *DefaultHardEvaluator) IsValid(
	varID scheduling.VariableID,
	a scheduling.Assignment,
	current scheduling.Solution,
) bool {
	for otherVar, otherAssign := range current {
		// Same timeslot conflict checks
		if otherAssign.TimeSlotID == a.TimeSlotID {
			// Teacher conflict
			if h.TeacherMap[varID] == h.TeacherMap[otherVar] {
				return false
			}

			// Student group conflict
			if h.GroupMap[varID] == h.GroupMap[otherVar] {
				return false
			}

			// Room conflict
			if otherAssign.RoomID == a.RoomID {
				return false
			}
		}
	}

	// Room capacity check
	// if h.RoomCapacityMap[a.RoomID] < h.GroupSizeMap[varID] {
	// 	return false
	// }

	return true
}
