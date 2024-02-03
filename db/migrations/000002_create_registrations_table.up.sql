CREATE TABLE registrations (
    id           VARCHAR(36) PRIMARY KEY,
    email        VARCHAR(255)              NOT NULL  UNIQUE,
    password     TEXT                      NOT NULL  CHECK ( password <> '' ),
    verified     BOOLEAN                   NOT NULL  DEFAULT FALSE
);

COMMENT ON COLUMN registrations.id IS 'User uniq id';
COMMENT ON COLUMN registrations.email IS 'User email';
COMMENT ON COLUMN registrations.password IS 'User hashed password';
COMMENT ON COLUMN registrations.verified IS 'Is registration verified';