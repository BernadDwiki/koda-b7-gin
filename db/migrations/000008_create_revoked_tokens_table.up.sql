CREATE TABLE revoked_tokens (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    user_id INT NOT NULL,

    token TEXT NOT NULL,

    expired_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
