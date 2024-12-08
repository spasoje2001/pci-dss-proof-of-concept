CREATE TABLE cardholders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    card_number VARCHAR(16) NOT NULL,
    expiration_date DATE NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
