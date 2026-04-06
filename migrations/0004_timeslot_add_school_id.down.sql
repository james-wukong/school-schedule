BEGIN;

ALTER TABLE timeslots 
	DROP COLUMN school_id;

COMMIT;