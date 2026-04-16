package report

import (
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
	"github.com/james-wukong/school-schedule/internal/types"
)

type ScheduleDetailsReport struct {
	ScheduleID    int64              `gorm:"column:id"`
	SchoolID      int64              `gorm:"column:school_id"`
	SemesterID    int64              `gorm:"column:semester_id"`
	RequirementID int64              `gorm:"column:requirement_id"`
	SubjectID     int64              `gorm:"column:subject_id"`
	TeacherID     int64              `gorm:"column:teacher_id"`
	ClassID       int64              `gorm:"column:class_id"`
	RoomID        *int64             `gorm:"column:room_id"`
	Version       float64            `gorm:"column:version"`
	TeacherName   string             `gorm:"column:teacher_name"`
	SubjectName   string             `gorm:"column:subject_name"`
	Grade         int                `gorm:"column:grade"`
	ClassName     string             `gorm:"column:class_name"`
	RoomName      *string            `gorm:"column:room_name"`
	DayOfWeek     timeslot.DayOfWeek `gorm:"column:day_of_week"`
	StartTime     types.ClockTime    `gorm:"column:start_time"`
	EndTime       types.ClockTime    `gorm:"column:end_time"`
}

func (ScheduleDetailsReport) TableName() string {
	return "vw_schedule_detailed_base"
}

type WeeklyClassScheduleReport struct {
	SemesterID  int64              `gorm:"column:semester_id"`
	ClassID     int64              `gorm:"column:class_id"`
	Version     float64            `gorm:"column:version"`
	Grade       int                `gorm:"column:grade"`
	ClassName   string             `gorm:"column:class_name"`
	TeacherName string             `gorm:"column:teacher_name"`
	SubjectName string             `gorm:"column:subject_name"`
	RoomName    *string            `gorm:"column:room_name"`
	DayOfWeek   timeslot.DayOfWeek `gorm:"column:day_of_week"`
	StartTime   types.ClockTime    `gorm:"column:start_time"`
	EndTime     types.ClockTime    `gorm:"column:end_time"`
}

func (WeeklyClassScheduleReport) TableName() string {
	return "vw_class_weekly_schedule"
}

type WeeklyTeacherScheduleReport struct {
	SemesterID  int64              `gorm:"column:semester_id"`
	TeacherID   int64              `gorm:"column:teacher_id"`
	Version     float64            `gorm:"column:version"`
	Grade       int                `gorm:"column:grade"`
	ClassName   string             `gorm:"column:class_name"`
	TeacherName string             `gorm:"column:teacher_name"`
	SubjectName string             `gorm:"column:subject_name"`
	RoomName    *string            `gorm:"column:room_name"`
	DayOfWeek   timeslot.DayOfWeek `gorm:"column:day_of_week"`
	StartTime   types.ClockTime    `gorm:"column:start_time"`
	EndTime     types.ClockTime    `gorm:"column:end_time"`
}

func (WeeklyTeacherScheduleReport) TableName() string {
	return "vw_teacher_weekly_schedule"
}
