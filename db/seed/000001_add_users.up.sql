CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;
LOCK TABLE users;
INSERT INTO users VALUES (gen_random_uuid(),'Ivan','Ivanov', 'Ivanovich', 'ivan@email.com');
INSERT INTO users VALUES (gen_random_uuid(),'Alise','Smith', 'Saint', 'alise@email.com');
COMMIT;