BEGIN;

CREATE INDEX idx_schedules_lookup 
ON schedules (semester_id, version, timeslot_id);

CREATE INDEX idx_timeslots_order 
ON timeslots (day_of_week, start_time);

COMMIT;
