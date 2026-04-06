BEGIN;

-- Drop old constraints
ALTER TABLE schedules DROP CONSTRAINT schedules_school_id_version_room_id_timeslot_id_key;
ALTER TABLE schedules DROP CONSTRAINT schedules_school_id_version_requirement_id_timeslot_id_key;

-- Update constraints
ALTER TABLE schedules ADD CONSTRAINT uq_schedules_room_key
	UNIQUE (school_id, semester_id, room_id, timeslot_id, version);

ALTER TABLE schedules ADD CONSTRAINT uq_schedules_requirement_key
	UNIQUE (school_id, semester_id, requirement_id, timeslot_id, version);


-- Add new indexes
CREATE INDEX idx_schedules_semester
ON schedules(semester_id);

CREATE INDEX idx_requirements_semester
ON requirements(semester_id);

COMMIT;