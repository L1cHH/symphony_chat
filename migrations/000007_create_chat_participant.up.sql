CREATE TABLE chat_participant (
    chat_id UUID REFERENCES chat(id),
    user_id UUID REFERENCES chat_user(id),
    role_id UUID REFERENCES chat_role(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);