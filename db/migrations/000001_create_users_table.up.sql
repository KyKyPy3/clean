CREATE TABLE users (
    id           VARCHAR(36) PRIMARY KEY,
    name         TEXT                      NOT NULL  CHECK ( name <> '' ),
    surname      TEXT,
    middlename   TEXT,
    email        VARCHAR(255)              NOT NULL  UNIQUE,
    password     TEXT                      NOT NULL  CHECK ( password <> '' ),
    created_at   TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE            DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN users.id IS 'User uniq id';
COMMENT ON COLUMN users.name IS 'User first name';
COMMENT ON COLUMN users.surname IS 'User last name';
COMMENT ON COLUMN users.middlename IS 'User middle name';
COMMENT ON COLUMN users.email IS 'User email';
COMMENT ON COLUMN users.password IS 'User hashed password';
COMMENT ON COLUMN users.created_at IS 'User created date';
COMMENT ON COLUMN users.updated_at IS 'User modified date';