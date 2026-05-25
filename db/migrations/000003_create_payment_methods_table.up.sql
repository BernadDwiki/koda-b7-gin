CREATE TABLE payment_methods (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    method_name VARCHAR(255) NOT NULL UNIQUE,
    tax_percent INT DEFAULT 0,
    admin_fee BIGINT DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP
);