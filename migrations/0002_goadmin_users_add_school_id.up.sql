BEGIN;

ALTER TABLE goadmin_users 
	ADD COLUMN school_id BIGINT REFERENCES schools(id) ON DELETE CASCADE;

UPDATE goadmin_users SET school_id = 1000;

ALTER TABLE goadmin_users 
	ALTER COLUMN school_id SET NOT NULL;

COMMIT;