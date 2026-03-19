package solver

import (
	"fmt"
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

// ─────────────────────────────────────────────
//  SHARED FIXTURES
//  Real entities reused across all test cases
// ─────────────────────────────────────────────

func makeTeacherFixtures() (
	alice, bob *model.Teacher,
) {
	alice = &model.Teacher{
		ID: 1000, TenantID: 1, Name: "Alice",
		Subjects: []model.SubjectID{101},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}
	bob = &model.Teacher{
		ID: 1001, TenantID: 1, Name: "Bob",
		Subjects: []model.SubjectID{102},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}

	return
}

func makeSubjectFixtures() (
	math, english, science *model.Subject,
) {
	math = &model.Subject{ID: 101, Name: "Math", RequiresLab: false, IsHeavy: true}
	english = &model.Subject{ID: 102, Name: "English", RequiresLab: false, IsHeavy: true}
	science = &model.Subject{ID: 103, Name: "Science", RequiresLab: true, IsHeavy: true}

	return
}

func makeClassFixtures() (
	classA, classB, classC *model.SchoolClass,
) {
	classA = model.NewSchoolClass(model.SchoolClass{
		ID: 101, TenantID: 1, StudentCount: 28, Grade: "1", Class: "01",
	})

	classB = model.NewSchoolClass(model.SchoolClass{
		ID: 102, TenantID: 1, StudentCount: 30, Grade: "1", Class: "02",
	})
	classC = model.NewSchoolClass(model.SchoolClass{
		ID: 103, TenantID: 1, StudentCount: 25, Grade: "1", Class: "03",
	})

	return
}

func makeRoomFixtures() (
	room101, room102, labA *model.Room,
) {
	room101 = &model.Room{ID: 101, TenantID: 1, Name: "Room-101", Capacity: 35, Type: model.Regular,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}
	room102 = &model.Room{ID: 102, TenantID: 1, Name: "Room-102", Capacity: 35, Type: model.Regular,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}
	labA = &model.Room{ID: 1001, TenantID: 1, Name: "Lab-101", Capacity: 32, Type: model.Lab,
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		)}

	return
}

// ─────────────────────────────────────────────
//  CASE 1 — NO VIOLATIONS (clean schedule)
//  3 classes, 3 teachers, 3 different slots
//  Expected: hard=0, soft=0
// ─────────────────────────────────────────────

func TestNoViolations(t *testing.T) {
	alice, bob := makeTeacherFixtures()
	math, english, science := makeSubjectFixtures()
	classA, classB, classC := makeClassFixtures()
	room101, room102, labA := makeRoomFixtures()

	carol := &model.Teacher{
		ID: 1010, TenantID: 1, Name: "Carol",
		Subjects: []model.SubjectID{103},
		AvailableTimes: model.AvailableTimeSlots([]model.DayOfWeek{
			model.Monday, model.Tuesday, model.Wednesday, model.Thursday, model.Friday,
		},
			[]string{"09:00", "10:00", "11", "13:00", "14:00", "15:00"},
		),
	}

	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classB,
		Subject:         english,
		Teacher:         bob,
		SessionsPerWeek: 1,
	}
	reqC := &model.Requirement{
		ID:              1002,
		SchoolClass:     classC,
		Subject:         science,
		Teacher:         carol,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Alice  | Mon-P1 | Room-101 | 10A | Math
		{
			Requirement: reqB, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "10:00"},
		}, // Bob    | Mon-P2 | Room-102 | 10B | English
		{
			Requirement: reqC, Room: labA,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "11:00"},
		}, // Carol  | Mon-P3 | Lab-A    | 11A | Science
	}

	hard := HardViolations(assignments)
	soft := SoftViolations(assignments)

	fmt.Printf("\n[Case 1] Clean schedule\n")
	fmt.Printf("  Hard violations: %d (want 0)\n", hard)
	fmt.Printf("  Soft penalty   : %.1f (want 0.0)\n", soft)

	if hard != 0 {
		t.Errorf("expected 0 hard violations, got %d", hard)
	}
	if soft != 0 {
		t.Errorf("expected 0.0 soft penalty, got %.1f", soft)
	}
}

// ─────────────────────────────────────────────
//  CASE 2 — TEACHER DOUBLE-BOOKED
//  Alice is assigned to TWO classes at Mon-09:00
//  Expected: hard=1
// ─────────────────────────────────────────────

