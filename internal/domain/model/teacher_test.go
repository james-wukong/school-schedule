package model

import "testing"

func TestTeacher_CanTakeMoreHours(t *testing.T) {
	tests := []struct {
		name         string
		teacher      Teacher
		currentHours int
		want         bool
	}{
		{
			name: "Can take more hours",
			teacher: Teacher{
				MaxHoursPerWeek: 10,
			},
			currentHours: 5,
			want:         true,
		},
		{
			name: "Can not take more hours",
			teacher: Teacher{
				MaxHoursPerWeek: 10,
			},
			currentHours: 15,
			want:         false,
		},
		{
			name: "Can nottake more hours",
			teacher: Teacher{
				MaxHoursPerWeek: 10,
			},
			currentHours: 10,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.teacher.CanTakeMoreHours(tt.currentHours); got != tt.want {
				t.Errorf("Teacher.CanTakeMoreHours() = %v, want %v", got, tt.want)
			}
		})
	}
}
