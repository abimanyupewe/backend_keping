-- Enable pgcrypto / uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Table: users (Public schema user extension)
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for User Email
CREATE INDEX IF NOT EXISTS idx_users_email ON public.users(email);

-- 2. Table: transactions
CREATE TABLE IF NOT EXISTS public.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    source VARCHAR(20) NOT NULL CHECK (source IN ('NOTIFICATION', 'OCR_SCAN', 'MANUAL')),
    type VARCHAR(10) NOT NULL CHECK (type IN ('CREDIT', 'DEBIT')),
    amount NUMERIC(15, 2) NOT NULL,
    description VARCHAR(255),
    raw_text TEXT,
    image_url TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for Transactions
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON public.transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON public.transactions(timestamp DESC);

-- 3. Table: wishlists
CREATE TABLE IF NOT EXISTS public.wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    target_amount NUMERIC(15, 2) NOT NULL,
    saved_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    auto_save_percentage NUMERIC(5, 2) NOT NULL DEFAULT 0,
    icon_name VARCHAR(50) DEFAULT 'star',
    label_color VARCHAR(10) DEFAULT '#2A9D8F',
    is_completed BOOLEAN NOT NULL DEFAULT FALSE
);

-- Index for Wishlists
CREATE INDEX IF NOT EXISTS idx_wishlists_user_id ON public.wishlists(user_id);

-- Optional: Supabase Auth Sync Trigger (Sync auth.users -> public.users automatically)
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO public.users (id, email, name, created_at)
  VALUES (new.id, new.email, new.raw_user_meta_data->>'name', new.created_at)
  ON CONFLICT (id) DO UPDATE
  SET email = EXCLUDED.email,
      name = EXCLUDED.name;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Trigger execution on Auth Signup
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;
CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