func TestTeacherDoubleBooked(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, english, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()

	// Both requirements use Alice
	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classB,
		Subject:         english,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // ← Alice again!

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Alice | Mon-P1 | Room-101 | 10A
		{
			Requirement: reqB, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Alice | Mon-P1 | Room-102 | 10B  ← CONFLICT
	}

	hard := HardViolations(assignments)

	fmt.Printf("\n[Case 2] Teacher double-booked\n")
	fmt.Printf("  Hard violations: %d (want 1)\n", hard)

	if hard != 1 {
		t.Errorf("expected 1 hard violation, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 3 — ROOM DOUBLE-BOOKED
//  Room-101 is used by two different classes at Mon-09:00
//  Expected: hard=1
// ─────────────────────────────────────────────

func TestRoomDoubleBooked(t *testing.T) {
	alice, bob := makeTeacherFixtures()
	math, english, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()

	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classB,
		Subject:         english,
		Teacher:         bob,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // Alice | Mon-P1 | Room-101 | 10A
		{
			Requirement: reqB, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // Bob   | Mon-P1 | Room-101 | 10B  ← CONFLICT
	}

	hard := HardViolations(assignments)

	fmt.Printf("\n[Case 3] Room double-booked\n")
	fmt.Printf("  Hard violations: %d (want 1)\n", hard)

	if hard != 1 {
		t.Errorf("expected 1 hard violation, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 4 — CLASS DOUBLE-BOOKED
//  Grade 1 Class 01 is scheduled for two subjects at Mon-09:00
//  Expected: hard=1
// ─────────────────────────────────────────────

func TestClassDoubleBooked(t *testing.T) {
	alice, bob := makeTeacherFixtures()
	math, english, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()

	// Both requirements belong to classA
	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classA,
		Subject:         english,
		Teacher:         bob,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // 10A | Mon-P1 | Math
		{
			Requirement: reqB, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // 10A | Mon-P1 | English  ← CONFLICT
	}

	hard := HardViolations(assignments)

	fmt.Printf("\n[Case 4] Class double-booked\n")
	fmt.Printf("  Hard violations: %d (want 1)\n", hard)

	if hard != 1 {
		t.Errorf("expected 1 hard violation, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 5 — MULTIPLE HARD VIOLATIONS AT ONCE
//  Alice double-booked AND Room-101 double-booked in same clash
//  Expected: hard=2
// ─────────────────────────────────────────────

func TestMultipleHardViolations(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, english, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()

	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classB,
		Subject:         english,
		Teacher:         alice,
		SessionsPerWeek: 1,
	} // <- Alice again

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		},
		{
			Requirement: reqB, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // same teacher + same room
	}

	hard := HardViolations(assignments)

	fmt.Printf("\n[Case 5] Teacher + Room double-booked simultaneously\n")
	fmt.Printf("  Hard violations: %d (want 2)\n", hard)

	if hard != 2 {
		t.Errorf("expected 2 hard violations, got %d", hard)
	}
}

// ─────────────────────────────────────────────
//  CASE 6 — SOFT: SAME SUBJECT TWICE IN ONE DAY
//  10A has Math at Mon-09:00 AND Mon-14:00
//  Expected: hard=0, soft=3.0
// ─────────────────────────────────────────────

func TestSameSubjectTwiceInOneDay(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()

	reqA1 := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqA2 := &model.Requirement{
		ID:              1001,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA1, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // 10A | Math | Mon-P1
		{
			Requirement: reqA2, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "14:00"},
		}, // 10A | Math | Mon-P5  ← same subject, same day
	}

	hard := HardViolations(assignments)
	soft := SoftViolations(assignments)

	fmt.Printf("\n[Case 6] Same subject twice in one day\n")
	fmt.Printf("  Hard violations: %d (want 0)\n", hard)
	fmt.Printf("  Soft penalty   : %.1f (want 3.0)\n", soft)

	if hard != 0 {
		t.Errorf("expected 0 hard violations, got %d", hard)
	}
	if soft != 3.0 {
		t.Errorf("expected 3.0 soft penalty, got %.1f", soft)
	}
}

// ─────────────────────────────────────────────
//  CASE 7 — SOFT: TEACHER GAP (window)
//  Alice teaches Mon-P1 and Mon-P4 with a gap in between
//  Expected: hard=0, soft=2.0
// ─────────────────────────────────────────────

func TestTeacherGap(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, classB, _ := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()

	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classB,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "9:00"},
		}, // Alice | Mon-09:00
		{
			Requirement: reqB, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "13:00"},
		}, // Alice | Mon-13:00  ← gap of 3 periods
	}

	hard := HardViolations(assignments)
	soft := SoftViolations(assignments)

	fmt.Printf("\n[Case 7] Teacher gap (window)\n")
	fmt.Printf("  Hard violations: %d (want 0)\n", hard)
	fmt.Printf("  Soft penalty   : %.1f (want 2.0)\n", soft)

	if hard != 0 {
		t.Errorf("expected 0 hard violations, got %d", hard)
	}
	if soft != 2.0 {
		t.Errorf("expected 2.0 soft penalty, got %.1f", soft)
	}
}

// ─────────────────────────────────────────────
//  CASE 8 — SOFT: HEAVY SUBJECT IN LATE PERIOD
//  Math scheduled at Period 7 (late)
//  Expected: hard=0, soft=1.0
// ─────────────────────────────────────────────

func TestHeavySubjectLatePeriod(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, _, _ := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()

	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "15:00"},
		}, // Math at period 15:00 — heavy subject, late period
	}

	hard := HardViolations(assignments)
	soft := SoftViolations(assignments)

	fmt.Printf("\n[Case 8] Heavy subject in late period\n")
	fmt.Printf("  Hard violations: %d (want 0)\n", hard)
	fmt.Printf("  Soft penalty   : %.1f (want 1.0)\n", soft)

	if hard != 0 {
		t.Errorf("expected 0 hard violations, got %d", hard)
	}
	if soft != 1.0 {
		t.Errorf("expected 1.0 soft penalty, got %.1f", soft)
	}
}

