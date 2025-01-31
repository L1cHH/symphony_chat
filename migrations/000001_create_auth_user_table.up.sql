CREATE TABLE auth_user (
    id UUID PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    registration_at TIMESTAMP NOT NULL DEFAULT now()
)
