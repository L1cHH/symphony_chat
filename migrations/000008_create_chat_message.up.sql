CREATE TABLE chat_message (
    id UUID PRIMARY KEY,
    chat_id UUID REFERENCES chat(id),
    sender_id UUID REFERENCES chat_user(id),
    context TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'sent'
);