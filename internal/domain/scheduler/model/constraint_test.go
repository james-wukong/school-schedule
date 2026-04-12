package model_test

import (
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

func TestConstraint_IsSoft(t *testing.T) {
	tests := []struct {
		name       string
		constraint model.Constraint
		want       bool
	}{
		{
			name: "soft constraint",
			constraint: model.Constraint{
				IsHard: false,
			},
			want: true,
		},
		{
			name: "hard constraint",
			constraint: model.Constraint{
				IsHard: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.constraint.IsSoft(); got != tt.want {
				t.Errorf("Constraint.IsSoft() = %v, want %v", got, tt.want)
			}
		})
	}
}
