CREATE TABLE games (
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(255)                NOT NULL,
    owner_id     VARCHAR(36) REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL   DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE               DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN games.id IS 'Game uniq id';
COMMENT ON COLUMN games.name IS 'Game name';
COMMENT ON COLUMN games.owner_id IS 'User uniq id';
COMMENT ON COLUMN games.created_at IS 'Game created date';
COMMENT ON COLUMN games.updated_at IS 'Game modified date';