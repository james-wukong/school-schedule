package solver

import (
	"math/rand/v2"
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

func fixedRNG() *rand.Rand {
	seed1 := uint64(10)
	seed2 := uint64(20)

	pcg := rand.NewPCG(seed1, seed2)
	return rand.New(pcg)
}

func randomRNG() *rand.Rand {
	num := rand.IntN(10)
	switch num {
	case 0:
		seed1, seed2 := uint64(10), uint64(19)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 1:
		seed1, seed2 := uint64(19), uint64(29)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 2:
		seed1, seed2 := uint64(29), uint64(39)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 3:
		seed1, seed2 := uint64(39), uint64(49)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 4:
		seed1, seed2 := uint64(49), uint64(59)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 5:
		seed1, seed2 := uint64(59), uint64(69)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 6:
		seed1, seed2 := uint64(69), uint64(79)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 7:
		seed1, seed2 := uint64(79), uint64(299)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 8:
		seed1, seed2 := uint64(199), uint64(89)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	case 9:
		seed1, seed2 := uint64(89), uint64(99)
		pcg := rand.NewPCG(seed1, seed2)
		return rand.New(pcg)
	default:
		return fixedRNG()
	}
}

// countPlaced returns how many assignments were placed for a specific requirement ID
func countPlaced(assignments []*model.Assignment, reqID model.RequirementID) int {
	count := 0
	for _, a := range assignments {
		if a.Requirement.ID == reqID {
			count++
		}
	}
	return count
}

// hasTeacherConflict returns true if any teacher appears twice in the same slot
func hasTeacherConflict(assignments []*model.Assignment) bool {
	seen := make(map[model.ConflictMapKey[model.TeacherID]]bool)
	for _, a := range assignments {
		key := model.ConflictMapKey[model.TeacherID]{
			ID:   a.Requirement.Teacher.ID,
			Day:  a.Slot.Day,
			Slot: a.Slot.StartTime,
		}
		if _, ok := seen[key]; ok {
			return true
		}
		seen[key] = true
	}
	return false
}

// hasRoomConflict returns true if any room is used twice in the same slot
func hasRoomConflict(assignments []*model.Assignment) bool {
	seen := make(map[model.ConflictMapKey[model.RoomID]]bool)
	for _, a := range assignments {
		key := model.ConflictMapKey[model.RoomID]{
			ID:   a.Room.ID,
			Day:  a.Slot.Day,
			Slot: a.Slot.StartTime,
		}
		if _, ok := seen[key]; ok {
			return true
		}
		seen[key] = true
	}
	return false
}

// hasClassConflict returns true if any class has two lessons at the same slot
func hasClassConflict(assignments []*model.Assignment) bool {
	seen := make(map[model.ConflictMapKey[model.ClassID]]bool)
	for _, a := range assignments {
		key := model.ConflictMapKey[model.ClassID]{
			ID:   a.Requirement.SchoolClass.ID,
			Day:  a.Slot.Day,
			Slot: a.Slot.StartTime,
		}
		if _, ok := seen[key]; ok {
			return true
		}
		seen[key] = true
	}
	return false
}

// usedRoomTypes returns the set of room types used by a requirement's assignments
func usedRoomTypes(assignments []*model.Assignment,
	reqID model.RequirementID,
) map[model.RoomType]int {
	types := make(map[model.RoomType]int)
	for _, a := range assignments {
		if a.Requirement.ID == reqID {
			types[a.Room.Type] += 1
		}
	}
	return types
}

// slotsUsedByTeacher returns all slots assigned to a teacher
func slotsUsedByTeacher(
	assignments []*model.Assignment,
	teacherID model.TeacherID,
) []model.TimeSlot {
	var slots []model.TimeSlot
	for _, a := range assignments {
		if a.Requirement.Teacher.ID == teacherID {
			slots = append(slots, a.Slot)
		}
	}
	return slots
}

// ─────────────────────────────────────────────
//  CASE 1 — BASIC PLACEMENT
//  Single requirement, 1 session/week
//  Expected: exactly 1 assignment placed
// ─────────────────────────────────────────────

func TestGreedy_SingleSession(t *testing.T) {
	rooms := make([]*model.Room, 0, 1)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()
	rooms = append(rooms, room101)
	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 1,
		}, // Grade 1 | Class 01 | Math | Alice | 1x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)
	t.Logf("len(Assignment) =%d", len(result))
	if len(result) > 0 {
		for _, r := range result {
			t.Logf("Assignment =%+v", r)
		}
	}

	if len(result) != 1 {
		t.Errorf("expected 1 assignment, got %d", len(result))
	}
}

// ─────────────────────────────────────────────
//  CASE 2 — SESSIONS PER WEEK ARE FULLY EXPANDED
//  One requirement needs 4 sessions/week
//  Expected: exactly 4 assignments placed
// ─────────────────────────────────────────────

