ALTER TABLE ustd_price
	ALTER COLUMN created_at TYPE TIMESTAMPTZ,
	ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
	ALTER COLUMN deleted_at TYPE TIMESTAMPTZ,
	ALTER COLUMN archived_at TYPE TIMESTAMPTZ;
