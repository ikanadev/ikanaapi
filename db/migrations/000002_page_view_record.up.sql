CREATE TABLE page_view_record (
	id UUID,
	app VARCHAR(255) NOT NULL,
	user_id VARCHAR(255) NOT NULL,
	url VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	PRIMARY KEY (id)
);