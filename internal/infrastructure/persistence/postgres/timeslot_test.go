package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/timeslot"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"github.com/james-wukong/school-schedule/internal/types"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newTimeslotRepo(t *testing.T) (timeslot.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewTimeslotRepository(gormDB, newLogger())
	return repo, mock
}

// timeslotColumns returns the standard column list used in SELECT expectations.
func timeslotColumns() []string {
	return []string{
		"id", "semester_id", "day_of_week", "start_time",
	}
}

func mockTimeslotRow(mock sqlmock.Sqlmock, s *timeslot.Timeslots) *sqlmock.Rows {
	return mock.NewRows(timeslotColumns()).
		AddRow(s.ID, s.SemesterID, s.DayOfWeek, s.StartTime)
}

// sampleTimeslot returns a deterministic Timeslots fixture.
func sampleTimeslot() *timeslot.Timeslots {
	start, _ := time.Parse(timeslot.TimeSlotLayout, "09:00")
	end, _ := time.Parse(timeslot.TimeSlotLayout, "09:45")
	return &timeslot.Timeslots{
		ID:         200,
		SemesterID: 100,
		DayOfWeek:  1,
		StartTime:  types.ClockTime(start),
		EndTime:    types.ClockTime(end),
	}
}

func sampleTimeslot1() *timeslot.Timeslots {
	start, _ := time.Parse(timeslot.TimeSlotLayout, "09:00")
	end, _ := time.Parse(timeslot.TimeSlotLayout, "09:45")
	return &timeslot.Timeslots{
		ID:         201,
		SemesterID: 100,
		DayOfWeek:  2,
		StartTime:  types.ClockTime(start),
		EndTime:    types.ClockTime(end),
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestTimeslotCreate_Success(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTimeslot()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotCreate_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTimeslot()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotCreate_DuplicateKey(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTimeslot()

	dupErr := errors.New(`ERROR: duplicate key value violates` +
		`unique constraint "timeslots_school_id_year_timeslot_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "timeslots"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────
func TestTimeslotGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()
	sch := sampleSchool()
	sem := sampleSemester()

	// EXPECTATION 1: The primary Timeslot query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE id = $1`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockTimeslotRow(mock, s))

		// EXPECTATION 2. Nested Level 1: Semesters (GORM uses IN clause for preloads)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE "semesters"."id" = $1`)).
		WithArgs(sem.ID).
		WillReturnRows(mockSemesterRow(mock, sem))

	// EXPECTATION 3: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.ID, result.ID)
	assert.Equal(t, s.DayOfWeek, result.DayOfWeek)
	assert.Equal(t, sem.Semester, result.Semester.Semester)
	assert.Equal(t, sch.Name, result.Semester.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetByID_NotFound(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, timeslot.ErrTimeslotNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetByID_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySemesterID ──────────────────────────────────────────────────────────────────

func TestTimeslotGetBySemesterID_FoundWithPreload(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()
	sch := sampleSchool()
	sem := sampleSemester()

	// EXPECTATION 1: The primary Timeslot query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE school_id = `)).
		WithArgs(s.ID).
		WillReturnRows(mockTimeslotRow(mock, s))

		// EXPECTATION 2. Nested Level 1: Semesters (GORM uses IN clause for preloads)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE `)).
		WithArgs(sem.ID).
		WillReturnRows(mockSemesterRow(mock, sem))

	// EXPECTATION 3: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE `)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetBySemesterID(ctx, s.SemesterID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.DayOfWeek, result[0].DayOfWeek)
	assert.Equal(t, sem.Semester, result[0].Semester.Semester)
	assert.Equal(t, sch.Name, result[0].Semester.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetBySemesterID_NotFound(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "semester_id"}))

	result, err := repo.GetBySemesterID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotGetBySemesterID_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySemesterID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestTimeslotUpdate_Success(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTimeslot()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "timeslots"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotUpdate_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTimeslot()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "timeslots"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestTimeslotDelete_Success(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotDelete_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotDelete_NonExistentID(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "timeslots"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseTimeslotFilter() *timeslot.TimeslotFilterEntity {
	return &timeslot.TimeslotFilterEntity{Page: 1, Limit: 10}
}

func TestTimeslotList_NoFilter(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WillReturnRows(mockTimeslotRow(mock, s))

	results, err := repo.List(ctx, baseTimeslotFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_FilterByDayOfWeek(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()

	filter := baseTimeslotFilter()
	filter.DayOfWeek = ptr(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE day_of_week = $1`)).
		WithArgs(1, filter.Limit).
		WillReturnRows(mockTimeslotRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_FilterByStartTime(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()

	filter := baseTimeslotFilter()
	filter.StartTime = ptr(time.Time(s.StartTime))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE start_time = $1`)).
		WithArgs(s.StartTime, filter.Limit).
		WillReturnRows(mockTimeslotRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_CombinedFilters(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()

	filter := baseTimeslotFilter()
	filter.DayOfWeek = ptr(1)
	filter.StartTime = ptr(time.Time(s.StartTime))

	// mock.ExpectQuery(`SELECT \* FROM "timeslots" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE`)).
		WillReturnRows(sqlmock.NewRows(timeslotColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTimeslot()

	filter := baseTimeslotFilter()
	filter.DayOfWeek = ptr(0)                      // Should not be skipped
	filter.StartTime = ptr(time.Time(s.StartTime)) // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WillReturnRows(sqlmock.NewRows(timeslotColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_Pagination(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	filter := &timeslot.TimeslotFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(timeslotColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTimeslotList_DBError(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseTimeslotFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestTimeslotCreate_CancelledContext(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleTimeslot())

	assert.Error(t, err)
}

func TestTimeslotGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTimeslotRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 200)

	assert.Nil(t, result)
	assert.Error(t, err)
}
