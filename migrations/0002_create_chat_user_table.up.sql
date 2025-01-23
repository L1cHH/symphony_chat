CREATE TABLE chat_user (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL DEFAULT "offline",
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMP NOT NULL DEFAULT now()
)

