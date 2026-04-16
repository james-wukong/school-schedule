package repository

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/schedule"
)

type CachedScheduleRepository struct {
	repo  schedule.Repository
	cache schedule.RedisCache
}

func NewCachedScheduleRepository(
	repo schedule.Repository, cache schedule.RedisCache,
) schedule.Repository {
	return &CachedScheduleRepository{
		repo:  repo,
		cache: cache,
	}
}

func (c *CachedScheduleRepository) Create(ctx context.Context, schedule *schedule.Schedules) error {
	return nil
}

func (c *CachedScheduleRepository) CreateInBatches(
	ctx context.Context, schedule []*schedule.Schedules,
) error {
	return c.repo.CreateInBatches(ctx, schedule)
}

func (c *CachedScheduleRepository) GetByID(
	ctx context.Context, id int64,
) (*schedule.Schedules, error) {
	return nil, nil
}

func (c *CachedScheduleRepository) GetBySchoolID(
	ctx context.Context, schoolID int64,
) ([]*schedule.Schedules, error) {
	return nil, nil
}

func (c *CachedScheduleRepository) GetByRequirementID(
	ctx context.Context, requirementID int64,
) ([]*schedule.Schedules, error) {
	return nil, nil
}

func (c *CachedScheduleRepository) GetByVersion(
	ctx context.Context, schoolID int64, version float64,
) ([]*schedule.Schedules, error) {
	return nil, nil
}

func (c *CachedScheduleRepository) Update(
	ctx context.Context, schedule *schedule.Schedules,
) error {
	return nil
}

func (c *CachedScheduleRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (c *CachedScheduleRepository) List(
	ctx context.Context, filter *schedule.ScheduleFilterEntity,
) ([]*schedule.Schedules, error) {
	return nil, nil
}
