BEGIN;

ALTER TABLE requirements 
	DROP COLUMN semester_id;

ALTER TABLE schedules 
	DROP COLUMN semester_id;

COMMIT;
