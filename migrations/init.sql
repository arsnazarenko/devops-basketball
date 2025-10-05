CREATE TABLE IF NOT EXISTS players (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL CHECK (LENGTH(name) >= 1),
    surname VARCHAR(50) NOT NULL CHECK (LENGTH(surname) >= 1),
    age INTEGER NOT NULL CHECK (age >= 15 AND age <= 50),
    height INTEGER NOT NULL CHECK (height >= 1500),
    weight INTEGER NOT NULL CHECK (weight >= 50000),
    citizenship VARCHAR(10) NOT NULL CHECK (LENGTH(citizenship) >= 2),
    role VARCHAR(2) NOT NULL CHECK (role IN ('PG', 'SG', 'SF', 'PF', 'C')),
    team_id BIGINT NOT NULL CHECK (team_id >= 1)
);
