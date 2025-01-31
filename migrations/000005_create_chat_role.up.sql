CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chat_role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE INDEX idx_chat_role_name ON chat_role(name);