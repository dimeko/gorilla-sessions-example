START TRANSACTION;

INSERT INTO users (username, password, firstname, lastname)
VALUES ('admin', '713bfda78870bf9d1b261f565286f85e97ee614efe5f0faf7c34e7ca4f65baca', 'Lakis', 'Mitsis');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('takis', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', 'Takis', 'Argyriou');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('sakis', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', 'Sakis', 'Petrovelegios');

COMMIT;