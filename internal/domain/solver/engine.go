package solver

import (
	"fmt"
)

// Schedule runs the CSP solver with local search
func (s *Scheduler) Solve() error {
	fmt.Println("Starting CSP-based scheduling for large school...")

	// Step 1: Initialize domains (possible assignments for each class)
	if err := s.initializeDomains(); err != nil {
		return err
	}
	fmt.Printf("Initialized domains for %d classes\n", len(s.Classes))

	// Step 2: Constraint propagation (AC-3 algorithm)
	if err := s.constraintPropagation(); err != nil {
		return err
	}
	fmt.Println("Constraint propagation completed")

	// Step 3: Backtracking search with MRV (Minimum Remaining Values) heuristic
	if err := s.backtrackingSearch(); err != nil {
		return err
	}
	fmt.Printf("Initial assignment found with %d conflicts\n", s.countHardViolations())

	// Step 4: Local search (Simulated Annealing) to improve solution
	s.localSearch()
	fmt.Printf("After local search: %d hard violations, %d soft violations\n",
		s.countHardViolations(), s.countSoftViolations())

	return nil
}
