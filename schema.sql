-- Kasir API Database Schema

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(200) NOT NULL,
    harga INTEGER NOT NULL DEFAULT 0,
    stok INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample categories
INSERT INTO categories (name, description) VALUES
    ('Minuman', 'Minuman dingin dan panas'),
    ('Makanan', 'Makanan ringan dan berat'),
    ('Snack', 'Camilan dan kudapan')
ON CONFLICT DO NOTHING;

-- Insert sample products
INSERT INTO products (nama, harga, stok) VALUES
    ('Es Teh Manis', 5000, 100),
    ('Nasi Goreng', 25000, 50),
    ('Keripik Kentang', 8000, 200)
ON CONFLICT DO NOTHING;
