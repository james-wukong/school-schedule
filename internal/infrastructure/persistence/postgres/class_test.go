package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/class"

	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newClassRepo(t *testing.T) (class.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewClassRepository(gormDB, newLogger())
	return repo, mock
}

// classColumns returns the standard column list used in SELECT expectations.
func classColumns() []string {
	return []string{
		"id", "semester_id", "grade", "class", "student_count",
	}
}

func mockClassRow(mock sqlmock.Sqlmock, s *class.Classes) *sqlmock.Rows {
	return mock.NewRows(classColumns()).
		AddRow(s.ID, s.SemesterID, s.Grade, s.ClassName, s.StudentCount)
}

// sampleClass returns a deterministic Classes fixture.
func sampleClass() *class.Classes {
	return &class.Classes{
		ID:           1000,
		SemesterID:   100,
		Grade:        1,
		ClassName:    "A",
		StudentCount: 35,
	}
}

func sampleClass1() *class.Classes {
	return &class.Classes{
		ID:           1001,
		SemesterID:   100,
		Grade:        1,
		ClassName:    "B",
		StudentCount: 35,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestClassCreate_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	entity := sampleClass()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassCreate_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	entity := sampleClass()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassCreate_DuplicateKey(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	entity := sampleClass()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "class_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "classes"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestClassGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	s := sampleClass()
	sch := sampleSchool()
	sem := sampleSemester()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockClassRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Semester
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE`)).
		WithArgs(sem.ID).
		WillReturnRows(mockSemesterRow(mock, sem))

	// EXPECTATION 3: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetByID(ctx, s.ID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Semester)
	require.NotNil(t, result.Semester.School)

	assert.Equal(t, sch.Name, result.Semester.School.Name)
	assert.Equal(t, sem.Year, result.Semester.Year)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetByID_NotFound(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, class.ErrClassNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetByID_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySemesterID ──────────────────────────────────────────────────────────────────

func TestClassGetBySemesterID_FoundWithPreload(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	s := sampleClass()
	sch := sampleSchool()
	sem := sampleSemester()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE semester_id`)).
		WithArgs(sem.ID).
		WillReturnRows(mockClassRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Semester
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE`)).
		WithArgs(sem.ID).
		WillReturnRows(mockSemesterRow(mock, sem))

	// EXPECTATION 3: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetBySemesterID(ctx, s.SemesterID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result)
	require.NotNil(t, result[0].Semester)
	require.NotNil(t, result[0].Semester.School)
	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, sem.Year, result[0].Semester.Year)
	assert.Equal(t, sch.Name, result[0].Semester.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetBySemesterID_NotFound(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySemesterID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassGetBySemesterID_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySemesterID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestClassUpdate_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	entity := sampleClass()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassUpdate_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()
	entity := sampleClass()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "classes"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestClassDelete_Success(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassDelete_DBError(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClassDelete_NonExistentID(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "classes"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestClassCreate_CancelledContext(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleClass())

	assert.Error(t, err)
}

func TestClassGetByID_CancelledContext(t *testing.T) {
	repo, mock := newClassRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
