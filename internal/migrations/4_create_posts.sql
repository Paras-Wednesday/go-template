-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE public.posts (
	id SERIAL UNIQUE PRIMARY KEY,
	author_id INT REFERENCES authors(id) NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE,
	updated_at TIMESTAMP WITH TIME ZONE,
	deleted_at TIMESTAMP WITH TIME ZONE
);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE posts;