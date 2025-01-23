CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE jwt_token (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth_user_id UUID NOT NULL REFERENCES auth_user(id),
    token VARCHAR(255) NOT NULL,
)