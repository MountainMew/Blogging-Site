package db

const Schema = `
	CREATE TABLE IF NOT EXISTS users (
    ID SERIAL PRIMARY KEY,
    Username VARCHAR(50) UNIQUE NOT NULL,
    Password TEXT NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    Role VARCHAR(50) NOT NULL
);
`
