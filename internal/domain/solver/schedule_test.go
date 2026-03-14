package solver

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

func TestScheduler_AddTeacher(t *testing.T) {
	tests := []struct {
		name    string
		teacher *model.Teacher
		want    bool
	}{
		{
			name: "successfully adds a new class",
			teacher: &model.Teacher{
				ID:   1000,
				Name: "Dr. Smith",
			},
			want: true,
		},
		{
			name: "overwrites existing class with same ID",
			teacher: &model.Teacher{
				ID:   1000,
				Name: "Dr. Johnson",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Teachers: make(map[int]*model.Teacher),
			}
			s.AddTeacher(&model.Teacher{ID: tt.teacher.ID, Name: tt.teacher.Name})
			if got := s.Teachers[tt.teacher.ID]; got.Name != tt.teacher.Name {
				t.Errorf("Scheduler.AddTeacher() = %v, want  %v", got, tt.teacher)
			}
		})
	}
}

func TestScheduler_AddTeacher_Concurrency(t *testing.T) {
	s := &Scheduler{
		Teachers: make(map[int]*model.Teacher),
	}

	const goroutines = 10
	const iterations = 15
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				id := fmt.Sprintf("1000%02d%02d", i, j)
				num, _ := strconv.Atoi(id)
				name := fmt.Sprintf("Teacher %d-%d", i, j)
				s.AddTeacher(&model.Teacher{ID: num, Name: name})
			}
		}(i)
	}

	wg.Wait()

	expectedTotal := goroutines * iterations
	if len(s.Teachers) != expectedTotal {
		t.Errorf("concurrency failure: expected %d classes, got %d", expectedTotal, len(s.Teachers))
	}
}
