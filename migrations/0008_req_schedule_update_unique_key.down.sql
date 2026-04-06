BEGIN;

-- Drop requirements unique constraint
ALTER TABLE requirements DROP CONSTRAINT uq_requirements_requirement_key;

ALTER TABLE schedules DROP CONSTRAINT uq_schedules_room_key;
ALTER TABLE schedules DROP CONSTRAINT uq_schedules_requirement_key;


-- Restore schedule constraints
ALTER TABLE schedules ADD CONSTRAINT uq_schedules_room_key
	UNIQUE (school_id, semester_id, room_id, timeslot_id, version);

ALTER TABLE schedules ADD CONSTRAINT uq_schedules_requirement_key
	UNIQUE (school_id, semester_id, requirement_id, timeslot_id, version);

COMMIT;