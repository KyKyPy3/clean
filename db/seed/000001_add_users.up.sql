BEGIN;
LOCK TABLE users;
INSERT INTO users VALUES ('2b0c8791-2136-46b6-bc38-b33038ca2e80','Ivan','Ivanov', 'Ivanovich', 'ivan@email.com');
INSERT INTO users VALUES ('c56ace69-ae54-4ecf-beb5-d3f314d3ee03','Alise','Smith', 'Saint', 'alise@email.com');
COMMIT;