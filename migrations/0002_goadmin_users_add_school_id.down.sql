BEGIN;

ALTER TABLE goadmin_users 
	DROP COLUMN school_id;

COMMIT;