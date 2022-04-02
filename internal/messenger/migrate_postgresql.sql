CREATE TABLE IF NOT EXISTS chat(
    id CHAR(24) PRIMARY KEY NOT NULL,
    chat_name VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS messages(
    id CHAR(24) PRIMARY KEY NOT NULL,
    chat_id CHAR(24) REFERENCES chat(id) on delete cascade on update cascade,
    username VARCHAR(100),
    text VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL
)