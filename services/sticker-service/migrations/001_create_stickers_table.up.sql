CREATE TABLE IF NOT EXISTS stickers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    preview_url VARCHAR(500),
    file_url VARCHAR(500) NOT NULL,
    price DECIMAL(10,2) DEFAULT 0,
    category_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS sticker_packs (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    preview_url VARCHAR(500),
    price DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL
);

