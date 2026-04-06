BEGIN;

-- Drop new constraints
ALTER TABLE schedules DROP CONSTRAINT uq_schedules_room_key;
ALTER TABLE schedules DROP CONSTRAINT uq_schedules_requirement_key;

-- Drop indexes
DROP INDEX idx_schedules_semester;
DROP INDEX idx_requirements_semester;

-- Recreate original constraints
ALTER TABLE schedules ADD CONSTRAINT schedules_school_id_version_room_id_timeslot_id_key
    UNIQUE (school_id, version, room_id, timeslot_id);

ALTER TABLE schedules ADD CONSTRAINT schedules_school_id_version_requirement_id_timeslot_id_key
    UNIQUE (school_id, version, requirement_id, timeslot_id);

COMMIT;