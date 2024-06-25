-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE roles DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;
ALTER TABLE authors DROP COLUMN deleted_at;
ALTER TABLE posts DROP COLUMN deleted_at;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE roles ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE authors ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE posts ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

