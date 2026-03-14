package constraint

import "github.com/james-wukong/school-schedule/internal/domain/scheduling"

type SoftConstraintEvaluator interface {
	Score(sol scheduling.Solution) int
}

type DefaultSoftEvaluator struct {
	PreferredTime map[scheduling.VariableID]string
}

func (s *DefaultSoftEvaluator) Score(sol scheduling.Solution) int {
	score := 0

	for varID, assign := range sol {
		// Preferred time bonus
		if preferred, ok := s.PreferredTime[varID]; ok {
			if preferred == assign.TimeSlotID {
				score += 10
			}
		}
	}

	return score
}
