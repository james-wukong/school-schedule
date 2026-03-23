// Package semester defines the semesters entity and related value objects.
// It represents how data looks in the database or business rules.
package semester

import (
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/school"
)

const (
	TimeDateLayout = "2006-01-02"
)

const (
	TermSpring = iota + 1 // 1
	TermSummer            // 2
	TermFall              // 3
	TermWinter            // 4
)

// Semesters represents the semester table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Semesters struct {
	// id is GENERATED ALWAYS. We use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID       int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`
	SchoolID int64 `gorm:"column:school_id;not null" json:"school_id"`

	Year     int `gorm:"not null" json:"year"`
	Semester int `gorm:"not null" json:"semester"`
	// StartDate and EndDate use time.Time.
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`

	// Relationships (Optional but recommended for Eager Loading)
	School *school.Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"-"`
}

type SemesterFilterEntity struct {
	SchoolID *int64
	Year     *int
	Semester *int
	Page     int
	Limit    int
}

func Newsemesters(schoolID int64,
	year, semester int,
) *Semesters {
	return &Semesters{
		SchoolID: schoolID,
		Year:     year,
		Semester: semester,
	}
}
