START TRANSACTION;

INSERT INTO users (username, password, firstname, lastname)
VALUES ('admin', 'adminpass', 'Lakis', 'Mitsis');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('takis', 'password', 'Takis', 'Argyriou');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('sakis', 'password', 'Sakis', 'Petrovelegios');

COMMIT;