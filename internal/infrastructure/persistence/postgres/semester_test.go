package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/semester"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newSemesterRepo(t *testing.T) (semester.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewSemesterRepository(gormDB, newLogger())
	return repo, mock
}

// semesterColumns returns the standard column list used in SELECT expectations.
func semesterColumns() []string {
	return []string{
		"id", "school_id", "year", "semester", "start_date", "end_date",
	}
}

func mockSemesterRow(mock sqlmock.Sqlmock, s *semester.Semesters) *sqlmock.Rows {
	return mock.NewRows(semesterColumns()).
		AddRow(s.ID, s.SchoolID, s.Year, s.Semester, s.StartDate, s.EndDate)
}

// sampleSemester returns a deterministic Semesters fixture.
func sampleSemester() *semester.Semesters {
	start, _ := time.Parse(semester.TimeDateLayout, "2025-01-15")
	end, _ := time.Parse(semester.TimeDateLayout, "2025-01-15")
	return &semester.Semesters{
		ID:        100,
		SchoolID:  10,
		Year:      2024,
		Semester:  3,
		StartDate: start,
		EndDate:   end,
	}
}

func sampleSemester1() *semester.Semesters {
	start, _ := time.Parse(semester.TimeDateLayout, "2025-01-15")
	end, _ := time.Parse(semester.TimeDateLayout, "2025-01-15")
	return &semester.Semesters{
		ID:        101,
		SchoolID:  10,
		Year:      2025,
		Semester:  3,
		StartDate: start,
		EndDate:   end,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestSemesterCreate_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	entity := sampleSemester()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterCreate_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	entity := sampleSemester()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterCreate_DuplicateKey(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	entity := sampleSemester()

	dupErr := errors.New(`ERROR: duplicate key value violates` +
		`unique constraint "semesters_school_id_year_semester_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "semesters"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────
func TestSemesterGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	s := sampleSemester()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE id = $1`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockSemesterRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.ID, result.ID)
	assert.Equal(t, s.Year, result.Year)
	assert.Equal(t, sch.Name, result.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetByID_NotFound(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, semester.ErrSemesterNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetByID_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(100), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestSemesterGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	s := sampleSemester()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE`)).
		WithArgs(s.ID).
		WillReturnRows(mockSemesterRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.Year, result[0].Year)
	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(int64(100)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 100)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestSemesterUpdate_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	entity := sampleSemester()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterUpdate_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	entity := sampleSemester()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "semesters"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestSemesterDelete_Success(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "semesters"`)).
		WithArgs(int64(100)).
		WillReturnResult(sqlmock.NewResult(100, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 100)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterDelete_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "semesters"`)).
		WithArgs(int64(100)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 100)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterDelete_NonExistentID(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "semesters"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseSemesterFilter() *semester.SemesterFilterEntity {
	return &semester.SemesterFilterEntity{Page: 1, Limit: 10}
}

func TestSemesterList_NoFilter(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	s := sampleSemester()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WillReturnRows(mockSemesterRow(mock, s))

	results, err := repo.List(ctx, baseSemesterFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_FilterByYear(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	s := sampleSemester()

	filter := baseSemesterFilter()
	filter.Year = ptr(2025)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE year = $1`)).
		WithArgs(2025, filter.Limit).
		WillReturnRows(mockSemesterRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_FilterBySemester(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()
	s := sampleSemester()

	filter := baseSemesterFilter()
	filter.Semester = ptr(3)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE semester = $1`)).
		WithArgs(3, filter.Limit).
		WillReturnRows(mockSemesterRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_CombinedFilters(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	filter := baseSemesterFilter()
	filter.Year = ptr(2025)
	filter.Semester = ptr(3)

	// mock.ExpectQuery(`SELECT \* FROM "semesters" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE`)).
		WillReturnRows(sqlmock.NewRows(semesterColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	filter := baseSemesterFilter()
	filter.Year = ptr(0)     // Should not be skipped
	filter.Semester = ptr(0) // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WillReturnRows(sqlmock.NewRows(semesterColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_Pagination(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	filter := &semester.SemesterFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(semesterColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSemesterList_DBError(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseSemesterFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestSemesterCreate_CancelledContext(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleSemester())

	assert.Error(t, err)
}

func TestSemesterGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSemesterRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 100)

	assert.Nil(t, result)
	assert.Error(t, err)
}
