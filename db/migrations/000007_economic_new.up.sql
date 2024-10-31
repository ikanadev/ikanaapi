create table economic_new (
	id UUID,
	title TEXT not null,
	date TIMESTAMPTZ,
	url TEXT not null,
	image TEXT,
	summary TEXT,
	company TEXT not null,
	tags VARCHAR(100)[] NOT NULL DEFAULT '{}',
	sentiment SMALLINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	PRIMARY KEY (id)
);
