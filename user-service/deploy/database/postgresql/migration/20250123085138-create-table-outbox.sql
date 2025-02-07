-- +migrate Up
-- Create a table for CDC
-- https://debezium.io/documentation/reference/3.0/transformations/outbox-event-router.html#basic-outbox-table
CREATE TABLE IF NOT EXISTS user_outbox (
    id UUID NOT NULL,
    aggregatetype VARCHAR(255) NOT NULL,
    aggregateid VARCHAR(255) NOT NULL,
    payload BYTEA NOT NULL,
    tracingspancontext JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_outbox_pkey PRIMARY KEY (id)
);

COMMENT ON COLUMN user_outbox.id IS 'The tracing ID or request ID associated with the event';

COMMENT ON COLUMN user_outbox.aggregatetype IS 'The aggregate event type, e.g., registered';

COMMENT ON COLUMN user_outbox.payload IS 'The outbox payload containing the event data in Protobuf binary format';

COMMENT ON COLUMN user_outbox.tracingspancontext IS 'The tracing span context associated with the event';

-- Create a publication for CDC (not needed, since Debezium will create it automatically)
-- CREATE PUBLICATION user_outbox_publication FOR TABLE user_outbox;
-- +migrate Down
-- Drop the publication
-- DROP PUBLICATION user_outbox_publication;
-- Drop the table
DROP TABLE IF EXISTS user_outbox;