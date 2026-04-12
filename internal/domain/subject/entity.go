// Package subject defines the subjects entity and related value objects.
// It represents how data looks in the database or business rules.
package subject

import (
	"github.com/james-wukong/school-schedule/internal/domain/school"

	solver "github.com/james-wukong/school-schedule/internal/domain/scheduler/model"
)

// Subjects represents the subjects table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Subjects struct {
	// id is GENERATED ALWAYS. We use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID       int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`
	SchoolID int64 `gorm:"column:school_id;not null" json:"school_id"`

	Name        string `gorm:"column:name;not null;unique" json:"name"`
	Code        string `gorm:"column:code;not null;unique" json:"code"`
	Description string `gorm:"column:description" json:"description"`

	RequiresLab bool `gorm:"column:requires_lab;default:false" json:"requires_lab"`
	IsHeavy     bool `gorm:"column:is_heavy;default:false" json:"is_heavy"`
	// Relationships (Optional but recommended for Eager Loading)
	School *school.Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"-"`
}

type SubjectFilterEntity struct {
	SchoolID *int64
	Name     *string
	Code     *string
	IsHeavy  *bool
	Page     int
	Limit    int
}

func NewSubjects(schoolID int64,
	name, code string,
	requiresLab, isHeavy bool,
) *Subjects {
	return &Subjects{
		SchoolID:    schoolID,
		Name:        name,
		Code:        code,
		RequiresLab: requiresLab,
		IsHeavy:     isHeavy,
	}
}

func (m *Subjects) ToSolverModel() *solver.Subject {
	return &solver.Subject{
		ID:          solver.SubjectID(m.ID),
		Name:        m.Name,
		RequiresLab: m.RequiresLab,
		IsHeavy:     m.IsHeavy,
	}
}
