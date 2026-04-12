package repository

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
)

type CachedTimeslotRepository struct {
	repo  timeslot.Repository
	cache timeslot.RedisCache
}

func NewCachedTimeslotRepository(repo timeslot.Repository, cache timeslot.RedisCache,
) timeslot.Repository {
	return &CachedTimeslotRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedTimeslotRepository) Create(ctx context.Context, slot *timeslot.Timeslots) error {
	return nil
}

func (r *CachedTimeslotRepository) GetByID(ctx context.Context, id int64,
) (*timeslot.Timeslots, error) {
	return nil, nil
}

func (r *CachedTimeslotRepository) GetBySemesterID(ctx context.Context, semesterID int64,
) ([]*timeslot.Timeslots, error) {
	return r.repo.GetBySemesterID(ctx, semesterID)
}

func (r *CachedTimeslotRepository) Update(ctx context.Context, timeslot *timeslot.Timeslots,
) error {
	return nil
}

func (r *CachedTimeslotRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *CachedTimeslotRepository) List(ctx context.Context, filter *timeslot.TimeslotFilterEntity,
) ([]*timeslot.Timeslots, error) {
	return nil, nil
}
