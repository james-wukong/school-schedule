package solver

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

// ─────────────────────────────────────────────
//  PHASE 4 — OUTPUT
// ─────────────────────────────────────────────

// Renders a day×period grid per class and a chronologically sorted schedule per teacher.

func PrintTimetable(assignments []*model.Assignment) {
	// Collect unique classes
	classSet := make(map[model.ClassID]*model.SchoolClass)
	for _, a := range assignments {
		classSet[a.Requirement.SchoolClass.ID] = a.Requirement.SchoolClass
	}
	classes := make([]*model.SchoolClass, 0, len(classSet))
	for _, c := range classSet {
		classes = append(classes, c)
	}

	slices.SortFunc(classes, func(a, b *model.SchoolClass) int {
		switch {
		case int(a.ID) < int(b.ID):
			return -1
		case int(a.ID) > int(b.ID):
			return 1
		default:
			return 0
		}
	})

	sep := strings.Repeat("═", 78)
	line := strings.Repeat("─", 78)

	for _, sc := range classes {
		fmt.Printf("\n%s\n", sep)
		fmt.Printf("  TIMETABLE — %s-%s  (%d students)\n", sc.Grade, sc.Class, sc.StudentCount)
		fmt.Printf("%s\n", sep)

		// Build lookup
		grid := make(map[model.TimeSlot]*model.Assignment)
		for _, a := range assignments {
			if a.Requirement.SchoolClass.ID == sc.ID {
				grid[a.Slot] = a
			}
		}

		// Header
		fmt.Printf("%-15s", "Period")
		days := []model.DayOfWeek{
			model.Monday, model.Tuesday,
			model.Wednesday, model.Thursday, model.Friday,
		}
		for _, d := range days {
			fmt.Printf("%-15d", d)
		}
		fmt.Println()
		fmt.Println(line)
		periods := []string{"09:00", "10:00", "11:00", "13:00", "14:00", "15:00"}
		for _, p := range periods {
			fmt.Printf("%-15s", p)
			for _, d := range days {
				key := model.TimeSlot{
					Day:       d,
					StartTime: p,
				}
				if a, ok := grid[key]; ok {
					subj := a.Requirement.Subject.Name
					if len(subj) > 6 {
						subj = subj[:6]
					}
					teacher := strings.Fields(a.Requirement.Teacher.Name)[0]
					cell := fmt.Sprintf("%s/%s", subj, teacher)
					fmt.Printf("%-15s", cell)
				} else {
					fmt.Printf("%-15s", "—")
				}
			}
			fmt.Println()
		}
	}
	fmt.Printf("\n%s\n", sep)
}

func PrintTeacherSchedules(assignments []*model.Assignment) {
	teacherSet := make(map[model.TeacherID]*model.Teacher)
	for _, a := range assignments {
		teacherSet[a.Requirement.Teacher.ID] = a.Requirement.Teacher
	}
	teachers := make([]*model.Teacher, 0, len(teacherSet))
	for _, t := range teacherSet {
		teachers = append(teachers, t)
	}
	slices.SortFunc(teachers, func(a, b *model.Teacher) int {
		switch {
		case int(a.ID) < int(b.ID):
			return -1
		case int(a.ID) > int(b.ID):
			return 1
		default:
			return 0
		}
	})

	fmt.Println("\n" + strings.Repeat("─", 70))
	fmt.Println("TEACHER SCHEDULES")
	fmt.Println(strings.Repeat("─", 70))

	for _, teacher := range teachers {
		fmt.Printf("\n  Teacher: %s  (subjects: %s)\n",
			teacher.Name, fmt.Sprintf("%+v", teacher.Subjects))

		var slots []*model.Assignment
		for _, a := range assignments {
			if a.Requirement.Teacher.ID == teacher.ID {
				slots = append(slots, a)
			}
		}
		slices.SortFunc(slots, func(a, b *model.Assignment) int {
			// 1. Compare Day Indices
			da := a.Slot.Day
			db := b.Slot.Day

			if da != db {
				return cmp.Compare(da, db)
			}

			// 2. If Days are equal, compare Periods
			return cmp.Compare(a.Slot.StartTime, b.Slot.StartTime)
		})
		for _, a := range slots {
			fmt.Printf("    %+v  →  %s-%s | %s | %s\n",
				a.Slot,
				a.Requirement.SchoolClass.Grade,
				a.Requirement.SchoolClass.Class,
				a.Requirement.Subject.Name,
				a.Room.Name)
		}
	}
}

func PrintSummary(assignments []*model.Assignment, total int) {
	fmt.Println("\n=== Summary ===")
	fmt.Printf("  Total assignments : %d / %d\n", len(assignments), total)
	fmt.Printf("  Hard violations   : %d\n", HardViolations(assignments))
	fmt.Printf("  Soft penalty      : %.1f\n", SoftViolations(assignments))
}
