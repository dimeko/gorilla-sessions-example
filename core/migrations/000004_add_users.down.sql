START TRANSACTION;

DELETE FROM users
WHERE username == 'takis';

DELETE FROM users
WHERE username == 'sakis';

COMMIT;