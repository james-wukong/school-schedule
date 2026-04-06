BEGIN;

ALTER TABLE timeslots 
	ADD COLUMN school_id BIGINT REFERENCES schools(id) ON DELETE CASCADE;

UPDATE timeslots SET school_id = 1000;

ALTER TABLE timeslots 
	ALTER COLUMN school_id SET NOT NULL;

COMMIT;