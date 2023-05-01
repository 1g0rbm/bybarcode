ALTER TABLE products ALTER COLUMN id TYPE bigint;
ALTER SEQUENCE products_id_seq RESTART WITH 5006351;
ALTER TABLE products ALTER COLUMN id SET DEFAULT nextval('products_id_seq');
