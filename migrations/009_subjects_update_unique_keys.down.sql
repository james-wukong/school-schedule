BEGIN;

--Add new unique constraint
ALTER TABLE subjects
	DROP CONSTRAINT subjects_school_id_name_key;

-- DROP old unique constraints
ALTER TABLE subjects
	ADD CONSTRAINT subjects_name_key UNIQUE(name), 
	ADD CONSTRAINT subjects_code_key UNIQUE(code);

COMMIT;