-- ============================================================================
-- REVERSE/DOWN MIGRATION SCRIPT
-- ============================================================================
-- Start transaction to ensure atomicity
BEGIN;

-- 1. Drop Analytical and Summary Views
-- These must go first as they depend on the tables below.
DROP VIEW IF EXISTS potential_conflicts;
DROP VIEW IF EXISTS classroom_occupancy;
DROP VIEW IF EXISTS teacher_workload;
DROP VIEW IF EXISTS course_schedule;
DROP VIEW IF EXISTS classroom_usage;
DROP VIEW IF EXISTS teacher_schedule;
DROP VIEW IF EXISTS schedule_summary;

-- 2. Drop Schedule Entries
-- Depends on: schedules, courses, teachers, classrooms, timeslots
DROP TABLE IF EXISTS schedule_entries;

-- 3. Drop Schedules
-- Depends on: schools
DROP TABLE IF EXISTS schedules;

-- 4. Drop Constraints and Rules
-- Depends on: schools
DROP TABLE IF EXISTS constraints;

-- 5. Drop Timeslots
-- Depends on: schools
DROP TABLE IF EXISTS timeslots;

-- 6. Drop Classrooms
-- Depends on: schools
DROP TABLE IF EXISTS classrooms;

-- 7. Drop Courses
-- Depends on: schools, teachers
DROP TABLE IF EXISTS courses;

-- 8. Drop Students
-- Depends on: schools
DROP TABLE IF EXISTS students;

-- 9. Drop Teachers
-- Depends on: schools
DROP TABLE IF EXISTS teachers;

-- 10. Drop Schools (The Root Entity)
DROP TABLE IF EXISTS schools;

-- Note: Postgres implicitly drops indexes associated with tables.
-- If you created standalone indexes on system tables (not the case here), 
-- they would be dropped here.

COMMIT;
-- ============================================================================
-- END OF REVERSAL
-- ============================================================================