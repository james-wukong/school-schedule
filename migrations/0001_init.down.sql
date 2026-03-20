-- ============================================================================
-- REVERSE/DOWN MIGRATION SCRIPT
-- ============================================================================
-- 1. Start fresh by rolling back any failed state
ROLLBACK; 

-- 2. Drop Views First (They depend on Tables)
DROP VIEW IF EXISTS view_schedule_audit_conflicts CASCADE;
DROP VIEW IF EXISTS v_room_conflicts CASCADE;
DROP VIEW IF EXISTS view_room_utilization CASCADE;
DROP VIEW IF EXISTS view_teacher_workload_stats CASCADE;
DROP VIEW IF EXISTS view_class_timetables CASCADE;
DROP VIEW IF EXISTS v_class_schedule CASCADE;
DROP VIEW IF EXISTS view_teacher_schedules CASCADE;

-- 3. Drop Tables (Child tables first to satisfy Foreign Keys)
-- Using CASCADE here is the "Industry Secret" to solving the Enum dependency
DROP TABLE IF EXISTS schedules CASCADE;
DROP TABLE IF EXISTS requirements CASCADE;
DROP TABLE IF EXISTS room_timeslots CASCADE;
DROP TABLE IF EXISTS teacher_timeslots CASCADE;
DROP TABLE IF EXISTS teacher_subjects CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;
DROP TABLE IF EXISTS timeslots CASCADE;
DROP TABLE IF EXISTS semesters CASCADE;
DROP TABLE IF EXISTS schools CASCADE;

-- 4. Drop the Custom Type
-- Now that 'schedules' is gone, this will succeed.
DROP TYPE IF EXISTS schedule_status_enum;

-- 5. Finalize
COMMIT;
-- ============================================================================
-- END OF REVERSAL
-- ============================================================================
