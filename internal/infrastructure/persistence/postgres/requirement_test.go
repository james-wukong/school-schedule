package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	req "github.com/james-wukong/school-schedule/internal/domain/requirement"
	infraPostgre "github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Helpers ─────────────────────────────────────────────────────────────────

func newRequirementRepo(t *testing.T) (req.Repository, sqlmock.Sqlmock) {
	t.Helper()
	gormDB, mock := setupMockDB(t)
	repo := infraPostgre.NewRequirementRepository(gormDB, newLogger())
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

// ── GetByVersion ────────────────────────────────────────────────────────────────

func TestRequirementGetByVersion_Found(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()
	s := sampleRequirement()
	// sem := sampleSemester()
	sch := sampleSchool()
	sub := sampleSubject()
	tch := sampleTeacher()
	cls := sampleClass()
	ts := sampleTimeslot()

	// EXPECTATION 1: The primary Requirement query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements" WHERE`)).
		WithArgs(s.SchoolID, s.SemesterID, s.Version).
		WillReturnRows(mockRequirementRow(mock, s))

	// EXPECTATION 2: The automatic Preload query for the Class
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "classes" WHERE `)).
		WithArgs(cls.ID).
		WillReturnRows(mockClassRow(mock, cls))

	// EXPECTATION 3: The automatic Preload query for the School
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schools" WHERE`)).
		WithArgs(sch.ID).
		WillReturnRows(mockSchoolRow(mock, sch))

	// EXPECTATION 3: The automatic Preload query for the Semester
	// mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "semesters" WHERE`)).
	// 	WithArgs(sem.ID).
	// 	WillReturnRows(mockSemesterRow(mock, sem))

	// EXPECTATION 4: The automatic Preload query for the Subjects
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "subjects" WHERE `)).
		WithArgs(sub.ID).
		WillReturnRows(mockSubjectRow(mock, sub))

	// EXPECTATION 5: The automatic Preload query for the Teacher
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teachers" WHERE `)).
		WithArgs(tch.ID).
		WillReturnRows(mockTeacherRow(mock, tch))

	// EXPECTATION 6: The automatic Preload query for the Timeslots
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "teacher_timeslots" WHERE`)).
		WithArgs(tch.ID).
		WillReturnRows(mockTimeslotRow(mock, ts))

	result, err := repo.GetByVersion(ctx, s.SchoolID, s.SemesterID, s.Version)

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
	// assert.Equal(t, ts, result[0].Teacher.Timeslots[0])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRequirementGetByVersion_NotFound(t *testing.T) {
	repo, mock := newRequirementRepo(t)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "requirements"`)).
		WithArgs(int64(999), float64(11.11)).
		WillReturnRows(sqlmock.NewRows(requirementColumns()))

	result, err := repo.GetByVersion(ctx, int64(999), int64(999), float64(11.11))

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

	result, err := repo.GetByVersion(ctx, int64(999), int64(999), float64(11.11))

	assert.Nil(t, result)
	assert.ErrorIs(t, err, dbErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}
