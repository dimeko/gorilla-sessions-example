START TRANSACTION;

INSERT INTO users (username, password, firstname, lastname)
VALUES ('admin', '$2a$04$ZHX0QBV/V/T9slKTUezQsu8WfTfYlrQnnEmpYzHd9VHAaibS/Uumq', 'Lakis', 'Mitsis');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('takis', '$2a$04$agVyVYhMgYs.bX.4MMmr4.gg5.ki8low/YN.UmCSQ4Lk.Ki3pIpUi', 'Takis', 'Argyriou');

INSERT INTO users (username, password, firstname, lastname)
VALUES ('sakis', '$2a$04$csLRza1bwJuUyv2RHsvgI.aC8F9jFO/TJIye1ehzQY4dSdSr9fRym', 'Sakis', 'Petrovelegios');

COMMIT;