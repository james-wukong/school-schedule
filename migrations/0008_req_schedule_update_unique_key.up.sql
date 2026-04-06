BEGIN;

-- Drop unique constraint in schedule table
ALTER TABLE schedules DROP CONSTRAINT uq_schedules_room_key;
ALTER TABLE schedules DROP CONSTRAINT uq_schedules_requirement_key;

-- Update constraints
ALTER TABLE schedules ADD CONSTRAINT uq_schedules_room_key
	UNIQUE (semester_id, room_id, timeslot_id, version);
ALTER TABLE schedules ADD CONSTRAINT uq_schedules_requirement_key
	UNIQUE (semester_id, requirement_id, timeslot_id, version);

-- Add constraint for requirements table
ALTER TABLE requirements ADD CONSTRAINT uq_requirements_requirement_key
	UNIQUE (semester_id, subject_id, teacher_id, class_id, version);

COMMIT;