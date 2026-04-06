BEGIN;

ALTER TABLE rooms 
DROP COLUMN available_days;

COMMIT;