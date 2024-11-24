CREATE TABLE game_user (
    game_id VARCHAR(36) REFERENCES games (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id VARCHAR(36) REFERENCES users (id) ON UPDATE CASCADE,
    role numeric NOT NULL,
    CONSTRAINT game_participant_pkey PRIMARY KEY (game_id, user_id)
);

COMMENT ON COLUMN game_user.game_id IS 'Game uniq id';
COMMENT ON COLUMN game_user.user_id IS 'Role uniq id';
COMMENT ON COLUMN game_user.role IS 'User role';