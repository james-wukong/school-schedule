package repository

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/requirement"
)

type CachedRequirementRepository struct {
	repo  requirement.Repository
	cache requirement.RedisCache
}

func NewCachedRequirementRepository(repo requirement.Repository, cache requirement.RedisCache,
) *CachedRequirementRepository {
	return &CachedRequirementRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedRequirementRepository) GetByVersion(ctx context.Context,
	schoolID, semesterID int64,
	version float64,
) ([]*requirement.Requirements, error) {
	return r.repo.GetByVersion(ctx, schoolID, semesterID, version)
}
