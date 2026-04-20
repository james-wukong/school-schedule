package repository

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/report"
)

type CachedReportRepository struct {
	repo  report.Repository
	cache report.RedisCache
}

func NewCachedReportRepository(
	repo report.Repository, cache report.RedisCache,
) *CachedReportRepository {
	return &CachedReportRepository{
		repo:  repo,
		cache: cache,
	}
}

func (c *CachedReportRepository) GetWeeklyClassReport(
	ctx context.Context, semesterID int64, version float64,
) ([]report.WeeklyClassScheduleReport, error) {
	return c.repo.GetWeeklyClassReport(ctx, semesterID, version)
}

func (c *CachedReportRepository) GetWeeklyTeacherReport(
	ctx context.Context, semesterID int64, version float64,
) ([]report.WeeklyTeacherScheduleReport, error) {
	return c.repo.GetWeeklyTeacherReport(ctx, semesterID, version)
}

func (c *CachedReportRepository) GetMaxDay(
	ctx context.Context, semesterID int64, version float64,
) int {
	return c.repo.GetMaxDay(ctx, semesterID, version)
}
