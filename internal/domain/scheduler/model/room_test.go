package model_test

import (
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

func TestClassroom_CanFit(t *testing.T) {
	tests := []struct {
		classroom model.Room
		name      string
		size      int
		want      bool
	}{
		{
			name: "fits exactly",
			classroom: model.Room{
				Capacity: 30,
			},
			size: 30,
			want: true,
		},
		{
			name: "fits smaller",
			classroom: model.Room{
				Capacity: 30,
			},
			size: 20,
			want: true,
		},
		{
			name: "does not fit larger",
			classroom: model.Room{
				Capacity: 30,
			},
			size: 40,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.classroom.CanFit(tt.size); got != tt.want {
				t.Errorf("Classroom.CanFit() = %v, want %v", got, tt.want)
			}
		})
	}
}
