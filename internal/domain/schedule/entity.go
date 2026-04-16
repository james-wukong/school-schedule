// Package schedule defines the schedules entity and related value objects.
// It represents how data looks in the database or business rules.
package schedule

import (
	"time"

	"github.com/james-wukong/school-schedule/internal/domain/requirement"
	"github.com/james-wukong/school-schedule/internal/domain/room"
	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
)

// ScheduleStatus represents the lifecycle state of a generated timetable.
type ScheduleStatus string

const (
	StatusDraft     ScheduleStatus = "Draft"
	StatusPublished ScheduleStatus = "Published"
	StatusActive    ScheduleStatus = "Active"
	StatusArchived  ScheduleStatus = "Archived"
)

// Schedules represents a finalized entry in the school's timetable.
// It maps a specific schedule to a physical room and a specific time.
type Schedules struct {
	// Identity column starting at 1000
	ID int64 `gorm:"primaryKey;column:id;default:nextval('schedules_id_seq');<-:false" json:"id"`

	// Foreign Keys
	SchoolID      int64  `gorm:"column:school_id;not null;index:idx_schedules_school" json:"school_id"`
	SemesterID    int64  `gorm:"column:semester_id;not null;uniqueIndex:uq_schedules_room_key;uniqueIndex:uq_schedules_requirement_key;index:idx_schedules_semester" json:"school_id"`
	RequirementID int64  `gorm:"column:requirement_id;not null;uniqueIndex:uq_schedules_requirement_key;index:idx_schedules_requirement" json:"requirement_id"`
	RoomID        *int64 `gorm:"column:room_id;default:null;uniqueIndex:uq_schedules_room_key;index:idx_schedules_room" json:"room_id"`
	TimeslotID    int64  `gorm:"column:timeslot_id;not null;uniqueIndex:uq_schedules_room_key;uniqueIndex:uq_schedules_requirement_key;index:idx_schedules_timeslot" json:"timeslot_id"`

	// Relationships (Belongs To)
	School      *school.Schools           `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE" json:"school,omitempty"`
	Requirement *requirement.Requirements `gorm:"foreignKey:RequirementID;constraint:OnDelete:CASCADE" json:"schedule,omitempty"`
	Room        *room.Rooms               `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"room,omitempty"`
	Timeslot    *timeslot.Timeslots       `gorm:"foreignKey:TimeslotID;constraint:OnDelete:CASCADE" json:"timeslot,omitempty"`

	// Metadata and Status
	// Note: We use the custom ScheduleStatus type for type safety in Go.
	Status  ScheduleStatus `gorm:"column:status;type:schedule_status_enum;default:Draft;index:idx_schedules_status" json:"status"`
	Version float64        `gorm:"column:version;type:numeric(10,2);default:1.00;uniqueIndex:idx_sch_room_time;uniqueIndex:idx_sch_req_time" json:"version"`

	// Audit Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ScheduleFilterEntity struct {
	SchoolID      *int64
	RequirementID *int64
	RoomID        *int64
	TimeslotID    *int64
	Version       *float64
	Status        *ScheduleStatus
	Page          int
	Limit         int
}

func NewSchedules(
	schoolID, requirementID, slotID int64,
	roomID *int64,
	version float64,
	status ScheduleStatus,
) *Schedules {
	return &Schedules{
		SchoolID:      schoolID,
		RequirementID: requirementID,
		RoomID:        roomID,
		TimeslotID:    slotID,
		Version:       version,
		Status:        status,
	}
}
