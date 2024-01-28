CREATE TABLE outbox (
    id           BIGSERIAL PRIMARY KEY,
    topic        VARCHAR(50)               NOT NULL,
    kind         VARCHAR(50)               NOT NULL,
    payload      BYTEA,
    consumed     BOOLEAN                   NOT NULL  DEFAULT FALSE
);

COMMENT ON COLUMN outbox.id IS 'Outbox message id';
COMMENT ON COLUMN outbox.topic IS 'Broker topic name';
COMMENT ON COLUMN outbox.topic IS 'Event kind';
COMMENT ON COLUMN outbox.payload IS 'Message payload';
COMMENT ON COLUMN outbox.consumed IS 'Is message sent';