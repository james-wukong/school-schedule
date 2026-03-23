package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	req "github.com/james-wukong/school-schedule/internal/domain/requirement"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newRequirementRepo(t *testing.T) (req.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := NewRequirementRepository(gormDB, newLogger())
	return repo, mock
}

// requirementColumns returns the standard column list used in SELECT expectations.
func requirementColumns() []string {
	return []string{
		"id", "school_id", "subject_id", "teacher_id", "class_id", "weekly_sessions",
		"min_day_gap", "preferred_days",
	}
}

func mockRequirementRow(mock sqlmock.Sqlmock, s *req.Requirements) *sqlmock.Rows {
	return mock.NewRows(requirementColumns()).
		AddRow(s.ID, s.SchoolID, s.SubjectID, s.TeacherID, s.ClassID,
			s.WeeklySessions, s.MinDayGap, s.PreferredDays,
		)
}

// sampleRequirement returns a deterministic Requirements fixture.
func sampleRequirement() *req.Requirements {
	return &req.Requirements{
		ID:             10000,
		SchoolID:       10,
		SubjectID:      1000,
		TeacherID:      300,
		ClassID:        1000,
		WeeklySessions: 10,
		MinDayGap:      0,
		PreferredDays:  "1,2,3,4,5",
	}
}

