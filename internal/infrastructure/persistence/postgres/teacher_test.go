package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/teacher"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newTeacherRepo(t *testing.T) (teacher.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewTeacherRepository(gormDB, newLogger())
	return repo, mock
}

// teacherColumns returns the standard column list used in SELECT expectations.
func teacherColumns() []string {
	return []string{
		"id", "school_id", "employee_id", "first_name", "is_active", "max_classes_per_day",
	}
}

func mockTeacherRow(mock sqlmock.Sqlmock, s *teacher.Teachers) *sqlmock.Rows {
	return mock.NewRows(teacherColumns()).
		AddRow(s.ID, s.SchoolID, s.EmployeeID, s.FirstName, s.IsActive, s.MaxClassesPerDay)
}

// sampleTeacher returns a deterministic Teachers fixture.
func sampleTeacher() *teacher.Teachers {
	return &teacher.Teachers{
		ID:               300,
		SchoolID:         10,
		EmployeeID:       1000,
		FirstName:        "Johnney",
		LastName:         "Zhang",
		IsActive:         true,
		MaxClassesPerDay: 6,
	}
}

func sampleTeacher1() *teacher.Teachers {
	return &teacher.Teachers{
		ID:               301,
		SchoolID:         10,
		EmployeeID:       1001,
		FirstName:        "Depper",
		LastName:         "Lee",
		IsActive:         true,
		MaxClassesPerDay: 6,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestTeacherCreate_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	entity := sampleTeacher()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnRows(mockTeacherRow(mock, entity))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherCreate_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	entity := sampleTeacher()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherCreate_DuplicateKey(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	entity := sampleTeacher()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teacher_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "teachers"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestTeacherGetByID_FoundWithPreload(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()
	sch := sampleSchool()
	sub := sampleSubject()
	sub1 := sampleSubject1()

	// EXPECTATION 1: The primary Teacher query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockTeacherRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 3: GORM queries the join table first
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE `)).
		WithArgs(s.ID).
		WillReturnRows(sqlmock.NewRows([]string{"teacher_id", "subject_id"}).
			AddRow(s.ID, sub.ID).
			AddRow(s.ID, sub1.ID),
		)

	// EXPECTATION 4: GORM then fetches the actual subject rows by collected IDs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"."id" IN ($1,$2)`)).
		WithArgs(sub.ID, sub1.ID).
		WillReturnRows(mockSubjectRow(mock, sub).
			AddRow(sub1.ID, sub1.SchoolID, sub1.Name, sub1.Code,
				sub1.Description, sub1.RequiresLab, sub1.IsHeavy,
			))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.ID, result.ID)
	assert.Equal(t, s.FirstName, result.FirstName)
	assert.Equal(t, sch.Name, result.School.Name)

	// Verify the Many-to-Many data was hydrated
	assert.Len(t, result.Subjects, 2)
	assert.Equal(t, sub.Name, result.Subjects[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByID_NotFound(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, int64(999))

	// Per implementation: (nil, err.NotFound) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, teacher.ErrTeacherNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByID_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetBySchoolID ──────────────────────────────────────────────────────────────────

func TestTeacherGetBySchoolID_FoundWithPreload(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()
	s1 := sampleTeacher1()
	sch := sampleSchool()
	sub := sampleSubject()
	sub1 := sampleSubject1()

	// EXPECTATION 1: The primary Teacher query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE school_id`)).
		WithArgs(s.SchoolID).
		WillReturnRows(mockTeacherRow(mock, s).
			AddRow(
				s1.ID, s1.SchoolID, s1.EmployeeID, s1.FirstName, s1.IsActive, s1.MaxClassesPerDay,
			))

	// EXPECTATION 2: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"."id"`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 3: GORM queries the join table first
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE`)).
		WithArgs(s.ID, s1.ID).
		WillReturnRows(sqlmock.NewRows([]string{"teacher_id", "subject_id"}).
			AddRow(s.ID, sub.ID).
			AddRow(s.ID, sub1.ID).
			AddRow(s1.ID, sub.ID),
		)

	// EXPECTATION 4: GORM then fetches the actual subject rows by collected IDs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"."id" IN ($1,$2)`)).
		WithArgs(sub.ID, sub1.ID).
		WillReturnRows(mockSubjectRow(mock, sub).
			AddRow(sub1.ID, sub1.SchoolID, sub1.Name, sub1.Code,
				sub1.Description, sub1.RequiresLab, sub1.IsHeavy,
			))

	result, err := repo.GetBySchoolID(ctx, s.SchoolID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.FirstName, result[0].FirstName)
	assert.Equal(t, sch.Name, result[0].School.Name)

	assert.Len(t, result[0].Subjects, 2)
	assert.Equal(t, sub.Name, result[0].Subjects[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetBySchoolID_NotFound(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "school_id"}))

	result, err := repo.GetBySchoolID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetBySchoolID_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetBySchoolID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByEmployeeID ──────────────────────────────────────────────────────────────────

func TestTeacherGetByEmployeeID_FoundWithPreload(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()
	s1 := sampleTeacher1()
	sch := sampleSchool()
	sub := sampleSubject()
	sub1 := sampleSubject1()

	// EXPECTATION 1: The primary Teacher query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE employee_id`)).
		WithArgs(s.EmployeeID).
		WillReturnRows(mockTeacherRow(mock, s).
			AddRow(
				s1.ID, s1.SchoolID, s1.EmployeeID, s1.FirstName, s1.IsActive, s1.MaxClassesPerDay,
			))

	// EXPECTATION 2: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"."id"`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 3: GORM queries the join table first
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE`)).
		WithArgs(s.ID, s1.ID).
		WillReturnRows(sqlmock.NewRows([]string{"teacher_id", "subject_id"}).
			AddRow(s.ID, sub.ID).
			AddRow(s.ID, sub1.ID).
			AddRow(s1.ID, sub.ID),
		)

	// EXPECTATION 4: GORM then fetches the actual subject rows by collected IDs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"."id" IN ($1,$2)`)).
		WithArgs(sub.ID, sub1.ID).
		WillReturnRows(mockSubjectRow(mock, sub).
			AddRow(sub1.ID, sub1.SchoolID, sub1.Name, sub1.Code,
				sub1.Description, sub1.RequiresLab, sub1.IsHeavy,
			))

	result, err := repo.GetByEmployeeID(ctx, s.EmployeeID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.FirstName, result[0].FirstName)
	assert.Equal(t, sch.Name, result[0].School.Name)

	assert.Len(t, result[0].Subjects, 2)
	assert.Equal(t, sub.Name, result[0].Subjects[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByEmployeeID_NotFound(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(teacherColumns()))

	result, err := repo.GetByEmployeeID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByEmployeeID_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)

	result, err := repo.GetByEmployeeID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByName ────────────────────────────────────────────────────────────────

func TestTeacherGetByName_Found(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()
	sch := sampleSchool()
	sub := sampleSubject()
	sub1 := sampleSubject1()

	// EXPECTATION 1: The primary Teacher query, Match the actual strings: "%John%" and "zhang"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE LOWER(first_name) LIKE`)).
		WithArgs("%John%", "zhang").
		WillReturnRows(mockTeacherRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE "schools"."id"`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 3: GORM queries the join table first
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE "teacher_subjects"`)).
		WithArgs(s.ID).
		WillReturnRows(sqlmock.NewRows([]string{"teacher_id", "subject_id"}).
			AddRow(s.ID, sub.ID).
			AddRow(s.ID, sub1.ID),
		)

	// EXPECTATION 4: GORM then fetches the actual subject rows by collected IDs
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"."id" IN ($1,$2)`)).
		WithArgs(sub.ID, sub1.ID).
		WillReturnRows(mockSubjectRow(mock, sub).
			AddRow(sub1.ID, sub1.SchoolID, sub1.Name, sub1.Code,
				sub1.Description, sub1.RequiresLab, sub1.IsHeavy,
			))

	result, err := repo.GetByName(ctx, "John", "zhang")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.ID, result[0].ID)
	assert.Equal(t, s.FirstName, result[0].FirstName)
	assert.Equal(t, sch.Name, result[0].School.Name)

	assert.Len(t, result[0].Subjects, 2)
	assert.Equal(t, sub.Name, result[0].Subjects[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByName_NotFound(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs("%MISSING%", "MISS").
		WillReturnRows(sqlmock.NewRows(teacherColumns()))

	result, err := repo.GetByName(ctx, "MISSING", "MISS")

	// Per implementation: wraps into domain sentinel
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherGetByName_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("db unavailable")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs("%John%", "Zhang").
		WillReturnError(dbErr)

	result, err := repo.GetByName(ctx, "John", "Zhang")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestTeacherUpdate_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	entity := sampleTeacher()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	require.NoError(t, err)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherUpdate_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	entity := sampleTeacher()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "teachers"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestTeacherDelete_Success(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "teachers"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherDelete_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "teachers"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherDelete_NonExistentID(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "teachers"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseTeacherFilter() *teacher.TeacherFilterEntity {
	return &teacher.TeacherFilterEntity{Page: 1, Limit: 10}
}

func TestTeacherList_NoFilter(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WillReturnRows(mockTeacherRow(mock, s))

	results, err := repo.List(ctx, baseTeacherFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_FilterByIsActive(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()

	filter := baseTeacherFilter()
	filter.IsActive = ptr(true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE is_active = $1`)).
		WithArgs(true, filter.Limit).
		WillReturnRows(mockTeacherRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_FilterByFirstName(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()
	s := sampleTeacher()

	filter := baseTeacherFilter()
	filter.FirstName = ptr("John")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE LOWER(first_name) LIKE`)).
		WithArgs("%john%", filter.Limit).
		WillReturnRows(mockTeacherRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_CombinedFilters(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	filter := baseTeacherFilter()
	filter.FirstName = ptr("John")
	filter.IsActive = ptr(true)

	// mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE`)).
		WillReturnRows(sqlmock.NewRows(teacherColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	filter := baseTeacherFilter()
	filter.FirstName = ptr("")   // Should not be skipped
	filter.IsActive = ptr(false) // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WillReturnRows(sqlmock.NewRows(teacherColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_Pagination(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	filter := &teacher.TeacherFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(teacherColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherList_DBError(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseTeacherFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestTeacherCreate_CancelledContext(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleTeacher())

	assert.Error(t, err)
}

func TestTeacherGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTeacherRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.Error(t, err)
}
