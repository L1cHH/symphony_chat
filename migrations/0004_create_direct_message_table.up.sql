CREATE TYPE MESSAGE_STATUS AS ENUM ('sent', 'read', 'undelivered');

CREATE TABLE direct_message(
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chat(id),
    sender_id UUID NOT NULL REFERENCES chat_user(id),
    recipient_id UUID NOT NULL REFERENCES chat_user(id),
    message_content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    status MESSAGE_STATUS NOT NULL DEFAULT 'sent',
    is_edited BOOLEAN NOT NULL DEFAULT false
)
