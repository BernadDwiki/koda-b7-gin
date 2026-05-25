CREATE TYPE transaction_type AS ENUM (
    'transfer',
    'top_up'
);

CREATE TYPE transaction_status AS ENUM (
    'pending',
    'success',
    'failed'
);
