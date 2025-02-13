ALTER SYSTEM
SET
    wal_level = 'logical';

ALTER SYSTEM
set
    wal_keep_size = "2GB";

ALTER SYSTEM
SET
    max_replication_slots = 10;

ALTER SYSTEM
SET
    max_wal_senders = 10;

SELECT
    pg_reload_conf();

CREATE DATABASE "mrs-user-service" IF NOT EXISTS;

\c "mrs-user-service";

CREATE EXTENSION "pgcrypto" IF NOT EXISTS;
