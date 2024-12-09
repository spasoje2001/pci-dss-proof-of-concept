CREATE TABLE users (
    id SERIAL PRIMARY KEY, -- Automatski generisani ID
    username VARCHAR(255) NOT NULL UNIQUE, -- Korisničko ime, jedinstveno
    password VARCHAR(255) NOT NULL, -- Šifrovana lozinka
    role VARCHAR(50) NOT NULL -- Uloga korisnika (npr. "admin", "user")
);