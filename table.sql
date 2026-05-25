CREATE TABLE users (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255),

    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    pin VARCHAR(255),

    picture VARCHAR(255),
    phone_number VARCHAR(20) UNIQUE,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP
);

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

CREATE TABLE payment_methods (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    method_name VARCHAR(255) NOT NULL UNIQUE,
    tax_percent INT DEFAULT 0,
    admin_fee BIGINT DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP
);

-- seed payment methods
INSERT INTO payment_methods (method_name, tax_percent, admin_fee)
VALUES
('BRI', 5, 1500),
('Dana', 3, 1000),
('BCA', 4, 1200),
('Gopay', 2, 800),
('Ovo', 3, 900);

CREATE TYPE transaction_type AS ENUM (
    'transfer',
    'top_up'
);

CREATE TYPE transaction_status AS ENUM (
    'pending',
    'success',
    'failed'
);

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
