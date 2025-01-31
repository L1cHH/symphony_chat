CREATE TABLE chat_role_permission (
    role_id UUID REFERENCES chat_role(id),
    permission VARCHAR(128) NOT NULL,
    PRIMARY KEY (role_id, permission)
);