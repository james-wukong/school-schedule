// Package teachersubject defines the teacher_subjects entity and related value objects.
// It represents how data looks in the database or business rules.
package teachersubject

import (
	"github.com/james-wukong/school-schedule/internal/domain/subject"
	"github.com/james-wukong/school-schedule/internal/domain/teacher"
)

// TeacherSubjects represents the teacher_subjects table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type TeacherSubjects struct {
	// TeacherID is part of the composite primary key and a foreign key
	TeacherID int64 `gorm:"primaryKey;column:teacher_id;not null" json:"teacher_id"`

	// SubjectID is the other part of the composite primary key and a foreign key
	SubjectID int64 `gorm:"primaryKey;column:subject_id;not null" json:"subject_id"`

	// Optional: Relationships for eager loading the actual entities from the link
	Teacher *teacher.Teachers `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE;" json:"teacher,omitempty"`
	Subject *subject.Subjects `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE;" json:"subject,omitempty"`
}

func NewTeacherSubjects(teacherID, subjectID int64) *TeacherSubjects {
	return &TeacherSubjects{
		TeacherID: teacherID,
		SubjectID: subjectID,
	}
}
