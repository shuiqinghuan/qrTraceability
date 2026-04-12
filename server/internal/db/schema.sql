-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('serverseed', 'servergrow', 'servermanager', 'clentcustomer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建种子信息表
CREATE TABLE IF NOT EXISTS seed_info (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    variety VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建播种定植时间表
CREATE TABLE IF NOT EXISTS planting (
    id SERIAL PRIMARY KEY,
    seed_id INTEGER REFERENCES seed_info(id),
    planting_date DATE NOT NULL,
    transplanting_date DATE,
    location VARCHAR(200),
    notes TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建生长媒体表（图片和视频）
CREATE TABLE IF NOT EXISTS growth_media (
    id SERIAL PRIMARY KEY,
    planting_id INTEGER REFERENCES planting(id),
    media_type VARCHAR(10) CHECK (media_type IN ('image', 'video')),
    file_path VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建产品品质指标表
CREATE TABLE IF NOT EXISTS product_quality (
    id SERIAL PRIMARY KEY,
    planting_id INTEGER REFERENCES planting(id),
    sugar_content DECIMAL(5,2),
    weight DECIMAL(10,2),
    taste_analysis TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    planting_id INTEGER REFERENCES planting(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, planting_id)
);

-- 创建红心表
CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    planting_id INTEGER REFERENCES planting(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, planting_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_seed_info_created_by ON seed_info(created_by);
CREATE INDEX IF NOT EXISTS idx_planting_seed_id ON planting(seed_id);
CREATE INDEX IF NOT EXISTS idx_planting_created_by ON planting(created_by);
CREATE INDEX IF NOT EXISTS idx_growth_media_planting_id ON growth_media(planting_id);
CREATE INDEX IF NOT EXISTS idx_growth_media_created_by ON growth_media(created_by);
CREATE INDEX IF NOT EXISTS idx_product_quality_planting_id ON product_quality(planting_id);
CREATE INDEX IF NOT EXISTS idx_product_quality_created_by ON product_quality(created_by);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_planting_id ON favorites(planting_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_planting_id ON likes(planting_id);