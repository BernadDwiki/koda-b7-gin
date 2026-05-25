CREATE TABLE ewallets (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    user_id INT NOT NULL UNIQUE,

    balance BIGINT DEFAULT 0,
    income BIGINT DEFAULT 0,
    expense BIGINT DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);