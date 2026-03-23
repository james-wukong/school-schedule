package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	tt "github.com/james-wukong/school-schedule/internal/domain/teacher_timeslot"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newTeacherTimeslotRepo(t *testing.T) (tt.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := NewTeacherTimeslotRepository(gormDB, newLogger())
	return repo, mock
}

// teacherTimeslotColumns returns the standard column list used in SELECT expectations.
func teacherTimeslotColumns() []string {
	return []string{
		"teacher_id", "timeslot_id",
	}
}

func mockTeacherTimeslotRow(mock sqlmock.Sqlmock, s *tt.TeacherTimeslots) *sqlmock.Rows {
	return mock.NewRows(teacherTimeslotColumns()).
		AddRow(s.TeacherID, s.TimeslotID)
}

// sampleTeacherTimeslot returns a deterministic Teachers fixture.
func sampleTeacherTimeslot() *tt.TeacherTimeslots {
	return &tt.TeacherTimeslots{
		TeacherID:  300,
		TimeslotID: 200,
	}
}
func sampleTeacherTimeslot1() *tt.TeacherTimeslots {
	return &tt.TeacherTimeslots{
		TeacherID:  300,
		TimeslotID: 201,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestTeacherTimeslotCreate_Success(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTeacherTimeslot()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_timeslots"`)).
		WithArgs(entity.TeacherID, entity.TimeslotID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotCreate_DBError(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTeacherTimeslot()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_timeslots"`)).
		WithArgs(entity.TeacherID, entity.TimeslotID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotCreate_DuplicateKey(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	entity := sampleTeacherTimeslot()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teacher_code_key"`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_timeslots"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByIDs ──────────────────────────────────────────────────────────────────

func TestTeacherTimeslotGetByIDs_Found(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTeacherTimeslot()
	tch := sampleTeacher()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Teacher query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots" WHERE`)).
		WithArgs(s.TeacherID, s.TimeslotID, 1).
		WillReturnRows(mockTeacherTimeslotRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherRow(mock, tch))

	// EXPECTATION 3: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE "timeslots"`)).
		WithArgs(s.TimeslotID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByIDs(ctx, s)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Guard against nil before accessing nested fields
	require.NotNil(t, result.Teacher)
	require.NotNil(t, result.Timeslot)
	assert.Equal(t, tch.ID, result.Teacher.ID)
	assert.Equal(t, slot.ID, result.Timeslot.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotGetByIDs_NotFound(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTeacherTimeslot()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots"`)).
		WithArgs(s.TeacherID, s.TimeslotID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByIDs(ctx, s)

	// Per implementation: (nil, err) on not-found
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, tt.ErrTeacherTimeslotNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotGetByID_DBError(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTeacherTimeslot()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots"`)).
		WithArgs(s.TeacherID, s.TimeslotID, 1).
		WillReturnError(dbErr)

	result, err := repo.GetByIDs(ctx, s)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByTeacherID ──────────────────────────────────────────────────────────────────

func TestTeacherTimeslotGetByTeacherID_Found(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()
	s := sampleTeacherTimeslot()
	tch := sampleTeacher()
	slot := sampleTimeslot()

	// EXPECTATION 1: The primary Teacher Timeslot query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots" WHERE`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherTimeslotRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherRow(mock, tch))

	// EXPECTATION 3: The automatic Preload query for the Timeslot
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "timeslots" WHERE "timeslots"`)).
		WithArgs(s.TimeslotID).
		WillReturnRows(mockTimeslotRow(mock, slot))

	result, err := repo.GetByTeacherID(ctx, s.TeacherID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.TimeslotID, result[0].TimeslotID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotGetByTeacherID_NotFound(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(teacherTimeslotColumns()))

	result, err := repo.GetByTeacherID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherTimeslotGetByTeacherID_DBError(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots"`)).
		WithArgs(int64(300)).
		WillReturnError(dbErr)

	result, err := repo.GetByTimeslotID(ctx, 300)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestTeacherTimeslotCreate_CancelledContext(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleTeacherTimeslot())

	assert.Error(t, err)
}

func TestTeacherTimeslotGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTeacherTimeslotRepo(t)
	entity := sampleTeacherTimeslot()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByIDs(ctx, entity)

	assert.Nil(t, result)
	assert.Error(t, err)
}
