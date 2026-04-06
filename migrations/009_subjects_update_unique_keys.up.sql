BEGIN;

-- DROP old unique constraints
ALTER TABLE subjects
	DROP CONSTRAINT subjects_name_key, 
	DROP CONSTRAINT subjects_code_key;

--Add new unique constraint
ALTER TABLE subjects
	ADD CONSTRAINT subjects_school_id_name_key UNIQUE (school_id, name);

COMMIT;