CREATE TABLE chat (
    id UUID PRIMARY KEY,
    user_one_id UUID NOT NULL REFERENCES chat_user(id),
    user_two_id UUID NOT NULL REFERENCES chat_user(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
)