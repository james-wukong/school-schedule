package repository

import (
	"context"

	"github.com/james-wukong/school-schedule/internal/domain/room"
)

type CachedRoomRepository struct {
	repo  room.Repository
	cache room.RedisCache
}

func NewCachedRoomRepository(repo room.Repository, cache room.RedisCache,
) room.Repository {
	return &CachedRoomRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedRoomRepository) Create(ctx context.Context, room *room.Rooms) error {
	return r.repo.Create(ctx, room)
}

// GetByID retrieves a room by their unique identifier.
// It returns the room or an error if the room is not found.
func (r *CachedRoomRepository) GetByID(ctx context.Context, id int64) (*room.Rooms, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *CachedRoomRepository) GetBySchoolID(
	ctx context.Context, schoolID int64,
) ([]*room.Rooms, error) {
	return r.repo.GetBySchoolID(ctx, schoolID)
}

// Update updates an existing room's information in the repository.
// It returns the updated room or an error if the operation fails.
func (r *CachedRoomRepository) Update(ctx context.Context, room *room.Rooms) error {
	return r.repo.Update(ctx, room)
}

// Delete removes a room from the repository by their unique identifier.
// It returns an error if the operation fails.
func (r *CachedRoomRepository) Delete(ctx context.Context, id int64) error {
	return r.repo.Delete(ctx, id)
}
