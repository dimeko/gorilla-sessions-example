START TRANSACTION;

DELETE FROM products
WHERE title == 'product_1';

DELETE FROM products
WHERE title == 'product_2';

COMMIT;