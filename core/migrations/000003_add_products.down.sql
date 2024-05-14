START TRANSACTION;

DELETE FROM products
WHERE name == 'spoon';

DELETE FROM products
WHERE name == 'fork';

DELETE FROM products
WHERE name == 'bigspoon';

DELETE FROM products
WHERE name == 'kettle';

DELETE FROM products
WHERE name == 'saucepan';

DELETE FROM products
WHERE name == 'pot';

DELETE FROM products
WHERE name == 'pan';

DELETE FROM products
WHERE name == 'glass';

DELETE FROM products
WHERE name == 'bottle';

COMMIT;
