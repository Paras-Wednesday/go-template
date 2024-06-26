-- +migrate Up
ALTER TABLE authors
ADD COLUMN email TEXT UNIQUE NOT NULL,
ADD COLUMN password TEXT NOT NULL;

-- +migrate Down
ALTER TABLE authors
DROP COLUMN email,
DROP COLUMN password;
