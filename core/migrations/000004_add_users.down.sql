START TRANSACTION;

DELETE FROM users
WHERE username == 'takis';

DELETE FROM users
WHERE username == 'sakis';

DELETE FROM users
WHERE username == 'admin';

COMMIT;