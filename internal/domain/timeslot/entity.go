// Package timeslot defines the timeslots entity and related value objects.
// It represents how data looks in the database or business rules.
package timeslot

import (
	"time"

	solver "github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/james-wukong/school-schedule/internal/domain/semester"
	"github.com/james-wukong/school-schedule/internal/types"
)

const (
	TimeSlotLayout = "15:04"
)

type DayOfWeek int

const (
	Monday    = iota + 1 // 1
	Tuesday              // 2
	Wednesday            // 3
	Thursday             // 4
	Friday               // 5
)

// Timeslots represents the timeslot table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Timeslots struct {
	// id is GENERATED ALWAYS. use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID         int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`
	SemesterID int64 `gorm:"not null;index:idx_timeslots_semester;index:idx_timeslot_unique,unique" json:"semester_id"`
	// DayOfWeek: 1->Monday, 2->Tuesday, etc.
	DayOfWeek DayOfWeek `gorm:"not null;index:idx_timeslots_day;index:idx_timeslot_unique,unique" json:"day_of_week"`
	// StartDate and EndDate use time.Time.
	StartTime types.ClockTime `gorm:"type:time;not null;index:idx_timeslot_unique,unique" json:"start_time"`
	EndTime   types.ClockTime `gorm:"type:time;not null" json:"end_time"`
	// Foreign Key to School
	SchoolID int64           `gorm:"column:school_id;not null;uniqueIndex:idx_rooms_school_code" json:"school_id"`
	School   *school.Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`

	// Relationships (Optional but recommended for Eager Loading)
	Semester *semester.Semesters `gorm:"foreignKey:SemesterID;constraint:OnDelete:CASCADE;" json:"-"`
}

type TimeslotFilterEntity struct {
	SemesterID *int64
	DayOfWeek  *int
	StartTime  *time.Time
	Page       int
	Limit      int
}

func Newtimeslots(semesterID int64,
	day DayOfWeek,
	start types.ClockTime,
) *Timeslots {
	return &Timeslots{
		SemesterID: semesterID,
		DayOfWeek:  day,
		StartTime:  start,
	}
}

func ToTimeslotMap(ts []*Timeslots) map[solver.DayOfWeek][]string {
	result := make(map[solver.DayOfWeek][]string)
	if len(ts) == 0 {
		return result
	}
	for _, slot := range ts {
		start := time.Time(slot.StartTime).Format(TimeSlotLayout)
		result[solver.DayOfWeek(slot.DayOfWeek)] = append(result[solver.DayOfWeek(slot.DayOfWeek)],
			start)
	}
	return result
}
