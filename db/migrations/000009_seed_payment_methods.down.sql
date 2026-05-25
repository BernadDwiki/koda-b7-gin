DELETE FROM payment_methods
WHERE method_name IN (
    'BRI',
    'Dana',
    'BCA',
    'Gopay',
    'Ovo'
);
