-- Create profiles table for storing user data
CREATE TABLE IF NOT EXISTS profiles (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    telegram_username VARCHAR(255),
    telegram_data JSONB,
    nickname VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index on telegram_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_profiles_telegram_id ON profiles(telegram_id);

-- Create index on nickname for faster lookups
CREATE INDEX IF NOT EXISTS idx_profiles_nickname ON profiles(nickname) WHERE nickname IS NOT NULL;

-- Enable Row Level Security
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;

-- Create policy to allow all operations (adjust based on your security needs)
-- For public read/write access with anon key:
CREATE POLICY "Allow all operations on profiles" ON profiles
    FOR ALL
    USING (true)
    WITH CHECK (true);

-- Or if using service role key, you can disable RLS:
-- ALTER TABLE profiles DISABLE ROW LEVEL SECURITY;

