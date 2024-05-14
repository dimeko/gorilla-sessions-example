START TRANSACTION;

-- Create a new products in the 'products' table
INSERT INTO products (name, title, description, price)
VALUES ('spoon', 'Spoon', 'Lorem ipsum spoon', 30);

INSERT INTO products (name, title, description, price)
VALUES ('fork', 'Fork', 'Lorem ipsum fork', 40);

INSERT INTO products (name, title, description, price)
VALUES ('bigspoon', 'Big spoon', 'Lorem ipsum bigspoon', 332);

INSERT INTO products (name, title, description, price)
VALUES ('kettle', 'Kettle', 'Lorem ipsum kettle', 23);

INSERT INTO products (name, title, description, price)
VALUES ('saucepan', 'Saucepan', 'Lorem ipsum saucepan', 43);

INSERT INTO products (name, title, description, price)
VALUES ('pot', 'Pot', 'Lorem ipsum pot', 65);

INSERT INTO products (name, title, description, price)
VALUES ('pan', 'Pan', 'Lorem ipsum pan', 14);

INSERT INTO products (name, title, description, price)
VALUES ('glass', 'Glass', 'Lorem ipsum glass', 43);

INSERT INTO products (name, title, description, price)
VALUES ('bottle', 'Bottle', 'Lorem ipsum bottle', 3);

COMMIT;