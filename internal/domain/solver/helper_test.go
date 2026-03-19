package solver

import (
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

func TestIsTeacherAvailable(t *testing.T) {
	tests := []struct {
		name    string
		teacher *model.Teacher
		slot    model.TimeSlot
		want    bool
	}{
		{
			name: "available teacher timeslot",
			teacher: &model.Teacher{
				ID:   1000,
				Name: "Dr. Smith",
				AvailableTimes: map[model.DayOfWeek][]string{
					model.Monday:    {"08:00", "09:00", "10:00", "13:00", "14:00", "15:00"},
					model.Tuesday:   {"08:00", "09:00", "10:00"},
					model.Wednesday: {"10:00", "13:00", "14:00", "15:00"},
					model.Friday:    {"09:00", "10:00", "13:00", "14:00"},
				},
			},
			slot: model.TimeSlot{Day: model.Monday, StartTime: "08:00"},
			want: true,
		},
		{
			name: "unavailable teacher timeslot",
			teacher: &model.Teacher{
				ID:   1000,
				Name: "Dr. Smith",
				AvailableTimes: map[model.DayOfWeek][]string{
					model.Monday:    {"08:00", "09:00", "10:00", "13:00", "14:00", "15:00"},
					model.Tuesday:   {"08:00", "09:00", "10:00"},
					model.Wednesday: {"10:00", "13:00", "14:00", "15:00"},
					model.Friday:    {"09:00", "10:00", "13:00", "14:00"},
				},
			},
			slot: model.TimeSlot{Day: model.Tuesday, StartTime: "13:00"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTeacherAvailable(tt.teacher, tt.slot); got != tt.want {
				t.Errorf("isTeacherAvailable() =%v, want=%v", got, tt.want)
			}
		})
	}
}

func TestIsRoomAvailable(t *testing.T) {
	tests := []struct {
		name      string
		classroom *model.Room
		slot      model.TimeSlot
		want      bool
	}{
		{
			name: "available classroom timeslot",
			classroom: &model.Room{
				ID: 1000,
				AvailableTimes: map[model.DayOfWeek][]string{
					model.Monday: {"08:00", "09:00", "10:00", "13:00", "14:00", "15:00"},
				},
			},
			slot: model.TimeSlot{Day: model.Monday, StartTime: "09:00"},
			want: true,
		},
		{
			name: "unavailable classroom timeslot",
			classroom: &model.Room{
				ID: 1000,
				AvailableTimes: map[model.DayOfWeek][]string{
					model.Friday: {"09:00", "10:00", "13:00", "14:00"},
				},
			},
			slot: model.TimeSlot{Day: model.Tuesday, StartTime: "13:00"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRoomAvailable(tt.classroom, tt.slot); got != tt.want {
				t.Errorf("isRoomAvailableInDomain() =%v, want=%v", got, tt.want)
			}
		})
	}
}

func TestDayDifference(t *testing.T) {
	tests := []struct {
		name  string
		start model.DayOfWeek
		end   model.DayOfWeek
		want  int
	}{
		{
			name:  "same day",
			start: model.Monday,
			end:   model.Monday,
			want:  0,
		},
		{
			name:  "positive day difference",
			start: model.Monday,
			end:   model.Friday,
			want:  4,
		},
		{
			name:  "negative day difference",
			start: model.Tuesday,
			end:   model.Monday,
			want:  10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dayDifference(tt.start, tt.end); got != tt.want {
				t.Errorf("dayDifference() =%v, want=%v", got, tt.want)
			}
		})
	}
}

func TestIsTimeAfter(t *testing.T) {
	tests := []struct {
		name      string
		slot      string
		benchmark string
		want      bool
	}{
		{
			name:      "before benchmark",
			slot:      "12:00",
			benchmark: "13:00",
			want:      false,
		},
		{
			name:      "after benchmark",
			slot:      "12:00",
			benchmark: "11:00",
			want:      true,
		},
		{
			name:      "wrong format",
			slot:      "12:00:00",
			benchmark: "13:00",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isTimeAfter(tt.slot, tt.benchmark)
			if got != tt.want {
				t.Errorf("isTimeAfter() got =%v, want=%v", got, tt.want)
			}
			if err != nil {
				if got != tt.want || got {
					t.Errorf("isTimeAfter() err =%v, got=%v, and want=%v", err, got, tt.want)
				}
			}
		})
	}
}
