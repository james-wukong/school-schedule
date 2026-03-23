package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	rt "github.com/james-wukong/school-schedule/internal/domain/room_timeslot"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newRoomTimeslotRepo(t *testing.T) (rt.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := NewRoomTimeslotRepository(gormDB, newLogger())
	return repo, mock
}

// roomTimeslotColumns returns the standard column list used in SELECT expectations.
func roomTimeslotColumns() []string {
	return []string{
		"room_id", "timeslot_id",
	}
}

func mockRoomTimeslotRow(mock sqlmock.Sqlmock, s *rt.RoomTimeslots) *sqlmock.Rows {
	return mock.NewRows(roomTimeslotColumns()).
		AddRow(s.RoomID, s.TimeslotID)
}

// sampleRoomTimeslot returns a deterministic Rooms fixture.
func sampleRoomTimeslot() *rt.RoomTimeslots {
	return &rt.RoomTimeslots{
		RoomID:     3000,
		TimeslotID: 200,
	}
}
func sampleRoomTimeslot1() *rt.RoomTimeslots {
	return &rt.RoomTimeslots{
		RoomID:     3000,
		TimeslotID: 201,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestRoomTimeslotCreate_Success(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleRoomTimeslot()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "room_timeslots"`)).
		WithArgs(entity.RoomID, entity.TimeslotID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotCreate_DBError(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleRoomTimeslot()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "room_timeslots"`)).
		WithArgs(entity.RoomID, entity.TimeslotID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotCreate_DuplicateKey(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleRoomTimeslot()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "room_code_key"`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "room_timeslots"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByIDs ──────────────────────────────────────────────────────────────────

func TestRoomTimeslotGetByIDs_Found(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	s := sampleRoomTimeslot()
	tch := sampleRoom()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Room query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots" WHERE`)).
		WithArgs(s.RoomID, s.TimeslotID, 1).
		WillReturnRows(mockRoomTimeslotRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Room
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"`)).
		WithArgs(s.RoomID).
		WillReturnRows(mockRoomRow(mock, tch))

	// EXPECTATION 3: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE "timeslots"`)).
		WithArgs(s.TimeslotID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByIDs(ctx, s)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Guard against nil before accessing nested fields
	require.NotNil(t, result.Room)
	require.NotNil(t, result.Timeslot)
	assert.Equal(t, tch.ID, result.Room.ID)
	assert.Equal(t, slot.ID, result.Timeslot.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotGetByIDs_NotFound(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	s := sampleRoomTimeslot()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots"`)).
		WithArgs(s.RoomID, s.TimeslotID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByIDs(ctx, s)

	// Per implementation: (nil, err) on not-found
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, rt.ErrRoomTimeslotNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotGetByID_DBError(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	s := sampleRoomTimeslot()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots"`)).
		WithArgs(s.RoomID, s.TimeslotID, 1).
		WillReturnError(dbErr)

	result, err := repo.GetByIDs(ctx, s)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByRoomID ──────────────────────────────────────────────────────────────────

func TestRoomTimeslotGetByRoomID_Found(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()
	s := sampleRoomTimeslot()
	tch := sampleRoom()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Room Timeslot query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots" WHERE`)).
		WithArgs(s.RoomID).
		WillReturnRows(mockRoomTimeslotRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Room
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE "rooms"`)).
		WithArgs(s.RoomID).
		WillReturnRows(mockRoomRow(mock, tch))

	// EXPECTATION 3: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE "timeslots"`)).
		WithArgs(s.TimeslotID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByRoomID(ctx, s.RoomID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.TimeslotID, result[0].TimeslotID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotGetByRoomID_NotFound(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(roomTimeslotColumns()))

	result, err := repo.GetByRoomID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoomTimeslotGetByRoomID_DBError(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots"`)).
		WithArgs(int64(300)).
		WillReturnError(dbErr)

	result, err := repo.GetByTimeslotID(ctx, 300)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestRoomTimeslotCreate_CancelledContext(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleRoomTimeslot())

	assert.Error(t, err)
}

func TestRoomTimeslotGetByID_CancelledContext(t *testing.T) {
	repo, mock := newRoomTimeslotRepo(t)
	entity := sampleRoomTimeslot()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "room_timeslots"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByIDs(ctx, entity)

	assert.Nil(t, result)
	assert.Error(t, err)
}
