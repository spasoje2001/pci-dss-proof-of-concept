CREATE TABLE cardholders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    card_number BYTEA NOT NULL,
    expiration_date VARCHAR(16) NOT NULL
);
