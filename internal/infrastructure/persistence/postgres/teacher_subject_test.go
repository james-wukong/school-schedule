package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	ts "github.com/james-wukong/school-schedule/internal/domain/teacher_subject"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newTeacherSubjectRepo(t *testing.T) (ts.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewTeacherSubjectRepository(gormDB, newLogger())
	return repo, mock
}

// teacherSubjectColumns returns the standard column list used in SELECT expectations.
func teacherSubjectColumns() []string {
	return []string{
		"teacher_id", "subject_id",
	}
}

func mockTeacherSubjectRow(mock sqlmock.Sqlmock, s *ts.TeacherSubjects) *sqlmock.Rows {
	return mock.NewRows(teacherSubjectColumns()).
		AddRow(s.TeacherID, s.SubjectID)
}

// sampleTeacherSubject returns a deterministic Teachers fixture.
func sampleTeacherSubject() *ts.TeacherSubjects {
	return &ts.TeacherSubjects{
		TeacherID: 300,
		SubjectID: 1000,
	}
}
func sampleTeacherSubject1() *ts.TeacherSubjects {
	return &ts.TeacherSubjects{
		TeacherID: 300,
		SubjectID: 1001,
	}
}

// ── Create ───────────────────────────────────────────────────────────────────

func TestTeacherSubjectCreate_Success(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	entity := sampleTeacherSubject()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_subjects"`)).
		WithArgs(entity.TeacherID, entity.SubjectID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, entity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectCreate_DBError(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	entity := sampleTeacherSubject()

	dbErr := errors.New("connection refused")

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_subjects"`)).
		WithArgs(entity.TeacherID, entity.SubjectID).
		WillReturnError(dbErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectCreate_DuplicateKey(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	entity := sampleTeacherSubject()

	dupErr := errors.New(`ERROR: duplicate key value violates unique constraint "teacher_code_key"`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "teacher_subjects"`)).
		WillReturnError(dupErr)
	mock.ExpectRollback()

	err := repo.Create(ctx, entity)

	assert.Error(t, err)
	assert.ErrorIs(t, err, dupErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByIDs ──────────────────────────────────────────────────────────────────

func TestTeacherSubjectGetByIDs_Found(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	s := sampleTeacherSubject()
	tch := sampleTeacher()
	sub := sampleSubject()

	// EXPECTATION 1: The primary Teacher query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE`)).
		WithArgs(s.TeacherID, s.SubjectID, 1).
		WillReturnRows(mockTeacherSubjectRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Subject
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"`)).
		WithArgs(s.SubjectID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 3: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetByIDs(ctx, s)

	require.NoError(t, err)
	require.NotNil(t, result)
	// Guard against nil before accessing nested fields
	require.NotNil(t, result.Teacher)
	require.NotNil(t, result.Subject)
	assert.Equal(t, tch.ID, result.Teacher.ID)
	assert.Equal(t, sub.ID, result.Subject.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectGetByIDs_NotFound(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	s := sampleTeacherSubject()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects"`)).
		WithArgs(s.TeacherID, s.SubjectID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByIDs(ctx, s)

	// Per implementation: (nil, err) on not-found
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ts.ErrTeacherSubjectNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectGetByID_DBError(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	s := sampleTeacherSubject()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects"`)).
		WithArgs(s.TeacherID, s.SubjectID, 1).
		WillReturnError(dbErr)

	result, err := repo.GetByIDs(ctx, s)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── GetByTeacherID ──────────────────────────────────────────────────────────────────

func TestTeacherSubjectGetByTeacherID_Found(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()
	s := sampleTeacherSubject()
	tch := sampleTeacher()
	sub := sampleSubject()

	// EXPECTATION 1: The primary Teacher Subject query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects" WHERE`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherSubjectRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Subject
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE "subjects"`)).
		WithArgs(s.SubjectID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 3: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE "teachers"`)).
		WithArgs(s.TeacherID).
		WillReturnRows(mockTeacherRow(mock, tch))

	result, err := repo.GetByTeacherID(ctx, s.TeacherID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, s.SubjectID, result[0].SubjectID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectGetByTeacherID_NotFound(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects"`)).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(teacherSubjectColumns()))

	result, err := repo.GetByTeacherID(ctx, int64(999))

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, 0, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTeacherSubjectGetByTeacherID_DBError(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx := context.Background()

	dbErr := errors.New("timeout")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects"`)).
		WithArgs(int64(300)).
		WillReturnError(dbErr)

	result, err := repo.GetBySubjectID(ctx, 300)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ── Context cancellation ─────────────────────────────────────────────────────

func TestTeacherSubjectCreate_CancelledContext(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	mock.ExpectBegin().WillReturnError(context.Canceled)

	err := repo.Create(ctx, sampleTeacherSubject())

	assert.Error(t, err)
}

func TestTeacherSubjectGetByID_CancelledContext(t *testing.T) {
	repo, mock := newTeacherSubjectRepo(t)
	entity := sampleTeacherSubject()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_subjects"`)).
		WillReturnError(context.Canceled)

	result, err := repo.GetByIDs(ctx, entity)

	assert.Nil(t, result)
	assert.Error(t, err)
}
