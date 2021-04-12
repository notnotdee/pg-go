CREATE TABLE villagers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    image TEXT NOT NULL,
    species TEXT NOT NULL,
    personality TEXT NOT NULL,
    birthday TEXT NOT NULL,
    quote TEXT NOT NULL
);
