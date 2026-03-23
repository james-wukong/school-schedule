// Package solver Verifies teacher load vs. their available slots,
// and that a suitable room exists for every requirement (lab-type and capacity checks)
// Phase 1 – FeasibilityCheck() — Verifies teacher load vs. their available slots,
// and that a suitable room exists for every requirement (lab-type and capacity checks).
//
// Phase 2 – GreedyConstruct() — Expands requirements into individual sessions,
// sorts by most-constrained-first, then assigns each to the first valid
// (slot, room) it finds. Conflict tracking uses three map[string]bool tables keyed
// as "entityID_day_period" for O(1) lookup.
//
// Phase 3 – SimulatedAnnealing() — Runs iterations rounds of random neighbour moves
// (swap two slots / relocate one slot / change room) with the standard Boltzmann
// acceptance: e^(-Δcost/T). Uses copyAssignments() to avoid mutating the current best.
//
// Cost function — TotalCost = 1000 × hardViolations + softPenalty.
// Hard violations count teacher/class/room double-bookings.
// Soft penalties cover teacher gap windows (+2), same subject twice in one day (+3),
// and heavy subjects placed in period 6+ (+1).
//
// Phase 4 – PrintTimetable() / PrintTeacherSchedules() — Renders a day×period grid per class
// and a chronologically sorted schedule per teacher.
package solver

import (
	"fmt"
)

// ─────────────────────────────────────────────
//  MAIN
// ─────────────────────────────────────────────

func Schedule() {
	rng := NewFastRNG()

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Printf("║     School Class Scheduler v%.1f      ║\n", 1.0)
	fmt.Println("╚══════════════════════════════════════╝")

	// Initialize requirements, rooms and teachers
	requirements, rooms, _ := Init()

	totalSessions := 0
	for _, r := range requirements {
		totalSessions += r.SessionsPerWeek
	}

	fmt.Println("\nSchool overview:")
	fmt.Printf("  Requirements  : %d\n", len(requirements))
	fmt.Printf("  Total sessions: %d\n", totalSessions)
	fmt.Printf("  Rooms         : %d\n", len(rooms))

	// ── Phase 1: Feasibility ──────────────────
	issues := FeasibilityCheck(requirements, rooms)
	if len(issues) > 0 {
		fmt.Println("\nStopping: feasibility issues found.")
		return
	}

	// ── Phase 2: Greedy Construction ──────────
	fmt.Println("\n=== Phase 2: Greedy Construction ===")
	initial := GreedyConstruct(rng, requirements, rooms)
	fmt.Printf("  Placed %d / %d sessions\n", len(initial), totalSessions)
	fmt.Printf("  Hard violations : %d\n", HardViolations(initial))
	fmt.Printf("  Soft penalty    : %.1f\n", SoftViolations(initial))

	// ── Phase 3: Simulated Annealing ──────────
	optimised := SimulatedAnnealing(
		rng,
		initial,
		rooms,
		650.0,   // initial temperature
		0.998,   // cooling rate
		100_000, // iterations
	)

	// TODO: Save assignments

	// ── Phase 4: Output ───────────────────────
	PrintTimetable(optimised)
	PrintTeacherSchedules(optimised)
	PrintSummary(optimised, totalSessions)
}
