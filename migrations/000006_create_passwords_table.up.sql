CREATE TABLE IF NOT EXISTS passwords (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE
);
