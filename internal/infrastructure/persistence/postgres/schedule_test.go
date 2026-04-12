package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/schedule"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newScheduleRepo(t *testing.T) (schedule.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewScheduleRepository(gormDB, newLogger())
	return repo, mock
}

// scheduleColumns returns the standard column list used in SELECT expectations.
func scheduleColumns() []string {
	return []string{
		"id", "school_id", "requirement_id", "room_id", "timeslot_id",
		"version", "status",
	}
}

func mockScheduleRow(mock sqlmock.Sqlmock, s *schedule.Schedules) *sqlmock.Rows {
	return mock.NewRows(scheduleColumns()).
		AddRow(s.ID, s.SchoolID, s.RequirementID, s.RoomID, s.TimeslotID,
			s.Version, s.Status,
		)
}

// sampleSchedule returns a deterministic Schedules fixture.
func sampleSchedule() *schedule.Schedules {
	return &schedule.Schedules{
		ID:            10000,
		SchoolID:      10,
		RequirementID: 10000,
		RoomID:        3000,
		TimeslotID:    200,
		Version:       1.00,
		Status:        schedule.StatusDraft,
	}
}

func sampleSchedule1() *schedule.Schedules {
	return &schedule.Schedules{
		ID:            10001,
		SchoolID:      10,
		RequirementID: 10000,
		RoomID:        3000,
		TimeslotID:    201,
		Version:       1.00,
		Status:        schedule.StatusDraft,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestScheduleCreate_Success(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	entity := sampleSchedule()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedules"`)).
		WillReturnRows(mockScheduleRow(mock, entity))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	require.NoError(t, err)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleCreate_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	entity := sampleSchedule()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedules"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleCreate_DuplicateKey(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	entity := sampleSchedule()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teacher_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedules"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestScheduleGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()
	sch := sampleSchool()
	req := sampleRequirement()
	room := sampleRoom()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Schedule query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockScheduleRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Requirement
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE `)).
		WithArgs(req.ID).
		WillReturnRows(mockRequirementRow(mock, req))

	// EXPECTATION 3: The automatic Preload query for the Room
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE`)).
		WithArgs(room.ID).
		WillReturnRows(mockRoomRow(mock, room))

	// EXPECTATION 4: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 5: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE `)).
		WithArgs(slot.ID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)
	require.NotNil(t, result.Requirement)
	require.NotNil(t, result.Room)
	require.NotNil(t, result.Timeslot)

	assert.Equal(t, s.ID, result.ID)

	// verify relationships
	assert.Equal(t, sch.Name, result.School.Name)
	assert.Equal(t, req.MinDayGap, result.Requirement.MinDayGap)
	assert.Equal(t, slot.DayOfWeek, result.Timeslot.DayOfWeek)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleGetByID_NotFound(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, err.NotFound) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, schedule.ErrScheduleNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleGetByID_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestScheduleGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()
	sch := sampleSchool()
	req := sampleRequirement()
	room := sampleRoom()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Schedule query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE`)).
		WithArgs(s.SchoolID).
		WillReturnRows(mockScheduleRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Requirement
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE `)).
		WithArgs(req.ID).
		WillReturnRows(mockRequirementRow(mock, req))

	// EXPECTATION 3: The automatic Preload query for the Room
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE`)).
		WithArgs(room.ID).
		WillReturnRows(mockRoomRow(mock, room))

	// EXPECTATION 4: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 5: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE `)).
		WithArgs(slot.ID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.NotNil(t, result[0].School)
	require.NotNil(t, result[0].Requirement)
	require.NotNil(t, result[0].Room)
	require.NotNil(t, result[0].Timeslot)

	assert.Equal(t, s.ID, result[0].ID)

	// verify relationships
	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.Equal(t, req.MinDayGap, result[0].Requirement.MinDayGap)
	assert.Equal(t, slot.DayOfWeek, result[0].Timeslot.DayOfWeek)
}

func TestScheduleGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByVersion ────────────────────────────────────────────────────────────────

func TestScheduleGetByVersion_Found(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()
	sch := sampleSchool()
	req := sampleRequirement()
	room := sampleRoom()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Schedule query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE`)).
		WithArgs(s.SchoolID, s.Version).
		WillReturnRows(mockScheduleRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Requirement
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE `)).
		WithArgs(req.ID).
		WillReturnRows(mockRequirementRow(mock, req))

	// EXPECTATION 3: The automatic Preload query for the Room
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rooms" WHERE`)).
		WithArgs(room.ID).
		WillReturnRows(mockRoomRow(mock, room))

	// EXPECTATION 4: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 5: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE `)).
		WithArgs(slot.ID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByVersion(ctx, s.SchoolID, s.Version)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.NotNil(t, result[0].School)
	require.NotNil(t, result[0].Requirement)
	require.NotNil(t, result[0].Room)
	require.NotNil(t, result[0].Timeslot)

	assert.Equal(t, s.ID, result[0].ID)

	// verify relationships
	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.Equal(t, req.MinDayGap, result[0].Requirement.MinDayGap)
	assert.Equal(t, slot.DayOfWeek, result[0].Timeslot.DayOfWeek)
}

func TestScheduleGetByVersion_NotFound(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(999), float64(11.11)).
		WillReturnRows(sqlmock.NewRows(scheduleColumns()))

	result, err := repo.GetByVersion(ctx, int64(999), float64(11.11))

	// Per implementation: wraps into domain sentinel
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleGetByVersion_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	dbErr := errors.New("db unavailable")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(int64(999), float64(11.11)).
		WillReturnError(dbErr)

	result, err := repo.GetByVersion(ctx, int64(999), float64(11.11))

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestScheduleUpdate_Success(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	entity := sampleSchedule()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schedules"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleUpdate_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	entity := sampleSchedule()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schedules"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestScheduleDelete_Success(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedules"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleDelete_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedules"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleDelete_NonExistentID(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedules"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseScheduleFilter() *schedule.ScheduleFilterEntity {
	return &schedule.ScheduleFilterEntity{Page: 1, Limit: 10}
}

func TestScheduleList_NoFilter(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WillReturnRows(mockScheduleRow(mock, s))

	results, err := repo.List(ctx, baseScheduleFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_FilterBySchoolID(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()

	filter := baseScheduleFilter()
	filter.SchoolID = ptr(s.SchoolID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE school_id = $1`)).
		WithArgs(filter.SchoolID, filter.Limit).
		WillReturnRows(mockScheduleRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_FilterByVersion(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()
	s := sampleSchedule()

	filter := baseScheduleFilter()
	filter.Version = ptr(float64(1.00))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE version`)).
		WithArgs(float64(1.00), filter.Limit).
		WillReturnRows(mockScheduleRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_CombinedFilters(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	filter := baseScheduleFilter()
	filter.SchoolID = ptr(int64(10))
	filter.Version = ptr(float64(1.00))

	// mock.ExpectQuery(`SELECT \* FROM "schedules" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules" WHERE`)).
		WillReturnRows(sqlmock.NewRows(scheduleColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	filter := baseScheduleFilter()
	filter.SchoolID = ptr(int64(0))    // Should not be skipped
	filter.Version = ptr(float64(0.0)) // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WillReturnRows(sqlmock.NewRows(scheduleColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_Pagination(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	filter := &schedule.ScheduleFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(scheduleColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestScheduleList_DBError(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseScheduleFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestScheduleCreate_CancelledContext(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleSchedule())

	assert.Error(t, err)
}

func TestScheduleGetByID_CancelledContext(t *testing.T) {
	repo, mock := newScheduleRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedules"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
