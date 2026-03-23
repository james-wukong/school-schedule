// Package teachertimeslot defines the teacher_subjects entity and related value objects.
// It represents how data looks in the database or business rules.
package teachertimeslot

import (
	"github.com/james-wukong/school-schedule/internal/domain/teacher"
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
)

// TeacherTimeslots represents the teacher_subjects table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type TeacherTimeslots struct {
	// Primary Key composition
	TeacherID  int64 `gorm:"primaryKey;column:teacher_id;not null" json:"teacher_id"`
	TimeslotID int64 `gorm:"primaryKey;column:timeslot_id;not null" json:"timeslot_id"`

	// Relationships (Belongs To)
	// These allow GORM to perform Preload("Teacher") or Preload("Timeslot")
	Teacher  *teacher.Teachers   `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE" json:"teacher,omitempty"`
	Timeslot *timeslot.Timeslots `gorm:"foreignKey:TimeslotID;constraint:OnDelete:CASCADE" json:"timeslot,omitempty"`
}

func NewTeacherTimeslots(teacherID, slotID int64) *TeacherTimeslots {
	return &TeacherTimeslots{
		TeacherID:  teacherID,
		TimeslotID: slotID,
	}
}
