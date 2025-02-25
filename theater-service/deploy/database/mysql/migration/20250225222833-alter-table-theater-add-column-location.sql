-- +migrate Up
ALTER TABLE
    theater
ADD
    COLUMN location POINT NOT NULL
AFTER
    website,
ADD
    SPATIAL INDEX idx_location (location);

-- +migrate Down
ALTER TABLE
    theater DROP COLUMN location,
    DROP INDEX idx_location;