func sampleRequirement1() *req.Requirements {
	return &req.Requirements{
		ID:             10001,
		SchoolID:       10,
		SubjectID:      1000,
		TeacherID:      300,
		ClassID:        10001,
		WeeklySessions: 10,
		MinDayGap:      0,
		PreferredDays:  "1,2,3,4,5",
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestRequirementCreate_Success(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	entity := sampleRequirement()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "requirements"`)).
		WillReturnRows(mockRequirementRow(mock, entity))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementCreate_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	entity := sampleRequirement()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "requirements"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementCreate_DuplicateKey(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	entity := sampleRequirement()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teacher_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "requirements"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestRequirementGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()
	sch := sampleSchool()
	sub := sampleSubject()
	tch := sampleTeacher()
	cls := sampleClass()

	// EXPECTATION 1: The primary Requirement query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockRequirementRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE `)).
		WithArgs(cls.ID).
		WillReturnRows(mockClassRow(mock, cls))

	// EXPECTATION 3: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 4: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE `)).
		WithArgs(sub.ID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 5: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.School)
	require.NotNil(t, result.Subject)
	require.NotNil(t, result.Teacher)
	require.NotNil(t, result.Class)

	assert.Equal(t, s.ID, result.ID)

	// verify relationships
	assert.Equal(t, sch.Name, result.School.Name)
	assert.Equal(t, cls.ClassName, result.Class.ClassName)
	assert.Equal(t, tch.FirstName, result.Teacher.FirstName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetByID_NotFound(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, err.NotFound) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, req.ErrRequirementNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetByID_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestRequirementGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()
	sch := sampleSchool()
	sub := sampleSubject()
	tch := sampleTeacher()
	cls := sampleClass()

	// EXPECTATION 1: The primary Requirement query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WithArgs(s.SchoolID).
		WillReturnRows(mockRequirementRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE `)).
		WithArgs(cls.ID).
		WillReturnRows(mockClassRow(mock, cls))

	// EXPECTATION 3: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 4: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE `)).
		WithArgs(sub.ID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 5: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result)
	assert.Len(t, result, 1)

	require.NotNil(t, result[0].School)
	require.NotNil(t, result[0].Subject)
	require.NotNil(t, result[0].Teacher)
	require.NotNil(t, result[0].Class)

	assert.Equal(t, s.ID, result[0].ID)

	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.Equal(t, cls.ClassName, result[0].Class.ClassName)
	assert.Equal(t, tch.FirstName, result[0].Teacher.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySubjectID ──────────────────────────────────────────────────────────────────

func TestRequirementGetBySubjectID_FoundWithPreload(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()
	sch := sampleSchool()
	sub := sampleSubject()
	tch := sampleTeacher()
	cls := sampleClass()

	// EXPECTATION 1: The primary Requirement query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WithArgs(s.SubjectID).
		WillReturnRows(mockRequirementRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE `)).
		WithArgs(cls.ID).
		WillReturnRows(mockClassRow(mock, cls))

	// EXPECTATION 3: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 4: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE `)).
		WithArgs(sub.ID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 5: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetBySubjectID(ctx, s.SubjectID)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result)
	assert.Len(t, result, 1)

	require.NotNil(t, result[0].School)
	require.NotNil(t, result[0].Subject)
	require.NotNil(t, result[0].Teacher)
	require.NotNil(t, result[0].Class)

	assert.Equal(t, s.ID, result[0].ID)

	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.Equal(t, cls.ClassName, result[0].Class.ClassName)
	assert.Equal(t, tch.FirstName, result[0].Teacher.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetBySubjectID_NotFound(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	result, err := repo.GetBySubjectID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetBySubjectID_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySubjectID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByVersion ────────────────────────────────────────────────────────────────

func TestRequirementGetByVersion_Found(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()
	sch := sampleSchool()
	sub := sampleSubject()
	tch := sampleTeacher()
	cls := sampleClass()

	// EXPECTATION 1: The primary Requirement query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WithArgs(s.SchoolID, s.Version).
		WillReturnRows(mockRequirementRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE `)).
		WithArgs(cls.ID).
		WillReturnRows(mockClassRow(mock, cls))

	// EXPECTATION 3: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 4: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE `)).
		WithArgs(sub.ID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 5: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetByVersion(ctx, s.SchoolID, s.Version)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result)
	assert.Len(t, result, 1)

	require.NotNil(t, result[0].School)
	require.NotNil(t, result[0].Subject)
	require.NotNil(t, result[0].Teacher)
	require.NotNil(t, result[0].Class)

	assert.Equal(t, s.ID, result[0].ID)

	assert.Equal(t, sch.Name, result[0].School.Name)
	assert.Equal(t, cls.ClassName, result[0].Class.ClassName)
	assert.Equal(t, tch.FirstName, result[0].Teacher.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetByVersion_NotFound(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999), float64(11.11)).
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	result, err := repo.GetByVersion(ctx, int64(999), float64(11.11))

	// Per implementation: wraps into domain sentinel
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetByVersion_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("db unavailable")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999), float64(11.11)).
		WillReturnError(dbErr)

	result, err := repo.GetByVersion(ctx, int64(999), float64(11.11))

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestRequirementUpdate_Success(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	entity := sampleRequirement()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "requirements"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementUpdate_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	entity := sampleRequirement()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "requirements"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestRequirementDelete_Success(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "requirements"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementDelete_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "requirements"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementDelete_NonExistentID(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "requirements"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseRequirementFilter() *req.RequirementFilterEntity {
	return &req.RequirementFilterEntity{Page: 1, Limit: 10}
}

func TestRequirementList_NoFilter(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WillReturnRows(mockRequirementRow(mock, s))

	results, err := repo.List(ctx, baseRequirementFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_FilterBySchoolID(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()

	filter := baseRequirementFilter()
	filter.SchoolID = ptr(int64(10))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE shcool_id = $1`)).
		WithArgs(filter.SchoolID, filter.Limit).
		WillReturnRows(mockRequirementRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_FilterByVersion(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()

	filter := baseRequirementFilter()
	filter.Version = ptr(float64(1.00))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE LOWER(first_name) LIKE`)).
		WithArgs(float64(1.00), filter.Limit).
		WillReturnRows(mockRequirementRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_CombinedFilters(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	filter := baseRequirementFilter()
	filter.SchoolID = ptr(int64(10))
	filter.Version = ptr(float64(1.00))

	// mock.ExpectQuery(`SELECT \* FROM "requirements" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	filter := baseRequirementFilter()
	filter.SchoolID = ptr(int64(0))    // Should not be skipped
	filter.Version = ptr(float64(0.0)) // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_Pagination(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	filter := &req.RequirementFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementList_DBError(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseRequirementFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestRequirementCreate_CancelledContext(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleRequirement())

	assert.Error(t, err)
}

func TestRequirementGetByID_CancelledContext(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
