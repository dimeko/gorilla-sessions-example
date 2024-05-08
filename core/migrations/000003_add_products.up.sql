START TRANSACTION;

-- Create a new products in the 'products' table
INSERT INTO products (title, description, price)
VALUES ('product_1', 'a nice product 1', 30);

INSERT INTO products (title, description, price)
VALUES ('product_2', 'a nice product 2', 40);

COMMIT;