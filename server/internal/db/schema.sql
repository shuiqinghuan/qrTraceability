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
    variety_code VARCHAR(20),
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

-- 创建种植标签表
CREATE TABLE IF NOT EXISTS planting_tags (
    id SERIAL PRIMARY KEY,
    planting_id INTEGER REFERENCES planting(id),
    tag_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(planting_id, tag_name)
);

-- 创建产品品质指标表
CREATE TABLE IF NOT EXISTS product_quality (
    id SERIAL PRIMARY KEY,
    planting_id INTEGER REFERENCES planting(id),
    harvest_start_date DATE,
    harvest_end_date DATE,
    sugar_content DECIMAL(5,2),
    weight DECIMAL(10,2),
    taste_description TEXT,
    suitable_for TEXT,
    quality_summary TEXT,
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

-- 创建IP点赞限制表
CREATE TABLE IF NOT EXISTS ip_like_restrictions (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    planting_id INTEGER REFERENCES planting(id),
    like_count INTEGER DEFAULT 0,
    last_like_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ip_address, planting_id)
);

-- 创建后台管理用户表
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认管理员账号（密码：123456，使用bcrypt哈希）
INSERT INTO admin_users (username, password_hash, role) 
VALUES ('lhseed', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin');

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_seed_info_created_by ON seed_info(created_by);
CREATE INDEX IF NOT EXISTS idx_seed_info_variety_code ON seed_info(variety_code);
CREATE INDEX IF NOT EXISTS idx_planting_seed_id ON planting(seed_id);
CREATE INDEX IF NOT EXISTS idx_planting_created_by ON planting(created_by);
CREATE INDEX IF NOT EXISTS idx_planting_location ON planting(location);
CREATE INDEX IF NOT EXISTS idx_growth_media_planting_id ON growth_media(planting_id);
CREATE INDEX IF NOT EXISTS idx_growth_media_created_by ON growth_media(created_by);
CREATE INDEX IF NOT EXISTS idx_planting_tags_planting_id ON planting_tags(planting_id);
CREATE INDEX IF NOT EXISTS idx_product_quality_planting_id ON product_quality(planting_id);
CREATE INDEX IF NOT EXISTS idx_product_quality_created_by ON product_quality(created_by);
CREATE INDEX IF NOT EXISTS idx_product_quality_harvest_dates ON product_quality(harvest_start_date, harvest_end_date);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_planting_id ON favorites(planting_id);
CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id);
CREATE INDEX IF NOT EXISTS idx_likes_planting_id ON likes(planting_id);
CREATE INDEX IF NOT EXISTS idx_ip_like_restrictions_ip ON ip_like_restrictions(ip_address);
CREATE INDEX IF NOT EXISTS idx_ip_like_restrictions_planting ON ip_like_restrictions(planting_id);
CREATE INDEX IF NOT EXISTS idx_ip_like_restrictions_last_time ON ip_like_restrictions(last_like_time);
CREATE INDEX IF NOT EXISTS idx_admin_users_username ON admin_users(username);
CREATE INDEX IF NOT EXISTS idx_admin_users_role ON admin_users(role);