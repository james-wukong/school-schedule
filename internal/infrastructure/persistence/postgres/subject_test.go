package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/subject"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newSubjectRepo(t *testing.T) (subject.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewSubjectRepository(gormDB, newLogger())
	return repo, mock
}

// subjectColumns returns the standard column list used in SELECT expectations.
func subjectColumns() []string {
	return []string{
		"id", "school_id", "name", "code", "description", "requires_lab", "is_heavy",
	}
}

func mockSubjectRow(mock sqlmock.Sqlmock, s *subject.Subjects) *sqlmock.Rows {
	return mock.NewRows(subjectColumns()).
		AddRow(s.ID, s.SchoolID, s.Name, s.Code, s.Description, s.RequiresLab, s.IsHeavy)
}

// sampleSubject returns a deterministic Subjects fixture.
func sampleSubject() *subject.Subjects {
	return &subject.Subjects{
		ID:          1000,
		SchoolID:    10,
		Name:        "Mathematics",
		Code:        "Math",
		Description: "sample description of a subject",
		RequiresLab: false,
		IsHeavy:     true,
	}
}

func sampleSubject1() *subject.Subjects {
	return &subject.Subjects{
		ID:          1001,
		SchoolID:    10,
		Name:        "Music",
		Code:        "Mus",
		Description: "sample description of a music subject",
		RequiresLab: false,
		IsHeavy:     false,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestSubjectCreate_Success(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	entity := sampleSubject()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectCreate_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	entity := sampleSubject()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectCreate_DuplicateKey(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	entity := sampleSubject()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "subject_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "subjects"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestSubjectGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockSubjectRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// the second argument is $limit
	// mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
	// 	WithArgs(s.ID, 1).
	// 	WillReturnRows(mockSubjectRow(mock, s))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.ID, result.ID)
	assert.Equal(t, s.Name, result.Name)
	assert.Equal(t, sch.Name, result.School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByID_NotFound(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, subject.ErrSubjectNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByID_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestSubjectGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()
	sch := sampleSchool()

	// EXPECTATION 1: The primary Semester query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE`)).
		WithArgs(s.ID).
		WillReturnRows(mockSubjectRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	// GORM uses an "IN" clause for preloading, even for single items.
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.Name, result[0].Name)
	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByCode ────────────────────────────────────────────────────────────────

func TestSubjectGetByCode_Found(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(s.Code, 1).
		WillReturnRows(mockSubjectRow(mock, s))

	result, err := repo.GetByCode(ctx, s.Code)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.Code, result.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByCode_NotFound(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs("MISSING", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByCode(ctx, "MISSING")

	// Per implementation: wraps into domain sentinel
	assert.Nil(t, result)
	assert.ErrorIs(t, err, subject.ErrSubjectNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectGetByCode_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("db unavailable")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs("TS001", 1).
		WillReturnError(dbErr)

	result, err := repo.GetByCode(ctx, "TS001")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestSubjectUpdate_Success(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	entity := sampleSubject()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subjects"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectUpdate_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	entity := sampleSubject()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "subjects"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestSubjectDelete_Success(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectDelete_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectDelete_NonExistentID(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "subjects"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseSubjectFilter() *subject.SubjectFilterEntity {
	return &subject.SubjectFilterEntity{Page: 1, Limit: 10}
}

func TestSubjectList_NoFilter(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WillReturnRows(mockSubjectRow(mock, s))

	results, err := repo.List(ctx, baseSubjectFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_FilterByIsHeavy(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()

	filter := baseSubjectFilter()
	filter.IsHeavy = ptr(true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE is_active = $1`)).
		WithArgs(true, filter.Limit).
		WillReturnRows(mockSubjectRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_FilterByName(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()

	filter := baseSubjectFilter()
	filter.Name = ptr("Test")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE name ILIKE $1`)).
		WithArgs("%Test%", filter.Limit).
		WillReturnRows(mockSubjectRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_FilterByCode(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()
	s := sampleSubject()

	filter := baseSubjectFilter()
	filter.Code = ptr("TS")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE code ILIKE $1`)).
		WithArgs("%TS%", filter.Limit).
		WillReturnRows(mockSubjectRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_CombinedFilters(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	filter := baseSubjectFilter()
	filter.Name = ptr("Test")
	filter.IsHeavy = ptr(true)

	// mock.ExpectQuery(`SELECT \* FROM "subjects" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE`)).
		WillReturnRows(sqlmock.NewRows(subjectColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	filter := baseSubjectFilter()
	filter.Name = ptr("") // Should not be skipped
	filter.Code = ptr("") // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WillReturnRows(sqlmock.NewRows(subjectColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_Pagination(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	filter := &subject.SubjectFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(subjectColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectList_DBError(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseSubjectFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestSubjectCreate_CancelledContext(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleSubject())

	assert.Error(t, err)
}

func TestSubjectGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSubjectRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
