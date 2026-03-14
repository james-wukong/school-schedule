package solver

import (
	"fmt"
)

// initializeDomains creates possible assignments for each class
func (s *Scheduler) initializeDomains() error {
	for classID, class := range s.Classes {
		var possibleAssignments []*ScheduleEntry

		teacher, ok := s.Teachers[class.TeacherID]
		if !ok {
			return fmt.Errorf("teacher %s not found", class.TeacherID)
		}

		// Generate all possible assignments for this class
		for _, room := range s.Classrooms {
			if room.Capacity < class.StudentCount {
				continue // Room too small
			}

			for _, slot := range s.TimeSlots {
				// Check teacher availability
				if !isTeacherAvailable(teacher, slot) {
					continue
				}

				// Check room availability in domain
				if !isRoomAvailableInDomain(room, slot) {
					continue
				}

				// Check preferred days if specified
				if len(class.PreferredDays) > 0 {
					found := false
					for _, day := range class.PreferredDays {
						if day == slot.Day {
							found = true
							break
						}
					}
					if !found {
						continue
					}
				}

				entry := &ScheduleEntry{
					Class:     class,
					Teacher:   teacher,
					Classroom: room,
					TimeSlot:  slot,
				}
				possibleAssignments = append(possibleAssignments, entry)
			}
		}

		if len(possibleAssignments) == 0 {
			return fmt.Errorf("no valid assignments for class %s", classID)
		}

		s.Domains[classID] = possibleAssignments
	}
	return nil
}
