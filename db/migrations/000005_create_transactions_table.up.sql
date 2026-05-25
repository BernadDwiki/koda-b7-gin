CREATE TABLE transactions (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    amount BIGINT NOT NULL CHECK(amount > 0),

    transaction_type transaction_type NOT NULL,

    note VARCHAR(255),

    status transaction_status
        DEFAULT 'pending',

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP
);