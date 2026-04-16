package redis

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/report"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type reportRedis struct {
	redis *redis.Client
	log   *zerolog.Logger
}

func NewReportRedis(redis *redis.Client, log *zerolog.Logger) report.RedisCache {
	return &reportRedis{redis: redis, log: log}
}

func (c *reportRedis) GetWeeklyClassReport(
	ctx context.Context, semesterID int64, version float64,
) ([]report.WeeklyClassScheduleReport, error) {
	return nil, nil
}
