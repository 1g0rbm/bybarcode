CREATE TABLE IF NOT EXISTS shopping_list_statistics
(
    id                     SERIAL PRIMARY KEY,
    shopping_list_id       INT REFERENCES shopping_lists (id),
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'::text),
    added_products_count   INT NOT NULL             DEFAULT 0,
    checked_products_count INT NOT NULL             DEFAULT 0
);