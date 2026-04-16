package report

import (
	"context"
	"time"
)

const (
	RedisUserPrefix = "class:"
	RedisUerTTL     = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetWeeklyClassReport(
		ctx context.Context, semesterID int64, version float64,
	) ([]WeeklyClassScheduleReport, error)
}