// ─────────────────────────────────────────────
//  CASE 9 — COMBINED: hard + soft violations
//  Room double-booked + teacher gap
//  Expected: hard=1, soft=2.0
// ─────────────────────────────────────────────

func TestCombinedViolations(t *testing.T) {
	alice, bob := makeTeacherFixtures()
	math, english, _ := makeSubjectFixtures()
	classA, classB, classC := makeClassFixtures()
	room101, room102, _ := makeRoomFixtures()

	reqA := &model.Requirement{
		ID:              1001,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1002,
		SchoolClass:     classB,
		Subject:         english,
		Teacher:         bob,
		SessionsPerWeek: 1,
	}
	reqC := &model.Requirement{
		ID:              1003,
		SchoolClass:     classC,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Alice | Mon-09:00" | Room-101 | 10A
		{
			Requirement: reqB, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Bob   | Mon-09:00 | Room-101 | 10B  ← room clash
		{
			Requirement: reqC, Room: room102,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "13:00"},
		}, // Alice | Mon-13:00 | Room-102 | 11A  ← gap from P1
	}

	hard := HardViolations(assignments)
	soft := SoftViolations(assignments)

	fmt.Printf("\n[Case 9] Combined: room conflict + teacher gap\n")
	fmt.Printf("  Hard violations: %d (want 1)\n", hard)
	fmt.Printf("  Soft penalty   : %.1f (want 2.0)\n", soft)

	if hard != 1 {
		t.Errorf("expected 1 hard violation, got %d", hard)
	}
	if soft != 2.0 {
		t.Errorf("expected 2.0 soft penalty, got %.1f", soft)
	}
}

// ─────────────────────────────────────────────
//  CASE 10 — TEACHER UNQUALIFIED
//  Alice is assigned to Science classes
//  Expected: hard=1
// ─────────────────────────────────────────────

func TestTeacherUnqualified(t *testing.T) {
	alice, _ := makeTeacherFixtures()
	math, _, science := makeSubjectFixtures()
	classA, _, _ := makeClassFixtures()
	room101, _, _ := makeRoomFixtures()

	// Both requirements use Alice
	reqA := &model.Requirement{
		ID:              1000,
		SchoolClass:     classA,
		Subject:         science,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}
	reqB := &model.Requirement{
		ID:              1001,
		SchoolClass:     classA,
		Subject:         math,
		Teacher:         alice,
		SessionsPerWeek: 1,
	}

	assignments := []*model.Assignment{
		{
			Requirement: reqA, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "09:00"},
		}, // Alice | Science | Mon-09:00 | Room-101 | 10A <- unqualified teaching
		{
			Requirement: reqB, Room: room101,
			Slot: model.TimeSlot{
				Day:       model.Monday,
				StartTime: "10:00"},
		}, // Alice | Math | Mon-09:00 | Room-102 | 10B
	}

	hard := HardViolations(assignments)

	fmt.Printf("\n[Case 10] Teacher unqualified\n")
	fmt.Printf("  Hard violations: %d (want 1)\n", hard)

	if hard != 1 {
		t.Errorf("expected 1 hard violation, got %d", hard)
	}
}
