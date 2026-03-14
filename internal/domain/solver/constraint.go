package solver

import (
	"fmt"
)

type pair struct {
	x, y int
}

// constraintPropagation uses AC-3 algorithm to reduce domains
func (s *Scheduler) constraintPropagation() error {
	queue := make([]pair, 0)

	// Build initial queue
	classIDs := make([]int, 0, len(s.Classes))
	for id := range s.Classes {
		classIDs = append(classIDs, id)
	}

	for i, id1 := range classIDs {
		for _, id2 := range classIDs[i+1:] {
			queue = append(queue, pair{id1, id2})
			queue = append(queue, pair{id2, id1})
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if s.revise(current.x, current.y) {
			if len(s.Domains[current.x]) == 0 {
				return fmt.Errorf("domain wipeout for class %s", current.x)
			}

			// Add neighbors to queue
			for _, neighbor := range classIDs {
				if neighbor != current.x && neighbor != current.y {
					queue = append(queue, pair{neighbor, current.x})
				}
			}
		}
	}
	return nil
}

// revise removes inconsistent values from domain of x
func (s *Scheduler) revise(x, y int) bool {
	revised := false
	xDomain := s.Domains[x]
	yDomain := s.Domains[y]

	i := 0
	for i < len(xDomain) {
		xAssignment := xDomain[i]

		// Check if there's a compatible value in y's domain
		compatible := false
		for _, yAssignment := range yDomain {
			if s.areCompatible(xAssignment, yAssignment) {
				compatible = true
				break
			}
		}

		if !compatible {
			xDomain = append(xDomain[:i], xDomain[i+1:]...)
			revised = true
		} else {
			i++
		}
	}
	s.Domains[x] = xDomain
	return revised
}

// areCompatible checks if two assignments conflict
func (s *Scheduler) areCompatible(entry1, entry2 *ScheduleEntry) bool {
	// Same class - not compatible
	if entry1.Class.ID == entry2.Class.ID {
		return true
	}

	// Teacher conflict - same teacher at same time
	if entry1.Teacher.ID == entry2.Teacher.ID &&
		entry1.TimeSlot.Day == entry2.TimeSlot.Day &&
		entry1.TimeSlot.StartTime == entry2.TimeSlot.StartTime {
		return false
	}

	// Room conflict - same room at same time
	if entry1.Classroom.ID == entry2.Classroom.ID &&
		entry1.TimeSlot.Day == entry2.TimeSlot.Day &&
		entry1.TimeSlot.StartTime == entry2.TimeSlot.StartTime {
		return false
	}

	return true
}

// countHardViolations counts the number of hard constraint violations in the current
// assignments
func (s *Scheduler) countHardViolations() int {
	violations := 0
	for _, entry := range s.Assignments {
		for otherID, otherEntry := range s.Assignments {
			if otherID == entry.Class.ID {
				continue
			}
			if entry.Teacher.ID == otherEntry.Teacher.ID &&
				entry.TimeSlot.Day == otherEntry.TimeSlot.Day &&
				entry.TimeSlot.StartTime == otherEntry.TimeSlot.StartTime {
				violations++
			}
			if entry.Classroom.ID == otherEntry.Classroom.ID &&
				entry.TimeSlot.Day == otherEntry.TimeSlot.Day &&
				entry.TimeSlot.StartTime == otherEntry.TimeSlot.StartTime {
				violations++
			}
		}
	}
	return violations / 2
}

// countSoftViolations counts the number of soft constraint violations in the current
// assignments
func (s *Scheduler) countSoftViolations() int {
	violations := 0
	for _, entry := range s.Assignments {
		if len(entry.Class.PreferredDays) > 0 {
			found := false
			for _, day := range entry.Class.PreferredDays {
				if day == entry.TimeSlot.Day {
					found = true
					break
				}
			}
			if !found {
				violations++
			}
		}
	}
	return violations
}
