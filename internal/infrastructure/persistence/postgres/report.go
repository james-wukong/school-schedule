package postgres

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/report"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type reportRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewReportRepository(db *gorm.DB, log *zerolog.Logger) report.Repository {
	return &reportRepository{db: db, log: log}
}

func (r *reportRepository) GetWeeklyClassReport(
	ctx context.Context, semesterID int64, version float64,
) ([]report.WeeklyClassScheduleReport, error) {
	// get report from VIEW: vw_class_weekly_schedule
	var schedules []report.WeeklyClassScheduleReport

	// GORM will execute:
	// SELECT * FROM v_weekly_schedules ORDER BY grade, class_name, start_time ASC
	if err := r.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "grade"}, Desc: false},
			{Column: clause.Column{Name: "class_name"}, Desc: false},
			{Column: clause.Column{Name: "start_time"}, Desc: false},
			{Column: clause.Column{Name: "day_of_week"}, Desc: false},
		}}).
		Find(&schedules, "semester_id = ? AND version = ?", semesterID, version).
		Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *reportRepository) GetWeeklyTeacherReport(
	ctx context.Context, semesterID int64, version float64,
) ([]report.WeeklyTeacherScheduleReport, error) {
	// get report from VIEW: vw_teacher_weekly_schedule
	var schedules []report.WeeklyTeacherScheduleReport

	// GORM will execute:
	// SELECT * FROM v_weekly_schedules ORDER BY grade, class_name, start_time ASC
	if err := r.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "teacher_id"}, Desc: false},
			{Column: clause.Column{Name: "start_time"}, Desc: false},
			{Column: clause.Column{Name: "day_of_week"}, Desc: false},
		}}).
		Find(&schedules, "semester_id = ? AND version = ?", semesterID, version).
		Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *reportRepository) GetMaxDay(
	ctx context.Context, semesterID int64, version float64,
) int {
	var maxDay int
	if err := r.db.WithContext(ctx).
		Model(&report.WeeklyClassScheduleReport{}).
		Select("MAX(day_of_week)").
		Where("semester_id = ? AND version = ?", semesterID, version).
		Row().
		Scan(&maxDay); err != nil {
		return 0
	}
	return maxDay
}
