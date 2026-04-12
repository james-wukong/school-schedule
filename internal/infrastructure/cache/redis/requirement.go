package redis

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/requirement"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type requirementRedis struct {
	redis *redis.Client
}

func NewRequirementRedis(redis *redis.Client, log *zerolog.Logger) requirement.RedisCache {
	return &requirementRedis{redis: redis}
}

// GetByID returns nil if the requirement is not found in cache,
// and returns error only if there is an actual error fetching from cache
// TODO implement some cache
func (r *requirementRedis) GetByVersion(ctx context.Context,
	schoolID, semesterID int64,
	version float64,
) ([]*requirement.Requirements, error) {
	// var res requirement.Requirements
	// // restaurant info hash key format: "restaurant:id"
	// key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, id)
	// // HGetALl return nil if the key doesn't exist.
	// err := r.redis.HGetAll(ctx, key).Scan(&res)
	// if err != nil {
	// 	r.log.Error().Err(err).Msg("Error scanning restaurant from cache")
	// 	return nil, err
	// } else if res.ID == uuid.Nil {
	// 	return nil, nil // Cache miss, return nil without error
	// }
	return nil, nil
}
