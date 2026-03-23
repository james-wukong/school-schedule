// Package requirement defines the requirements entity and related value objects.
// It represents how data looks in the database or business rules.
package requirement

import (
	"github.com/james-wukong/school-schedule/internal/domain/class"
	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/james-wukong/school-schedule/internal/domain/subject"
	"github.com/james-wukong/school-schedule/internal/domain/teacher"
)

// Requirements defines the scheduling constraints for a specific teacher,
// subject, and class combination (e.g., Teacher A must teach Math to Class 10B 3 times a week).
type Requirements struct {
	// Identity column starting at 1000
	ID int64 `gorm:"primaryKey;column:id;default:nextval('requirements_id_seq')" json:"id"`

	// Foreign Keys with Composite Unique Index
	// UNIQUE(subject_id, teacher_id, class_id)
	// Foreign Keys
	SchoolID  int64 `gorm:"column:school_id;not null;index:idx_requirements_school" json:"school_id"`
	SubjectID int64 `gorm:"column:subject_id;not null;index:idx_requirements_subject" json:"subject_id"`
	TeacherID int64 `gorm:"column:teacher_id;not null;index:idx_requirements_teacher" json:"teacher_id"`
	ClassID   int64 `gorm:"column:class_id;not null;index:idx_requirements_class" json:"class_id"`

	// Relationships (Belongs To)
	// These allow GORM to perform Preload("School"), Preload("Subject"), etc.
	School  *school.Schools   `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`
	Subject *subject.Subjects `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE" json:"subject,omitempty"`
	Teacher *teacher.Teachers `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE" json:"teacher,omitempty"`
	Class   *class.Classes    `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class,omitempty"`

	// Scheduling Logic
	WeeklySessions int    `gorm:"column:weekly_sessions;not null;default:1" json:"weekly_sessions"`
	MinDayGap      int    `gorm:"column:min_day_gap;not null;default:0" json:"min_day_gap"`
	PreferredDays  string `gorm:"column:preferred_days;type:varchar(100)" json:"preferred_days"`

	// Versioning (NUMERIC 10,2)
	Version float64 `gorm:"column:version;type:numeric(10,2);default:1.00" json:"version"`
}

type RequirementFilterEntity struct {
	SchoolID  *int64
	SubjectID *int64
	TeacherID *int64
	ClassID   *int64
	Version   *float64
	Page      int
	Limit     int
}

func NewRequirements(
	schoolID, subjectID, teacherID, classID int64,
	weeklySessions, minDayGap int,
	version float64,
) *Requirements {
	return &Requirements{
		SchoolID:       schoolID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
		ClassID:        classID,
		WeeklySessions: weeklySessions,
		MinDayGap:      minDayGap,
		Version:        version,
	}
}
