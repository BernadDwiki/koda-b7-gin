CREATE TABLE transfer_details (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    transaction_id INT NOT NULL UNIQUE,

    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,

    FOREIGN KEY (transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    FOREIGN KEY (sender_id)
        REFERENCES users(id),

    FOREIGN KEY (receiver_id)
        REFERENCES users(id),

    CHECK (sender_id <> receiver_id)
);
