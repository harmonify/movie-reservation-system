ALTER SYSTEM
SET
    wal_level = 'logical';

ALTER SYSTEM
SET
    max_replication_slots = 10;

ALTER SYSTEM
SET
    max_wal_senders = 10;

CREATE DATABASE "mvs-user-service" IF NOT EXISTS;

\c "mvs-user-service";

CREATE EXTENSION "pgcrypto" IF NOT EXISTS;
