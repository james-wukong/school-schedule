package solver

// ─────────────────────────────────────────────
//  PHASE 3 — Optimization: SIMULATED ANNEALING
// ─────────────────────────────────────────────
// Runs iterations rounds of random neighbour moves
// (swap two slots / relocate one slot / change room)
// with the standard Boltzmann acceptance: e^(-Δcost/T).
// Uses copyAssignments() to avoid mutating the current best.
//
// Start with any timetable (even a bad one)
// Repeat:
//     Randomly swap two assignments
//     If better → always accept
//     If worse  → accept with probability e^(-ΔCost / Temperature)
//     Slowly decrease Temperature over time

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

func copyAssignments(src []*model.Assignment) []*model.Assignment {
	dst := make([]*model.Assignment, len(src))
	for i, a := range src {
		cp := *a
		dst[i] = &cp
	}
	return dst
}

// Three neighbour moves: swap slots, move slot, change room
func randomNeighbour(
	rng *rand.Rand,
	assignments []*model.Assignment,
	rooms []*model.Room,
	slots []model.TimeSlot,
	excludeRooms bool,
) []*model.Assignment {
	neighbour := copyAssignments(assignments)
	// TODO ADD THIS
	AllSlots := slots
	n := len(neighbour)
	if n == 0 {
		return neighbour
	}

	move := rng.IntN(3)
	switch move {
	case 0: // Swap slots of two random assignments
		if n < 2 {
			break
		}
		i := rng.IntN(n)
		j := rng.IntN(n)
		for j == i {
			j = rng.IntN(n)
		}
		neighbour[i].Slot, neighbour[j].Slot = neighbour[j].Slot, neighbour[i].Slot

	case 1: // Move one assignment to a random slot
		i := rng.IntN(n)
		neighbour[i].Slot = AllSlots[rng.IntN(len(AllSlots))]

	case 2: // Change room of one assignment
		if !excludeRooms {
			i := rng.IntN(n)
			a := neighbour[i]
			suitable := suitableRooms(a.Requirement, rooms)
			if len(suitable) > 0 {
				neighbour[i].Room = suitable[rng.IntN(len(suitable))]
			}
		}
	}

	return neighbour
}

func SimulatedAnnealing(
	rng *rand.Rand,
	initial []*model.Assignment,
	rooms []*model.Room,
	slots []model.TimeSlot,
	excludeRooms bool,
	initialTemp float64,
	coolingRate float64,
	iterations int,
) []*model.Assignment {
	current := copyAssignments(initial)
	currentCost := TotalCost(current, excludeRooms)
	best := copyAssignments(current)
	bestCost := currentCost
	temp := initialTemp

	fmt.Printf("\n=== Phase 3: Simulated Annealing (iterations=%d) ===\n", iterations)
	fmt.Printf("  Initial cost : %.1f  (hard=%d, soft=%.1f)\n",
		currentCost, HardViolations(current, excludeRooms), SoftViolations(current))

	logInterval := iterations / 10

	for it := 1; it <= iterations; it++ {
		neighbour := randomNeighbour(rng, current, rooms, slots, excludeRooms)
		nCost := TotalCost(neighbour, excludeRooms)
		delta := nCost - currentCost

		if delta < 0 || rng.Float64() < math.Exp(-delta/math.Max(temp, 1e-9)) {
			current = neighbour
			currentCost = nCost
		}

		if currentCost < bestCost {
			best = copyAssignments(current)
			bestCost = currentCost
		}

		temp *= coolingRate

		if it%logInterval == 0 {
			pct := float64(it) / float64(iterations) * 100
			fmt.Printf("  [%5.1f%%] temp=%.4f  best=%.1f  current=%.1f\n",
				pct, temp, bestCost, currentCost)
		}
	}

	fmt.Printf("\n  Final cost   : %.1f  (hard=%d, soft=%.1f)\n",
		bestCost, HardViolations(best, excludeRooms), SoftViolations(best))

	return best
}
