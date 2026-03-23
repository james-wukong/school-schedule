package solver

import (
	"testing"
)

func TestFeasibilityCheck(t *testing.T) {
	requirements, rooms, _ := Init()
	issues := FeasibilityCheck(requirements, rooms)

	if len(issues) > 0 {
		for _, iss := range issues {
			t.Errorf("issue found %s", iss.Description)
		}
	}
}
