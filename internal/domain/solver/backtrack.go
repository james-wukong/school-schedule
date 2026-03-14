package solver

import (
	"fmt"
	"math"
)

// backtrackingSearch uses MRV heuristic to find initial solution
func (s *Scheduler) backtrackingSearch() error {
	return s.backtrack(0)
}

func (s *Scheduler) backtrack(depth int) error {
	// Get unassigned class with MRV (fewest remaining values)
	var selectedClassID int
	minDomain := math.MaxInt

	for classID, domain := range s.Domains {
		if _, assigned := s.Assignments[classID]; !assigned {
			if len(domain) < minDomain {
				minDomain = len(domain)
				selectedClassID = classID
			}
		}
	}

	if selectedClassID == 0 {
		// All classes assigned - solution found
		return nil
	}

	// Try each value in domain
	for _, assignment := range s.Domains[selectedClassID] {
		if s.isConsistent(selectedClassID, assignment) {
			s.Assignments[selectedClassID] = assignment

			if err := s.backtrack(depth + 1); err == nil {
				return nil
			}

			delete(s.Assignments, selectedClassID)
		}
	}

	return fmt.Errorf("no solution found at depth %d", depth)
}

// isConsistent checks if assignment violates hard constraints
func (s *Scheduler) isConsistent(classID int, assignment *ScheduleEntry) bool {
	class := s.Classes[classID]

	// Check day gap constraint
	for assignedID, entry := range s.Assignments {
		if assignedID == classID {
			continue
		}

		if entry.Teacher.ID == assignment.Teacher.ID {
			dayDiff := dayDifference(entry.TimeSlot.Day, assignment.TimeSlot.Day)
			if dayDiff < class.MinDayGap && dayDiff > 0 {
				return false
			}
		}
	}

	// Check time slot conflicts
	for _, entry := range s.Assignments {
		if entry.Class.ID == assignment.Class.ID {
			continue
		}

		// Teacher conflict
		if entry.Teacher.ID == assignment.Teacher.ID &&
			entry.TimeSlot.Day == assignment.TimeSlot.Day &&
			entry.TimeSlot.StartTime == assignment.TimeSlot.StartTime {
			return false
		}

		// Room conflict
		if entry.Classroom.ID == assignment.Classroom.ID &&
			entry.TimeSlot.Day == assignment.TimeSlot.Day &&
			entry.TimeSlot.StartTime == assignment.TimeSlot.StartTime {
			return false
		}
	}

	return true
}
