-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE public.authors (
	id SERIAL UNIQUE PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT,
	created_at TIMESTAMP WITH TIME ZONE,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE authors;