package solver

import (
	"fmt"
)

// localSearch uses Simulated Annealing to improve solution
func (s *Scheduler) localSearch() {
	temperature := s.Constraints.Temperature

	for iteration := 0; iteration < s.Constraints.MaxIterations; iteration++ {
		// Make random change
		newAssignments := s.copyAssignments()
		classID := s.randomClassID()
		randomEntry := s.Domains[classID][randomInt(len(s.Domains[classID]))]
		newAssignments[classID] = randomEntry

		currentCost := costFunction(s.Assignments)
		newCost := costFunction(newAssignments)
		delta := currentCost - newCost

		// Accept or reject
		if delta > 0 || acceptWithProbability(delta, temperature) {
			s.Assignments = newAssignments
		}

		temperature *= s.Constraints.CoolingRate

		if iteration%100 == 0 {
			fmt.Printf("Iteration %d: cost=%.2f, temp=%.2f\n",
				iteration, float64(costFunction(s.Assignments)), temperature)
		}
	}
}

// randomClassID returns a random class ID from the scheduler's classes
func (s *Scheduler) randomClassID() int {
	i := 0
	idx := randomInt(len(s.Classes))
	for classID := range s.Classes {
		if i == idx {
			return classID
		}
		i++
	}
	return 0
}

// copyAssignments creates a deep copy of the current assignments map
func (s *Scheduler) copyAssignments() map[int]*ScheduleEntry {
	clonedCopy := make(map[int]*ScheduleEntry)
	for k, v := range s.Assignments {
		clonedCopy[k] = v
	}
	return clonedCopy
}
