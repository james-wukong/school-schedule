package report

import (
	"context"
)

type Repository interface {
	GetWeeklyClassReport(
		ctx context.Context, semesterID int64, version float64,
	) ([]WeeklyClassScheduleReport, error)

	GetWeeklyTeacherReport(
		ctx context.Context, semesterID int64, version float64,
	) ([]WeeklyTeacherScheduleReport, error)

	GetMaxDay(
		ctx context.Context, semesterID int64, version float64,
	) int
}
