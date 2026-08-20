-- Seed Data for Keping Backend Database

-- 1. Seed Users (Public extension schema)
INSERT INTO public.users (id, email, name, created_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'budi@keping.app', 'Budi Santoso', NOW() - INTERVAL '30 days'),
    ('22222222-2222-2222-2222-222222222222', 'clara@keping.app', 'Clara Financial', NOW() - INTERVAL '15 days')
ON CONFLICT (id) DO UPDATE 
SET email = EXCLUDED.email, 
    name = EXCLUDED.name;

-- 2. Seed Wishlists
INSERT INTO public.wishlists (id, user_id, name, target_amount, saved_amount, auto_save_percentage, icon_name, label_color, is_completed)
VALUES 
    ('a1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'Laptop Gaming Baru', 15000000, 3500000, 10, 'laptop', '#4A90E2', FALSE),
    ('a2222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Tiket Konser', 2500000, 2500000, 5, 'ticket', '#8E44AD', TRUE),
    ('a3333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222', 'Dana Darurat 6 Bulan', 30000000, 12000000, 15, 'shield', '#2A9D8F', FALSE)
ON CONFLICT (id) DO UPDATE 
SET target_amount = EXCLUDED.target_amount,
    saved_amount = EXCLUDED.saved_amount, 
    auto_save_percentage = EXCLUDED.auto_save_percentage,
    is_completed = EXCLUDED.is_completed;

-- 3. Seed Transactions
INSERT INTO public.transactions (id, user_id, source, type, amount, description, raw_text, image_url, timestamp)
VALUES 
    (
        'c1111111-1111-1111-1111-111111111111', 
        '11111111-1111-1111-1111-111111111111', 
        'NOTIFICATION', 
        'CREDIT', 
        500000, 
        'Isi Saldo DANA via BRI', 
        'Kamu telah berhasil melakukan isi saldo sebesar Rp500.000 via BRI', 
        NULL, 
        NOW() - INTERVAL '5 days'
    ),
    (
        'c2222222-2222-2222-2222-222222222222', 
        '11111111-1111-1111-1111-111111111111', 
        'OCR_SCAN', 
        'DEBIT', 
        45000, 
        'Kopi Kenangan Mantan', 
        'Merchant: Kopi Kenangan, Total: Rp 45.000', 
        'https://res.cloudinary.com/dmzt2vecd/image/upload/v1/receipts/sample_receipt.jpg', 
        NOW() - INTERVAL '3 days'
    ),
    (
        'c3333333-3333-3333-3333-333333333333', 
        '11111111-1111-1111-1111-111111111111', 
        'NOTIFICATION', 
        'CREDIT', 
        1000000, 
        'Transfer Masuk ShopeePay', 
        'Transfer sebesar Rp1.000.000 dari Sdr. Budi telah masuk ke ShopeePay', 
        NULL, 
        NOW() - INTERVAL '1 day'
    ),
    (
        'c4444444-4444-4444-4444-444444444444', 
        '22222222-2222-2222-2222-222222222222', 
        'MANUAL', 
        'DEBIT', 
        750000, 
        'Belanja Mingguan Supermarket', 
        'Pembelian sembako dan kebutuhan rumah tangga', 
        NULL, 
        NOW() - INTERVAL '2 days'
    )
ON CONFLICT (id) DO UPDATE
SET amount = EXCLUDED.amount;
