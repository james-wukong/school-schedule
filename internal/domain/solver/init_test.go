package solver

import (
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

func TestInit(t *testing.T) {
	requirements, rooms, teachers := Init()
	tests := []struct {
		name         string
		requirements []*model.Requirement
		rooms        []*model.Room
		teachers     []*model.Teacher
		want         bool
	}{
		{
			name:         "checking requirements",
			requirements: requirements,
			rooms:        rooms,
			teachers:     teachers,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.requirements) != 20 {
				t.Error("Expected 20 requirements, got", len(tt.requirements))
			}
			if tt.requirements[19].Subject.ID != 110 {
				t.Errorf("Expected Subject ID: 110, got %d", tt.requirements[19].Subject.ID)
			}
		})
	}
}
