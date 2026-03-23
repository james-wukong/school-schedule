package solver

// ─────────────────────────────────────────────
//  PHASE 2 — GREEDY CONSTRUCTION
// ─────────────────────────────────────────────
// Expands requirements into individual sessions, sorts by most-constrained-first,
// then assigns each to the first valid (slot, room) it finds. Conflict tracking uses
// three map[string]bool tables keyed as "entityID_day_period" for O(1) lookup

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

func suitableRooms(req *model.Requirement, rooms []*model.Room) []*model.Room {
	var result []*model.Room
	for _, r := range rooms {
		if r.Capacity < req.SchoolClass.StudentCount {
			continue
		}
		if req.Subject.RequiresLab && r.Type != model.Lab {
			continue
		}
		result = append(result, r)
	}
	return result
}

func GreedyConstruct(
	rng *rand.Rand,
	requirements []*model.Requirement,
	rooms []*model.Room,
) []*model.Assignment {
	// Expand requirements into individual sessions
	var sessions []*model.Requirement
	for _, req := range requirements {
		for i := 0; i < req.SessionsPerWeek; i++ {
			sessions = append(sessions, req)
		}
	}

	// Sort: most sessions/week first (most constrained)
	slices.SortStableFunc(sessions, func(a, b *model.Requirement) int {
		// Descending order
		switch {
		case a.SessionsPerWeek > b.SessionsPerWeek:
			return -1
		case a.SessionsPerWeek < b.SessionsPerWeek:
			return 1
		default:
			return 0
		}
	})

	// Shuffle to break ties
	rng.Shuffle(len(sessions), func(i, j int) {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	})

	// Conflict tracking maps: conflictkey → occupied?
	teacherOccupied := make(map[model.ConflictMapKey[model.TeacherID]]bool)
	classOccupied := make(map[model.ConflictMapKey[model.ClassID]]bool)
	roomOccupied := make(map[model.ConflictMapKey[model.RoomID]]bool)

	var assignments []*model.Assignment
	slotOrder := model.AllTimeSlots()

	for _, req := range sessions {
		suitable := suitableRooms(req, rooms)
		if len(suitable) == 0 {
			// fmt.Printf("  ⚠ No room available for %s-%s / %s\n",
			// 	req.SchoolClass.Grade, req.SchoolClass.Class, req.Subject.Name)
			continue
		}

		// Shuffle slot & room order for variety
		rng.Shuffle(len(slotOrder), func(i, j int) {
			slotOrder[i], slotOrder[j] = slotOrder[j], slotOrder[i]
		})
		rng.Shuffle(len(suitable), func(i, j int) {
			suitable[i], suitable[j] = suitable[j], suitable[i]
		})

		placed := false
		for _, slot := range slotOrder {
			if !isTeacherAvailable(req.Teacher, slot) {
				// fmt.Printf("  ⚠ Teacher Occupied %+v at TimeSlot: %+v\n",
				// 	req.Teacher.Name, slot)
				continue
			}
			tk := model.ConflictMapKey[model.TeacherID]{
				ID:   req.Teacher.ID,
				Day:  slot.Day,
				Slot: slot.StartTime,
			}
			ck := model.ConflictMapKey[model.ClassID]{
				ID:   req.SchoolClass.ID,
				Day:  slot.Day,
				Slot: slot.StartTime,
			}

			if teacherOccupied[tk] || classOccupied[ck] {
				// fmt.Printf("  ⚠ Teacher Occupied %+v\n", tk)
				// fmt.Printf("  ⚠ Or Class Occupied %+v\n", ck)
				continue
			}
			// TODO: bind rooms to classes with a switch
			for _, room := range suitable {
				rk := model.ConflictMapKey[model.RoomID]{
					ID:   room.ID,
					Day:  slot.Day,
					Slot: slot.StartTime,
				}
				if roomOccupied[rk] {
					// fmt.Printf("  ⚠ Room Occupied %d at %+v\n",
					// 	room.ID, slot)
					continue
				}
				a := &model.Assignment{Requirement: req, Room: room, Slot: slot}
				assignments = append(assignments, a)
				teacherOccupied[tk] = true
				classOccupied[ck] = true
				roomOccupied[rk] = true
				placed = true
				break
			}
			if placed {
				break
			}
		}

		if !placed {
			fmt.Printf("  ⚠ Could not place: %s-%s | %s | %s\n",
				req.SchoolClass.Grade, req.SchoolClass.Class,
				req.Subject.Name,
				req.Teacher.Name,
			)
		}
	}

	return assignments
}
