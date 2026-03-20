package solver

// Hard violations count teacher/class/room double-bookings.
// Soft penalties cover teacher gap windows (+2), same subject twice in one day (+3),
// and heavy subjects placed in period 6+ (+1).

import (
	"fmt"
	"slices"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

// BuildConflictIndex creates conflict index
func BuildConflictIndex(assignments []*model.Assignment) (*model.ConflictIndex, int) {
	idx := model.NewConflictIndex()
	count := 0
	for _, a := range assignments {
		tk := model.ConflictKey(a.Requirement.Teacher.ID, a.Slot)
		ck := model.ConflictKey(a.Requirement.SchoolClass.ID, a.Slot)
		rk := model.ConflictKey(a.Room.ID, a.Slot)

		idx.TeacherSlot[tk] = append(idx.TeacherSlot[tk], a)
		idx.ClassSlot[ck] = append(idx.ClassSlot[ck], a)
		idx.RoomSlot[rk] = append(idx.RoomSlot[rk], a)
		// A teacher can only teach subjects they're qualified for
		if !slices.Contains(a.Requirement.Teacher.Subjects, a.Requirement.Subject.ID) {
			count++
		}
	}

	return idx, count
}

// HardViolations returns total hard violations
// A teacher can't be in two rooms at the same time
// A class can't have two lessons simultaneously
// A room can't host two classes at the same time
// A teacher can only teach subjects they're qualified for
func HardViolations(assignments []*model.Assignment) int {
	idx, unqualified := BuildConflictIndex(assignments)
	count := 0
	for _, lst := range idx.TeacherSlot {
		if len(lst) > 1 {
			count += len(lst) - 1
		}
	}
	for _, lst := range idx.ClassSlot {
		if len(lst) > 1 {
			count += len(lst) - 1
		}
	}
	for _, lst := range idx.RoomSlot {
		if len(lst) > 1 {
			count += len(lst) - 1
		}
	}
	return count + unqualified
}

// SoftViolations returns total soft violations
// Teachers prefer not to have isolated single periods ("windows")
// Classes shouldn't have the same subject twice in one day
// Heavy subjects (Math, Science) are better placed in the morning
// Teachers shouldn't have more than N consecutive periods
func SoftViolations(assignments []*model.Assignment) float64 {
	penalty := 0.0

	teacherAssignedSlots := make(map[model.TeacherDaySlotsKey][]string, 0)
	nonheavyDaySubjects := make(map[model.ClassDaySubjectKey][]model.SubjectID, 0)

	for _, a := range assignments {
		tk := model.TeacherDaySlotsKey{
			TeacherID: a.Requirement.Teacher.ID,
			Day:       a.Slot.Day,
		}
		teacherAssignedSlots[tk] = append(teacherAssignedSlots[tk], a.Slot.StartTime)

		ck := model.ClassDaySubjectKey{
			ClassID: a.Requirement.SchoolClass.ID,
			Day:     a.Slot.Day,
		}
		if !a.Requirement.Subject.IsHeavy {
			nonheavyDaySubjects[ck] = append(nonheavyDaySubjects[ck], a.Requirement.Subject.ID)
		}
		// 1. Heavy subjects placed in late periods (P6+)
		isAfter, err := isTimeAfter(a.Slot.StartTime, "14:00")
		if err != nil {
			fmt.Printf("error caught transiting string to time: %v\n", err)
		}
		if a.Requirement.Subject.IsHeavy && isAfter {
			penalty += 1.0
		}

		// 2. Preferred Day in requirements is not met
		if slices.Contains(a.Requirement.PreferredDays, a.Slot.Day) {
			penalty += 2.0
		}
	}

	// 3. Teacher gap (window) detection per day
	for _, slots := range teacherAssignedSlots {
		slices.Sort(slots)
		for i := 1; i < len(slots); i++ {
			if minuteDifference(slots[i], slots[i-1]) > 10 {
				penalty += 3.0
			}
		}
	}

	// 3. Same subject twice in one day for a non-heavy subject
	for _, subjects := range nonheavyDaySubjects {
		slices.Sort(subjects)
		for k, v := range subjects {
			if k >= 1 && subjects[k-1] == v {
				penalty += 5.0
			}
		}
	}

	return penalty
}

const hardWeight = 1000.0

func TotalCost(assignments []*model.Assignment) float64 {
	return hardWeight*float64(HardViolations(assignments)) + SoftViolations(assignments)
}
