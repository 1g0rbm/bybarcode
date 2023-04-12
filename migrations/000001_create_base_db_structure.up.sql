CREATE TABLE IF NOT EXISTS account
(
    id         BIGINT PRIMARY KEY,
    first_name VARCHAR(150) NOT NULL,
    last_name  VARCHAR(150) NOT NULL,
    username   VARCHAR(150) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS shopping_lists
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    account_id BIGINT REFERENCES account (id)
);

CREATE TABLE IF NOT EXISTS categories
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS brands
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    upcean      BIGINT       NOT NULL UNIQUE,
    category_id INT REFERENCES categories (id),
    brand_id    INT REFERENCES brands (id)
);

CREATE TABLE IF NOT EXISTS shopping_list__products
(
    shopping_list_id INT REFERENCES shopping_lists (id),
    product_id       INT REFERENCES products (id),
    PRIMARY KEY (shopping_list_id, product_id)
);