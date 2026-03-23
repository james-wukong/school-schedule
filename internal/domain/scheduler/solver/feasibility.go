package solver

import (
	"fmt"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

// ─────────────────────────────────────────────
//  PHASE 1 — FEASIBILITY CHECK
// ─────────────────────────────────────────────

// FeasibilityCheck performs a quick check to identify any obvious issues
// before we dive into the full scheduling algorithm.
func FeasibilityCheck(requirements []*model.Requirement, rooms []*model.Room) []*model.Issue {
	var issues []*model.Issue

	// Count total sessions per teacher
	teacherLoad := make(map[model.TeacherID]int)
	teacherRef := make(map[model.TeacherID]*model.Teacher)
	for _, req := range requirements {
		teacherLoad[req.Teacher.ID] += req.SessionsPerWeek
		teacherRef[req.Teacher.ID] = req.Teacher
	}

	for tID, load := range teacherLoad {
		teacher := teacherRef[tID]
		avail := 0
		for _, s := range model.AllTimeSlots() {
			if isTeacherAvailable(teacher, s) {
				avail++
			}
		}
		// Checks if the workload is more than the teacher's availability
		if load > avail {
			issues = append(issues, &model.Issue{Description: fmt.Sprintf(
				"Teacher %s needs %d slots but only has %d available",
				teacher.Name, load, avail)})
		}
	}

	// Check suitable rooms exist for each requirement
	for _, req := range requirements {
		neededLab := req.Subject.RequiresLab
		found := false
		for _, room := range rooms {
			if room.Capacity >= req.SchoolClass.StudentCount {
				if !neededLab || room.Type == model.Lab {
					found = true
					break
				}
			}
		}
		if !found {
			issues = append(issues, &model.Issue{Description: fmt.Sprintf(
				"No suitable room for %s-%s / %s (requiresLab=%v)",
				req.SchoolClass.Grade, req.SchoolClass.Class, req.Subject.Name, neededLab)})
		}
	}

	fmt.Println("\n=== Phase 1: Feasibility Check ===")
	if len(issues) == 0 {
		fmt.Println("  ✓ All checks passed")
	} else {
		for _, iss := range issues {
			fmt.Println("  x", iss)
		}
	}
	return issues
}
