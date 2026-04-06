BEGIN;

ALTER TABLE schedules 
	ADD COLUMN semester_id BIGINT NOT NULL REFERENCES semesters(id) ON DELETE CASCADE;

UPDATE schedules SET semester_id = 1000;

ALTER TABLE schedules 
	ALTER COLUMN semester_id SET NOT NULL;


ALTER TABLE requirements 
	ADD COLUMN semester_id BIGINT REFERENCES semesters(id) ON DELETE CASCADE;

UPDATE requirements SET semester_id = 1000;

ALTER TABLE requirements 
	ALTER COLUMN semester_id SET NOT NULL;

COMMIT;