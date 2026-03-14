package solver

import (
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

const (
	TimeLayout        = "2006-01-02 15:04:05"
	TimeMinutesLayout = "15:04"
)

// isTeacherAvailable checks if the teacher is available at the given time slot
func isTeacherAvailable(teacher *model.Teacher, slot model.TimeSlot) bool {
	if times, ok := teacher.AvailableHours[slot.Day]; ok {
		for _, t := range times {
			if t == slot.StartTime {
				return true
			}
		}
	}
	return false
}

// isRoomAvailableInDomain checks if the room is available at the given time slot
// based on the current domain (possible assignments)
func isRoomAvailableInDomain(classroom *model.Classroom, slot model.TimeSlot) bool {
	if times, ok := classroom.AvailableTimes[slot.Day]; ok {
		for _, t := range times {
			if t == slot.StartTime {
				return true
			}
		}
	}
	return false
}

// dayDifference returns the number of days between two days of the week
func dayDifference(startDay, endDay model.DayOfWeek) int {
	if endDay < startDay {
		return 10
	}
	return int(endDay) - int(startDay)
}

// countClassesPerDay counts how many classes a teacher has on a given day
func countClassesPerDay(assignments map[int]*ScheduleEntry, teacherID int, day model.DayOfWeek) int {
	count := 0
	for _, entry := range assignments {
		if entry.Teacher.ID == teacherID && entry.TimeSlot.Day == day {
			count++
		}
	}
	return count
}

// randomInt returns a random integer from 0 to max-1
func randomInt(maxInt int) int {
	if maxInt <= 0 {
		return 0
	}
	return int(time.Now().UnixNano()) % maxInt
}

// acceptWithProbability decides whether to accept a worse solution
// based on the change in violations and current temperature
func acceptWithProbability(delta int, temperature float64) bool {
	if temperature == 0 {
		return false
	}
	probability := 1.0 / (1.0 + float64(delta)/temperature)
	return probability > 0.5
}

// costFunction calculates total violations (lower is better)
func costFunction(assignments map[int]*ScheduleEntry) int {
	cost := 0

	// Hard constraints (weight: 1000)
	for _, entry := range assignments {
		class := entry.Class
		teacher := entry.Teacher

		// Check conflicts with other assignments
		for otherID, otherEntry := range assignments {
			if otherID == entry.Class.ID {
				continue
			}

			// Teacher conflict
			if teacher.ID == otherEntry.Teacher.ID &&
				entry.TimeSlot.Day == otherEntry.TimeSlot.Day &&
				entry.TimeSlot.StartTime == otherEntry.TimeSlot.StartTime {
				cost += 1000
			}

			// Room conflict
			if entry.Classroom.ID == otherEntry.Classroom.ID &&
				entry.TimeSlot.Day == otherEntry.TimeSlot.Day &&
				entry.TimeSlot.StartTime == otherEntry.TimeSlot.StartTime {
				cost += 1000
			}

			// Day gap violation
			if teacher.ID == otherEntry.Teacher.ID {
				dayDiff := dayDifference(entry.TimeSlot.Day, otherEntry.TimeSlot.Day)
				if dayDiff < class.MinDayGap && dayDiff > 0 {
					cost += 1000
				}
			}
		}

		// Teacher workload balance (soft constraint, weight: 10)
		classesPerDay := countClassesPerDay(assignments, teacher.ID, entry.TimeSlot.Day)
		if classesPerDay > teacher.MaxClassesPerDay {
			cost += 10 * (classesPerDay - teacher.MaxClassesPerDay)
		}
	}

	// Soft constraints - preferred days (weight: 5)
	for _, entry := range assignments {
		if len(entry.Class.PreferredDays) > 0 {
			found := false
			for _, day := range entry.Class.PreferredDays {
				if day == entry.TimeSlot.Day {
					found = true
					break
				}
			}
			if !found {
				cost += 5
			}
		}
	}

	// Soft constraints - teacher preferences (weight: 2)
	for _, entry := range assignments {
		if preference, ok := entry.Teacher.Preferences[entry.TimeSlot.Day]; ok {
			cost += (10 - preference) * 2 // Prefer high preference values
		}
	}

	return cost
}
