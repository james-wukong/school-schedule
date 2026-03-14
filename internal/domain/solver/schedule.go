// Package solver
package solver

import (
	"sync"

	"github.com/james-wukong/school-schedule/internal/domain/model"
)

type ScheduleStatus string

const (
	ScheduleDraft     ScheduleStatus = "draft"
	SchedulePublished ScheduleStatus = "published"
	ScheduleArchived  ScheduleStatus = "archived"
)

type Scheduler struct {
	ID         int
	TenantID   int
	Classes    map[int]*model.Class
	Teachers   map[int]*model.Teacher
	Classrooms map[int]*model.Classroom
	TimeSlots  []model.TimeSlot
	Schedule   []*ScheduleEntry
	Status     ScheduleStatus

	// Domain for each class: possible assignments
	Domains map[int][]*ScheduleEntry

	// Tracking assignments
	Assignments map[int]*ScheduleEntry
	Constraints *model.ConstraintManager

	mu sync.Mutex
}

// type Schedule struct {
// 	ID        string
// 	TenantID  string
// 	Version   int
// 	Status    ScheduleStatus
// 	Entries   []ScheduleEntry
// 	CreatedAt time.Time
// }

type ScheduleEntry struct {
	Class     *model.Class
	Teacher   *model.Teacher
	Classroom *model.Classroom
	TimeSlot  model.TimeSlot
}

// type ScheduleEntry struct {
// 	ID          string
// 	CourseID    string
// 	TeacherID   string
// 	GroupID     string
// 	ClassroomID string
// 	TimeSlotID  string
// }

// NewScheduler creates a new scheduler
func NewScheduler(maxIterations int) *Scheduler {
	return &Scheduler{
		Classes:     make(map[int]*model.Class),
		Teachers:    make(map[int]*model.Teacher),
		Classrooms:  make(map[int]*model.Classroom),
		TimeSlots:   make([]model.TimeSlot, 0),
		Schedule:    make([]*ScheduleEntry, 0),
		Status:      ScheduleDraft,
		Domains:     make(map[int][]*ScheduleEntry),
		Assignments: make(map[int]*ScheduleEntry),
		Constraints: &model.ConstraintManager{
			MaxIterations: maxIterations,
			Temperature:   100.0,
			CoolingRate:   0.95,
		},
	}
}

func (s *Scheduler) Publish() {
	s.Status = SchedulePublished
}

func (s *Scheduler) IsPublished() bool {
	return s.Status == SchedulePublished
}

// AddTeacher adds a teacher to the scheduler
func (s *Scheduler) AddTeacher(teacher *model.Teacher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Teachers[teacher.ID] = teacher
}

// AddRoom adds a room to the scheduler
func (s *Scheduler) AddRoom(room *model.Classroom) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Classrooms[room.ID] = room
}

// AddClass adds a class to be scheduled
func (s *Scheduler) AddClass(class *model.Class) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Classes[class.ID] = class
}

// AddTimeSlot adds an available time slot
func (s *Scheduler) AddTimeSlot(slot model.TimeSlot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TimeSlots = append(s.TimeSlots, slot)
}
