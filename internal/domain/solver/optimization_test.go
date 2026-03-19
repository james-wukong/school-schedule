package solver

import (
	"slices"
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

// ─────────────────────────────────────────────
//  HELPERS
// ─────────────────────────────────────────────

func makeMoreTeacherFixtures() (
	carol, dave, eve *model.Teacher,
) {
	carol = &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	dave = &model.Teacher{
		ID: 1011, TenantID: 1, Name: "Dave",
		Subjects: []model.SubjectID{104},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	eve = &model.Teacher{
		ID: 1012, TenantID: 1, Name: "Eve",
		Subjects: []model.SubjectID{102, 105},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}

	return
}

func makeMoreSubjectFixtures() (
	history, pe *model.Subject,
) {
	history = &model.Subject{ID: 104, Name: "History", RequiresLab: false, IsHeavy: false}
	pe = &model.Subject{ID: 105, Name: "PE", RequiresLab: false, IsHeavy: false}

	return
}

func makeMoreRoomFixtures() (
	room103, labB, gimA *model.Room,
) {
	room103 = &model.Room{ID: 103, TenantID: 1, Name: "Room-103", Capacity: 35, Type: model.Regular,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}
	labB = &model.Room{ID: 1002, TenantID: 1, Name: "Lab-102", Capacity: 32, Type: model.Lab,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}
	gimA = &model.Room{ID: 1010, TenantID: 1, Name: "Gym-100", Capacity: 60, Type: model.GYM,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}

	return
}

// buildCleanSchedule returns a conflict-free set of assignments
// using distinct slots for each session — a known-good starting point.
func buildCleanSchedule() ([]*model.Assignment, []*model.Room) {
	rooms := make([]*model.Room, 0, 3)
	alice, bob := makeTeacherFixtures()
	carol, _, _ := makeMoreTeacherFixtures()
	math, eng, science := makeSubjectFixtures()
	room101, room102, labA := makeRoomFixtures()
	cls10A, cls10B, _ := makeClassFixtures()
	rooms = append(rooms, room101, room102, labA)

	reqA := &model.Requirement{
		ID:              1001,
		SchoolClass:     cls10A,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 01 | Math | Alice | 1x
	reqB := &model.Requirement{
		ID:              1002,
		SchoolClass:     cls10B,
		Subject:         eng,
		Teacher:         bob,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 01 | English | Bob | 1x
	reqC := &model.Requirement{
		ID:              1003,
		SchoolClass:     cls10A,
		Subject:         science,
		Teacher:         carol,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 01 | Science | Carol | 1x

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Grade 1 | Class 01 | Math | Alice | Mon-09:00 | Room-101
		{
			Requirement: reqB, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "10:00"},
		}, // Grade 1 | Class 02 | English | Bob | Mon-10:00 | Room-102
		{
			Requirement: reqC, Room: labA,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "11:00"},
		}, // Grade 1 | Class 03 | Science | Carol | Mon-11:00 | Lab-1001
	}

	return assignments, rooms
}

// buildConflictedSchedule returns assignments that have known hard violations:
// Alice is double-booked and Room-101 is double-booked at Mon-09:00.
func buildConflictedSchedule() ([]*model.Assignment, []*model.Room) {
	rooms := make([]*model.Room, 0, 3)
	alice, bob := makeTeacherFixtures()
	math, eng, _ := makeSubjectFixtures()
	room101, room102, _ := makeRoomFixtures()
	room103, _, _ := makeMoreRoomFixtures()
	cls10A, cls10B, cls11A := makeClassFixtures()
	rooms = append(rooms, room101, room102, room103)

	reqA := &model.Requirement{
		ID:              1001,
		SchoolClass:     cls10A,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 01 | Math | Alice | 1x
	reqB := &model.Requirement{
		ID:              1002,
		SchoolClass:     cls10B,
		Subject:         eng,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 02 | English | Alice | 1x <- Alice conflict
	reqC := &model.Requirement{
		ID:              1003,
		SchoolClass:     cls11A,
		Subject:         math,
		Teacher:         bob,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 03 | Math | Bob | 1x

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Grade 1 | Class 01 | Math | Alice | Mon-09:00 | Room-101
		{
			Requirement: reqB, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Grade 1 | Class 02 | English | Alice | Mon-09:00 | Room-101 ← teacher + room clash
		{
			Requirement: reqC, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "10:00"},
		}, // Grade 1 | Class 03 | Math | Bob | Mon-10:00 | Room-102

	}

	return assignments, rooms
}

// buildLargeSchedule builds a realistic multi-class schedule via GreedyConstruct
// giving SA something meaningful to optimise.
func buildLargeSchedule() ([]*model.Assignment, []*model.Room) {
	rooms := make([]*model.Room, 0, 6)
	alice, bob := makeTeacherFixtures()
	carol, dave, eve := makeMoreTeacherFixtures()
	math, eng, science := makeSubjectFixtures()
	history, pe := makeMoreSubjectFixtures()
	cls10A, cls10B, cls11A := makeClassFixtures()
	room101, room102, labA := makeRoomFixtures()
	room103, labB, gymA := makeMoreRoomFixtures()
	rooms = append(rooms, room101, room102, room103, labA, labB, gymA)

	reqs := []*model.Requirement{
		{
			ID:              1001,
			SchoolClass:     cls10A,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 1,
		}, // Grade 1 | Class 01 | Math | Alice | 1x
		{
			ID:              1002,
			SchoolClass:     cls10A,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 1,
		}, // Grade 1 | Class 01 | English | Bob | 1x
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
			SchoolClass:     cls10A,
			Subject:         pe,
			Teacher:         eve,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 01 | PE | Eve | 2x
		{
			ID:              1006,
			SchoolClass:     cls10B,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 02 | Math | Alice | 4x
		{
			ID:              1007,
			SchoolClass:     cls10B,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 02 | English | Bob | 4x
		{
			ID:              1008,
			SchoolClass:     cls10B,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 02 | Science | Carol | 3x
		{
			ID:              1009,
			SchoolClass:     cls10B,
			Subject:         history,
			Teacher:         dave,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 02 | History | Dave | 2x
		{
			ID:              1010,
			SchoolClass:     cls10B,
			Subject:         pe,
			Teacher:         eve,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 02 | PE | eve | 2x
		{
			ID:              1011,
			SchoolClass:     cls11A,
			Subject:         math,
			Teacher:         alice,
			SessionsPerWeek: 5,
		}, // Grade 1 | Class 03 | Math | Alice | 5x
		{
			ID:              1012,
			SchoolClass:     cls11A,
			Subject:         eng,
			Teacher:         bob,
			SessionsPerWeek: 4,
		}, // Grade 1 | Class 03 | English | Bob | 4x
		{
			ID:              1013,
			SchoolClass:     cls11A,
			Subject:         science,
			Teacher:         carol,
			SessionsPerWeek: 3,
		}, // Grade 1 | Class 03 | Science | Carol | 3x
		{
			ID:              1014,
			SchoolClass:     cls11A,
			Subject:         history,
			Teacher:         dave,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 03 | History | Dave | 2x
		{
			ID:              1015,
			SchoolClass:     cls11A,
			Subject:         pe,
			Teacher:         eve,
			SessionsPerWeek: 2,
		}, // Grade 1 | Class 03 | PE | Eve | 2x
	}

	rng := fixedRNG()
	assignments := GreedyConstruct(rng, reqs, rooms)
	return assignments, rooms
}

// totalSessions counts the number of assignments for a given requirement ID
func totalSessions(assignments []*model.Assignment, reqID model.RequirementID) int {
	count := 0
	for _, a := range assignments {
		if a.Requirement.ID == reqID {
			count++
		}
	}
	return count
}

// requirementIDs extracts the set of unique requirement IDs from assignments
func requirementIDs(assignments []*model.Assignment) []model.RequirementID {
	ids := make([]model.RequirementID, 0)
	for _, a := range assignments {
		if !slices.Contains(ids, a.Requirement.ID) {
			ids = append(ids, a.Requirement.ID)
		}
	}
	return ids
}

// ─────────────────────────────────────────────
//  CASE 1 — SESSION COUNT IS PRESERVED
//  SA must never add or drop sessions — only move them.
//  Expected: len(result) == len(initial)
// ─────────────────────────────────────────────

func TestSA_PreservesSessionCount(t *testing.T) {
	initial, rooms := buildCleanSchedule()
	rng := fixedRNG()

	result := SimulatedAnnealing(rng, initial, rooms, 100.0, 0.99, 500)

	if len(result) != len(initial) {
		t.Errorf("SA changed session count: started with %d, ended with %d",
			len(initial), len(result))
	}
}

// ─────────────────────────────────────────────
//  CASE 2 — REQUIREMENT IDENTITY IS PRESERVED
//  SA only moves assignments to new slots/rooms;
//  it must not swap which requirement an assignment belongs to.
//  Expected: same set of requirement IDs before and after
// ─────────────────────────────────────────────

func TestSA_PreservesRequirementIdentity(t *testing.T) {
	initial, rooms := buildLargeSchedule()
	rng := fixedRNG()

	before := requirementIDs(initial)
	result := SimulatedAnnealing(rng, initial, rooms, 500.0, 0.997, 2000)
	after := requirementIDs(result)

	for _, id := range before {
		if !slices.Contains(after, id) {
			t.Errorf("requirement %d was present before SA but missing after", id)
		}
	}
	for _, id := range after {
		if !slices.Contains(before, id) {
			t.Errorf("requirement %d appeared after SA but was not in initial", id)
		}
	}
}

// ─────────────────────────────────────────────
//  CASE 3 — PERFECT SCHEDULE STAYS PERFECT
//  A conflict-free, penalty-free initial schedule
//  should not be made worse by SA.
//  Expected: hard violations remain 0
// ─────────────────────────────────────────────

func TestSA_PerfectScheduleRemainsValid(t *testing.T) {
	initial, rooms := buildCleanSchedule()
	rng := fixedRNG()

	// Confirm starting point is clean
	if HardViolations(initial) != 0 {
		t.Fatal("test setup error: initial schedule already has hard violations")
	}

	result := SimulatedAnnealing(rng, initial, rooms, 10.0, 0.99, 500)

	// SA returns the BEST seen, so hard violations must not increase
	if HardViolations(result) > 0 {
		t.Errorf("SA made a perfect schedule worse: %d hard violations after SA",
			HardViolations(result))
	}
}

// ─────────────────────────────────────────────
//  CASE 4 — SA IMPROVES A CONFLICTED SCHEDULE
//  Start with known hard violations;
//  SA should reduce total cost.
//  Expected: TotalCost(result) < TotalCost(initial)
// ─────────────────────────────────────────────

func TestSA_ImprovesConflictedSchedule(t *testing.T) {
	initial, rooms := buildConflictedSchedule()
	rng := fixedRNG()

	initialCost := TotalCost(initial)
	if initialCost == 0 {
		t.Fatal("test setup error: initial schedule has no violations to improve")
	}

	result := SimulatedAnnealing(rng, initial, rooms, 500.0, 0.99, 3000)
	finalCost := TotalCost(result)

	if finalCost >= initialCost {
		t.Errorf("SA did not improve the schedule: initial cost=%.1f, final cost=%.1f",
			initialCost, finalCost)
	}
}

// ─────────────────────────────────────────────
//  CASE 5 — ZERO ITERATIONS RETURNS INITIAL
//  With 0 iterations SA has no chance to change anything.
//  Expected: result is identical to initial (same slots & rooms)
// ─────────────────────────────────────────────

func TestSA_ZeroIterationsReturnsInitial(t *testing.T) {
	initial, rooms := buildCleanSchedule()
	rng := fixedRNG()

	result := SimulatedAnnealing(rng, initial, rooms, 100.0, 0.99, 0)

	if len(result) != len(initial) {
		t.Fatalf("expected %d assignments, got %d", len(initial), len(result))
	}
	for i, a := range result {
		b := initial[i]
		if a.Slot != b.Slot || a.Room.ID != b.Room.ID {
			t.Errorf("assignment %d changed despite 0 iterations", i)
		}
	}
}

// ─────────────────────────────────────────────
//  CASE 6 — MORE ITERATIONS YIELDS LOWER OR EQUAL COST
//  Running SA longer should not produce a worse result.
//  Expected: cost(2000 iters) <= cost(200 iters)
// ─────────────────────────────────────────────

func TestSA_MoreIterationsDoesNotWorsen(t *testing.T) {
	initial, rooms := buildConflictedSchedule()

	rng1 := fixedRNG()
	rng2 := fixedRNG()

	short := SimulatedAnnealing(rng1, initial, rooms, 500.0, 0.99, 200)
	long := SimulatedAnnealing(rng2, initial, rooms, 500.0, 0.99, 2000)

	costShort := TotalCost(short)
	costLong := TotalCost(long)

	if costLong > costShort {
		t.Errorf(
			"longer SA produced worse result: 200 iters cost=%.1f, 2000 iters cost=%.1f",
			costShort, costLong,
		)
	}
}

// ─────────────────────────────────────────────
//  CASE 7 — RESULT NEVER WORSE THAN INITIAL
//  SA tracks the best-seen solution internally.
//  Even if the random walk goes uphill, the returned
//  result must be <= the starting cost.
//  Expected: TotalCost(result) <= TotalCost(initial)
// ─────────────────────────────────────────────

func TestSA_ResultNeverWorseThanInitial(t *testing.T) {
	initial, rooms := buildLargeSchedule()
	rng := fixedRNG()

	initialCost := TotalCost(initial)
	result := SimulatedAnnealing(rng, initial, rooms, 800.0, 0.997, 3000)
	finalCost := TotalCost(result)

	if finalCost > initialCost {
		t.Errorf(
			"SA returned a solution worse than the initial: initial=%.1f, final=%.1f",
			initialCost, finalCost,
		)
	}
}

// ─────────────────────────────────────────────
//  CASE 8 — SA ON LARGE REALISTIC SCHEDULE
//  A greedy-constructed schedule of 46 sessions across 3 classes.
//  SA should eliminate all hard violations.
//  Expected: hard violations == 0 after SA
// ─────────────────────────────────────────────

func TestSA_EliminatesHardViolationsOnLargeSchedule(t *testing.T) {
	initial, rooms := buildLargeSchedule()
	rng := fixedRNG()

	// Greedy may leave hard violations; SA should clean them up
	result := SimulatedAnnealing(rng, initial, rooms, 800.0, 0.997, 20_000)

	hard := HardViolations(result)
	if hard != 0 {
		t.Errorf("expected 0 hard violations after SA on large schedule, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 9 — ROOM ASSIGNMENTS STAY VALID AFTER SA
//  SA's "change room" move must never assign a class
//  to a room that is too small or wrong type.
//  Expected: no assignment violates capacity or lab requirement
// ─────────────────────────────────────────────

func TestSA_RoomAssignmentsStayValid(t *testing.T) {
	initial, rooms := buildLargeSchedule()
	rng := fixedRNG()

	result := SimulatedAnnealing(rng, initial, rooms, 500.0, 0.997, 5000)

	for _, a := range result {
		sc := a.Requirement.SchoolClass
		subj := a.Requirement.Subject
		room := a.Room

		if room.Capacity < sc.StudentCount {
			t.Errorf(
				"room %s (cap %d) used for %s-%s (%d students)",
				room.Name, room.Capacity, sc.Grade, sc.Class, sc.StudentCount,
			)
		}
		if subj.RequiresLab && room.Type != model.Lab {
			t.Errorf(
				"subject %s requires lab but placed in %s (%s)",
				subj.Name, room.Name, room.Type,
			)
		}
	}
}

// ─────────────────────────────────────────────
//  CASE 10 — COOLING RATE OF 1.0 (NO COOLING)
//  A cooling rate of 1.0 means temperature never drops —
//  SA behaves like pure random walk. It must still return
//  a result no worse than initial (best-tracking guarantee).
//  Expected: TotalCost(result) <= TotalCost(initial)
// ─────────────────────────────────────────────

func TestSA_NoCoolingStillTrackseBest(t *testing.T) {
	initial, rooms := buildConflictedSchedule()
	rng := fixedRNG()

	initialCost := TotalCost(initial)
	result := SimulatedAnnealing(rng, initial, rooms, 100.0, 1.0, 1000)
	finalCost := TotalCost(result)

	if finalCost > initialCost {
		t.Errorf(
			"no-cooling SA returned worse result: initial=%.1f, final=%.1f",
			initialCost, finalCost,
		)
	}
}

// ─────────────────────────────────────────────
//  CASE 11 — SINGLE ASSIGNMENT SCHEDULE
//  Minimal edge case: only one assignment to work with.
//  SA has nothing to swap, so result must equal initial.
//  Expected: len(result)==1, same requirement ID
// ─────────────────────────────────────────────

func TestSA_SingleAssignment(t *testing.T) {
	rooms := make([]*model.Room, 0, 1)
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()
	rooms = append(rooms, room101)

	req := &model.Requirement{
		ID:              1001,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // Grade 1 | Class 01 | Math | Alice | 1x
	initial := []*model.Assignment{
		{
			Requirement: req, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Grade 1 | Class 01 | Math | Alice | Mon-09:00 | Room-101
	}

	rng := fixedRNG()
	result := SimulatedAnnealing(rng, initial, rooms, 100.0, 0.99, 500)

	if len(result) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(result))
	}
	if result[0].Requirement.ID != 1001 {
		t.Errorf("requirement identity changed after SA on single assignment")
	}
}

// ─────────────────────────────────────────────
//  CASE 12 — DIFFERENT SEEDS BOTH CONVERGE
//  Two runs with different seeds should both
//  reach hard=0 on a solvable schedule.
//  Expected: hard==0 for both runs
// ─────────────────────────────────────────────

func TestSA_DifferentSeedsBothConverge(t *testing.T) {
	initial, rooms := buildLargeSchedule()

	for range 3 {
		rng := randomRNG()
		result := SimulatedAnnealing(rng, initial, rooms, 800.0, 0.997, 20_000)
		hard := HardViolations(result)
		if hard != 0 {
			t.Errorf("expected 0 hard violations, got %d", hard)
		}
	}
}
