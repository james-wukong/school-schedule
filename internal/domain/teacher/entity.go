// Package teacher defines the teachers entity and related value objects.
// It represents how data looks in the database or business rules.
package teacher

import (
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/james-wukong/school-schedule/internal/domain/subject"
)

// Teachers represents the teachers table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Teachers struct {
	// Primary Key with Identity configuration (START WITH 1000)
	ID int64 `gorm:"primaryKey;autoIncrement:true;autoIncrementIncrement:1;<-:false" json:"id"`

	// Foreign Key to Schools
	SchoolID int64 `gorm:"not null;index:idx_teachers_school;uniqueIndex:idx_teacher_emp_unique" json:"school_id"`

	// Employee ID - Unique per school
	EmployeeID int64 `gorm:"not null;uniqueIndex:idx_teacher_emp_unique" json:"employee_id"`

	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`

	// Indexed for fast lookups during login or directory searches
	Email *string `gorm:"type:varchar(100);index:idx_teachers_email" json:"email"`
	Phone *string `gorm:"type:varchar(20)" json:"phone"`

	HireDate       time.Time `gorm:"type:date;not null" json:"hire_date"`
	EmploymentType string    `gorm:"type:varchar(50)" json:"employment_type"` // 'Full-time', 'Part-time', 'Contract'

	MaxClassesPerDay int  `gorm:"default:5" json:"max_classes_per_day"`
	IsActive         bool `gorm:"not null;default:true;index:idx_teachers_active" json:"is_active"`

	// Standard GORM tracking timestamps
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// --- Relationships ---

	// Many-to-Many relationship with Subjects
	// The 'many2many' tag points to the join table name in PostgreSQL
	// foreignKey: Primary Key of "Source" -> Teachers.ID
	// joinForeignKey: Column in Join Table for Source -> teacher_subjects.teacher_id
	// references: Primary Key of "Target" -> Subjects.ID
	// joinReferences: Column in Join Table for Target -> teacher_subjects.subject_id
	Subjects []subject.Subjects `gorm:"many2many:teacher_subjects;foreignKey:ID;joinForeignKey:TeacherID;References:ID;joinReferences:SubjectID;constraint:OnDelete:CASCADE;" json:"subjects,omitempty"`

	// Belongs To School
	School *school.Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"school,omitempty"`
}

type TeacherFilterEntity struct {
	SchoolID   *int64
	EmployeeID *int64
	FirstName  *string
	LastName   *string
	IsActive   *bool
	Page       int
	Limit      int
}

func NewTeachers(schoolID, employeeID int64,
	firstName, lastName string,
	isActive bool,
) *Teachers {
	return &Teachers{
		SchoolID:   schoolID,
		EmployeeID: employeeID,
		FirstName:  firstName,
		LastName:   lastName,
		IsActive:   isActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