func TestGreedy_MultipleSessionsExpanded(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	rooms = append(rooms, room101, room102)

	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | Math | Alice | 4x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	placed := countPlaced(result, 1000)
	if placed != 4 {
		t.Errorf("expected 4 sessions placed for REQ1, got %d", placed)
	}
}

// ─────────────────────────────────────────────
//  CASE 3 — NO TEACHER DOUBLE-BOOKING
//  Two requirements sharing the same teacher
//  Expected: no teacher conflict in result
// ─────────────────────────────────────────────

func TestGreedy_NoTeacherConflict(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	rooms = append(rooms, room101, room102)
	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Math | Alice | 3x
		{
			ID:              1001,
			SchoolClass:     classB,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 02 | Math | Alice | 3x <- same teacher
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if hasTeacherConflict(result) {
		t.Error("teacher conflict detected: Alice placed in two classes at the same slot")
	}
}

// ─────────────────────────────────────────────
//  CASE 4 — NO ROOM DOUBLE-BOOKING
//  Multiple requirements competing for limited rooms
//  Expected: no room conflict in result
// ─────────────────────────────────────────────

func TestGreedy_NoRoomConflict(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	alice, bob := makeTeacherFixtures()
	math, eng, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	rooms = append(rooms, room101, room102)

	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | Math | Alice | 4x
		{
			ID:              1001,
			SchoolClass:     classB,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 02 | English | Bob | 4x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if hasRoomConflict(result) {
		t.Error("room conflict detected: same room used by two classes at the same slot")
	}
}

// ─────────────────────────────────────────────
//  CASE 5 — NO CLASS DOUBLE-BOOKING
//  One class has multiple subjects
//  Expected: class never has two lessons at once
// ─────────────────────────────────────────────

func TestGreedy_NoClassConflict(t *testing.T) {
	rooms := make([]*model.Room, 0, 3)
	alice, bob := makeTeacherFixtures()
	carol := &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	math, eng, sci := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, lab101 := makeRoomFixtures()
	rooms = append(rooms, room101, room102, lab101)

	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | Math | Alice | 4x
		{
			ID:              1001,
			SchoolClass:     classA,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | English | Bob | 4x <- same class
		{
			ID:              1002,
			SchoolClass:     classA,
			Subject:         sci,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Science | Carol | 4x <- same class
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if hasClassConflict(result) {
		t.Error("class conflict detected: Year 10A has two lessons at the same slot")
	}
}

// ─────────────────────────────────────────────
//  CASE 6 — TEACHER AVAILABILITY IS RESPECTED
//  Alice is only available on Monday
//  Expected: all her assignments land on Monday
// ─────────────────────────────────────────────

func TestGreedy_TeacherAvailabilityRespected(t *testing.T) {
	rooms := make([]*model.Room, 0, 1)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()
	rooms = append(rooms, room101)

	// Block all days except Monday
	for _, day := range []model.DayOfWeek{
		model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
	} {
		delete(alice.AvailableTimes, day)
	}
	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Math | Alice | 3x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	for _, a := range result {
		if a.Requirement.Teacher.ID == 1000 && a.Slot.Day != model.Monday {
			t.Errorf("Alice was placed on %d but is only available Monday", a.Slot.Day)
		}
	}
}

// ─────────────────────────────────────────────
//  CASE 7 — ROOM CAPACITY IS RESPECTED
//  Class has 30 students; only Room-102 (cap 35) fits
//  Room-101 has capacity 20 and should never be used
//  Expected: all assignments use Room-102
// ─────────────────────────────────────────────

func TestGreedy_RoomCapacityRespected(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	room101.Capacity = 20
	room102.Capacity = 30
	rooms = append(rooms, room101, room102)

	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 01 | Math | Alice | 2x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	for _, a := range result {
		if a.Room.ID == 101 {
			t.Errorf("Room-101 (capacity 20) was used for a class of 30 students")
		}
	}
	if len(result) != 2 {
		t.Errorf("expected 2 sessions placed, got %d", len(result))
	}
}

// ─────────────────────────────────────────────
//  CASE 8 — LAB REQUIREMENT IS RESPECTED
//  Science requires a lab; regular rooms must be skipped
//  Expected: all science assignments use lab rooms
// ─────────────────────────────────────────────

func TestGreedy_LabRequirementRespected(t *testing.T) {
	rooms := make([]*model.Room, 0, 3)
	_, _, science := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, labA := makeRoomFixtures()
	rooms = append(rooms, room101, room102, labA)
	carol := &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	reqs := []*model.Requirement{
		{
			ID:              1000,
			SchoolClass:     classA,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Science | Carol | 3x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if len(result) != 3 {
		t.Errorf("expected 3 science sessions placed, got %d", len(result))
	}
	for _, a := range result {
		if a.Room.Type != model.Lab {
			t.Errorf("science session placed in non-lab room: %s (%s)",
				a.Room.Name, a.Room.Type)
		}
	}
}

// ─────────────────────────────────────────────
//  CASE 9 — FULL WEEK MULTI-CLASS SCHEDULE
//  4 classes, 4 teachers, 5 subjects
//  Expected: all sessions placed, zero hard violations
// ─────────────────────────────────────────────

func TestGreedy_FullWeekNoConflicts(t *testing.T) {
	rooms := make([]*model.Room, 0, 4)
	alice, bob := makeTeacherFixtures()
	carol := &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	dave := &model.Teacher{
		ID: 1011, TenantID: 1, Name: "Dave",
		Subjects: []model.SubjectID{104},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	math, eng, science := makeSubjectFixtures()
	history := &model.Subject{
		ID:          104,
		Name:        "History",
		RequiresLab: false,
		IsHeavy:     false,
	}
	cls10A, cls10B, _ := makeClassFixtures()
	room101, room102, labA := makeRoomFixtures()
	room103 := &model.Room{ID: 103, TenantID: 1, Name: "Room-103", Capacity: 35, Type: model.Regular,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}
	rooms = append(rooms, room101, room102, room103, labA)

	reqs := []*model.Requirement{
		{
			ID:              1001,
			SchoolClass:     cls10A,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | Math | Alice | 4x
		{
			ID:              1002,
			SchoolClass:     cls10A,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | English | Bob | 3x
		{
			ID:              1003,
			SchoolClass:     cls10A,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Science | Carol | 3x
		{
			ID:              1004,
			SchoolClass:     cls10A,
			Subject:         history,
			Teacher:         dave,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 01 | History | Dave | 2x
		{
			ID:              1005,
			SchoolClass:     cls10B,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 02 | Math | Alice | 4x
		{
			ID:              1006,
			SchoolClass:     cls10B,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 02 | English | Bob | 4x
		{
			ID:              1007,
			SchoolClass:     cls10B,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 02 | Science | Carol | 3x
		{
			ID:              1008,
			SchoolClass:     cls10B,
			Subject:         history,
			Teacher:         dave,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 02 | History | Dave | 2x
	}

	totalSessions := 0
	for _, r := range reqs {
		totalSessions += r.SessionsPerWeek
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	// All sessions placed
	if len(result) != totalSessions {
		t.Errorf("expected %d sessions placed, got %d", totalSessions, len(result))
	}

	// Zero hard violations
	hard := HardViolations(result)
	if hard != 0 {
		t.Errorf("expected 0 hard violations after greedy construct, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 10 — BLOCKED TEACHER PRODUCES FEWER SESSIONS
//  Alice is fully blocked (no availability at all)
//  Expected: 0 sessions placed for her requirement
// ─────────────────────────────────────────────

func TestGreedy_FullyBlockedTeacherPlacesNothing(t *testing.T) {
	rooms := make([]*model.Room, 0, 1)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()
	rooms = append(rooms, room101)

	// Block every single slot
	alice.AvailableTimes = nil
	reqs := []*model.Requirement{
		{
			ID:              1001,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 01 | Math | Alice | 4x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	placed := countPlaced(result, 1001)
	if placed != 0 {
		t.Errorf("expected 0 sessions for fully blocked teacher, got %d", placed)
	}
}

// ─────────────────────────────────────────────
//  CASE 11 — ALL ROOMS TOO SMALL
//  Class has 50 students but all rooms hold 20
//  Expected: 0 sessions placed
// ─────────────────────────────────────────────

func TestGreedy_AllRoomsTooSmall(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	rooms = append(rooms, room101, room102)
	classA.StudentCount = 50

	reqs := []*model.Requirement{
		{
			ID:              1001,
			SchoolClass:     classA,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Math | Alice | 3x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if len(result) != 0 {
		t.Errorf("expected 0 sessions when all rooms are too small, got %d", len(result))
	}
}

// ─────────────────────────────────────────────
//  CASE 12 — NO LAB AVAILABLE FOR SCIENCE
//  Science requires a lab but none exist
//  Expected: 0 sessions placed
// ─────────────────────────────────────────────

func TestGreedy_NoLabForSciencePlacesNothing(t *testing.T) {
	rooms := make([]*model.Room, 0, 2)
	_, _, science := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()
	rooms = append(rooms, room101, room102)

	carol := &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}

	reqs := []*model.Requirement{
		{
			ID:              1001,
			SchoolClass:     classA,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 01 | Math | Alice | 3x
	}

	result := GreedyConstruct(fixedRNG(), reqs, rooms)

	if len(result) != 0 {
		t.Errorf("expected 0 sessions when no lab exists, got %d", len(result))
	}
}
