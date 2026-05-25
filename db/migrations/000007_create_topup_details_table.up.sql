CREATE TABLE top_up_details (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    transaction_id INT NOT NULL UNIQUE,

    receiver_id INT NOT NULL,

    payment_method_id INT NOT NULL,

    tax BIGINT DEFAULT 0,
    admin_fee BIGINT DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,

    FOREIGN KEY (transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    FOREIGN KEY (receiver_id)
        REFERENCES users(id),

    FOREIGN KEY (payment_method_id)
        REFERENCES payment_methods(id),

    CHECK (tax >= 0),
    CHECK (admin_fee >= 0)
);