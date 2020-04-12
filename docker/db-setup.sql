CREATE DATABASE links;

\c links;

CREATE TABLE link (
    id UUID PRIMARY KEY,
    long_url TEXT,
    short_url VARCHAR(8) UNIQUE
);
