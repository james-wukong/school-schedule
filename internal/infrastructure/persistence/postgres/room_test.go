package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/room"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newRoomRepo(t *testing.T) (room.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewRoomRepository(gormDB, newLogger())
	return repo, mock
}

// roomColumns returns the standard column list used in SELECT expectations.
func roomColumns() []string {
	return []string{
		"id", "school_id", "name", "code", "room_type", "capacity", "is_active",
	}
}

func mockRoomRow(mock sqlmock.Sqlmock, s *room.Rooms) *sqlmock.Rows {
	return mock.NewRows(roomColumns()).
		AddRow(s.ID, s.SchoolID, s.Name, s.Code, s.RoomType, s.Capacity, s.IsActive)
}

// sampleRoom returns a deterministic Rooms fixture.
func sampleRoom() *room.Rooms {
	return &room.Rooms{
		ID:       3000,
		SchoolID: 10,
		Name:     "Room1-101",
		Code:     "Room1-101",
		RoomType: room.Regular,
		Capacity: 50,
		IsActive: true,
	}
}

func sampleRoom1() *room.Rooms {
	return &room.Rooms{
		ID:       3001,
		SchoolID: 10,
		Name:     "Room1-102",
		Code:     "Room1-102",
		RoomType: room.Regular,
		Capacity: 50,
		IsActive: true,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestRoomCreate_Success(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	entity := sampleRoom()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomCreate_DBError(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	entity := sampleRoom()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomCreate_DuplicateKey(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	entity := sampleRoom()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "room_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "rooms"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestRoomGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	s := sampleRoom()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockRoomRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetByID(ctx, s.ID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)

	assert.Equal(t, sch.Name, result.School.Name)
	assert.Equal(t, s.RoomType, result.RoomType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomGetByID_NotFound(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, room.ErrRoomNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomGetByID_DBError(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestRoomGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	s := sampleRoom()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE school_id`)).
		WithArgs(sch.ID).
		WillReturnRows(mockRoomRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Semester
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result)
	require.NotNil(t, result[0].School)

	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(roomColumns()))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestRoomUpdate_Success(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	entity := sampleRoom()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomUpdate_DBError(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()
	entity := sampleRoom()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "rooms"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestRoomDelete_Success(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "rooms"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomDelete_DBError(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "rooms"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomDelete_NonExistentID(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "rooms"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestRoomCreate_CancelledContext(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleRoom())

	assert.Error(t, err)
}

func TestRoomGetByID_CancelledContext(t *testing.T) {
	repo, mock := newRoomRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
