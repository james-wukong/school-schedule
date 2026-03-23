package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/james-wukong/school-schedule/internal/domain/school"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func ptr[T any](v T) *T { return &v }

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })

	return gormDB, mock
}

func newLogger() *zerolog.Logger {
	l := zerolog.Nop()
	return &l
}

func newSchoolRepo(t *testing.T) (school.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := NewSchoolRepository(gormDB, newLogger())
	return repo, mock
}

// schoolColumns returns the standard column list used in SELECT expectations.
func schoolColumns() []string {
	return []string{
		"id", "name", "code", "email", "is_active", "created_at", "updated_at",
	}
}

func mockSchoolRow(mock sqlmock.Sqlmock, s *school.Schools) *sqlmock.Rows {
	return mock.NewRows(schoolColumns()).
		AddRow(s.ID, s.Name, s.Code, s.Email, s.IsActive, s.CreatedAt, s.UpdatedAt)
}

// sampleSchool returns a deterministic Schools fixture.
func sampleSchool() *school.Schools {
	return &school.Schools{
		ID:        10,
		Name:      "Test School",
		Code:      "TS001",
		Email:     "test@school.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func sampleSchool1() *school.Schools {
	return &school.Schools{
		ID:        11,
		Name:      "Another Test School",
		Code:      "TS002",
		Email:     "test1@school.com",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestSchoolCreate_Success(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	entity := sampleSchool()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schools"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(entity.ID))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolCreate_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	entity := sampleSchool()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schools"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolCreate_DuplicateKey(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	entity := sampleSchool()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "schools_code_key"`)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schools"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func TestSchoolGetByID_Found(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	// the second argument is $limit
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(s.ID, 1).
		WillReturnRows(mockSchoolRow(mock, s))

	result, err := repo.GetByID(ctx, s.ID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.ID, result.ID)
	assert.Equal(t, s.Name, result.Name)
	assert.Equal(t, s.Code, result.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByID_NotFound(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(int64(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, 999)

	// Per implementation: (nil, nil) on not-found
	assert.Nil(t, result)
	assert.ErrorIs(t, err, school.ErrSchoolNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByID_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(int64(1), 1).
		WillReturnError(dbErr)

	result, err := repo.GetByID(ctx, 1)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByCode ────────────────────────────────────────────────────────────────

func TestSchoolGetByCode_Found(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(s.Code, 1).
		WillReturnRows(mockSchoolRow(mock, s))

	result, err := repo.GetByCode(ctx, s.Code)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, s.Code, result.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByCode_NotFound(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs("MISSING", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByCode(ctx, "MISSING")

	// Per implementation: wraps into domain sentinel
	assert.Nil(t, result)
	assert.ErrorIs(t, err, school.ErrSchoolNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolGetByCode_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	dbErr := errors.New("db unavailable")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs("TS001", 1).
		WillReturnError(dbErr)

	result, err := repo.GetByCode(ctx, "TS001")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestSchoolUpdate_Success(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	entity := sampleSchool()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schools"`)).
		WillReturnResult(sqlmock.NewResult(entity.ID, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolUpdate_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	entity := sampleSchool()

	dbErr := errors.New("update failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schools"`)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Update(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Delete ───────────────────────────────────────────────────────────────────

func TestSchoolDelete_Success(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schools"`)).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, int64(10))

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolDelete_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	dbErr := errors.New("delete failed")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schools"`)).
		WithArgs(int64(1)).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Delete(ctx, int64(10))

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolDelete_NonExistentID(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	// DELETE on a missing ID succeeds at DB level; GORM does not error.
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schools"`)).
		WithArgs(int64(9999)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(ctx, 9999)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── List ─────────────────────────────────────────────────────────────────────

func baseSchoolFilter() *school.SchoolFilterEntity {
	return &school.SchoolFilterEntity{Page: 1, Limit: 10}
}

func TestSchoolList_NoFilter(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WillReturnRows(mockSchoolRow(mock, s))

	results, err := repo.List(ctx, baseSchoolFilter())

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, s.ID, results[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_FilterByEmail(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	filter := baseSchoolFilter()
	filter.Email = ptr("test@school.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE email = $1`)).
		WithArgs(*filter.Email, filter.Limit).
		WillReturnRows(mockSchoolRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_FilterByIsActive(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	filter := baseSchoolFilter()
	filter.IsActive = ptr(true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE is_active = $1`)).
		WithArgs(true, filter.Limit).
		WillReturnRows(mockSchoolRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_FilterByName(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	filter := baseSchoolFilter()
	filter.Name = ptr("Test")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE name ILIKE $1`)).
		WithArgs("%Test%", filter.Limit).
		WillReturnRows(mockSchoolRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_FilterByCode(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()
	s := sampleSchool()

	filter := baseSchoolFilter()
	filter.Code = ptr("TS")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE code ILIKE $1`)).
		WithArgs("%TS%", filter.Limit).
		WillReturnRows(mockSchoolRow(mock, s))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_CombinedFilters(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	filter := baseSchoolFilter()
	filter.Name = ptr("Test")
	filter.IsActive = ptr(true)

	// mock.ExpectQuery(`SELECT \* FROM "schools" WHERE`).
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WillReturnRows(sqlmock.NewRows(schoolColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_EmptyStringFiltersIgnored(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	filter := baseSchoolFilter()
	filter.Email = ptr("") // Should not be skipped
	filter.Name = ptr("")  // Should not be skipped
	filter.Code = ptr("")  // Should not be skipped

	// Expect a plain SELECT without any WHERE conditions
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WillReturnRows(sqlmock.NewRows(schoolColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_Pagination(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	filter := &school.SchoolFilterEntity{Page: 3, Limit: 5} // offset = (3-1)*5 = 10

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WithArgs(5, 10). // LIMIT 5 OFFSET 10
		WillReturnRows(sqlmock.NewRows(schoolColumns()))

	results, err := repo.List(ctx, filter)

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSchoolList_DBError(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx := context.Background()

	dbErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WillReturnError(dbErr)

	results, err := repo.List(ctx, baseSchoolFilter())

	assert.Nil(t, results)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestSchoolCreate_CancelledContext(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleSchool())

	assert.Error(t, err)
}

func TestSchoolGetByID_CancelledContext(t *testing.T) {
	repo, mock := newSchoolRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByID(ctx, int64(10))

	assert.Nil(t, result)
	assert.Error(t, err)
}
