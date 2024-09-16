CREATE TABLE IF NOT EXISTS users (
    guid VARCHAR(36) PRIMARY KEY,
    email VARCHAR(100),
    client_ip VARCHAR(45),
    refresh_token TEXT
);
