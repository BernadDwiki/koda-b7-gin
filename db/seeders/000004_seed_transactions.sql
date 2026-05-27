INSERT INTO transactions (
    amount,
    transaction_type,
    note,
    status
)
VALUES
(
    100000,
    'transfer',
    'Bayar makan',
    'success'
),
(
    200000,
    'transfer',
    'Bayar hutang',
    'success'
),
(
    300000,
    'top_up',
    'Top up saldo',
    'success'
),
(
    150000,
    'top_up',
    'Isi saldo',
    'pending'
);
