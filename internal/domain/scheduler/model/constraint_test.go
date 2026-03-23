package model

import "testing"

func TestConstraint_IsSoft(t *testing.T) {
	tests := []struct {
		name       string
		constraint Constraint
		want       bool
	}{
		{
			name: "soft constraint",
			constraint: Constraint{
				IsHard: false,
			},
			want: true,
		},
		{
			name: "hard constraint",
			constraint: Constraint{
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
