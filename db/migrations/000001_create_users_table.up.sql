CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION citext;

CREATE DOMAIN EMAIL AS CITEXT
    CHECK ( value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );

CREATE TABLE users (
    id           UUID PRIMARY KEY                    DEFAULT uuid_generate_v4(),
    name         TEXT                      NOT NULL  CHECK ( name <> '' ),
    surname      TEXT                      NOT NULL  CHECK ( surname <> '' ),
    middlename   TEXT                      NOT NULL  CHECK ( middlename <> '' ),
    email        EMAIL                     NOT NULL  UNIQUE,
    created_at   TIMESTAMP WITH TIME ZONE  NOT NULL  DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE            DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN users.id IS 'User uniq id';
COMMENT ON COLUMN users.name IS 'User first name';
COMMENT ON COLUMN users.surname IS 'User last name';
COMMENT ON COLUMN users.middlename IS 'User middle name';
COMMENT ON COLUMN users.email IS 'User email';
COMMENT ON COLUMN users.created_at IS 'User created date';
COMMENT ON COLUMN users.updated_at IS 'User modified date';