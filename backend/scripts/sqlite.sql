PRAGMA foreign_keys = on;

-- user
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    nickname TEXT NOT NULL DEFAULT '',
    department TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

-- prompt
CREATE TABLE prompt (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    latest_version TEXT NOT NULL,
    is_publish INTEGER NOT NULL DEFAULT 0,
    created_by TEXT NOT NULL,
    username TEXT NOT NULL,
    category TEXT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_prompt_path ON prompt(path);
CREATE INDEX idx_prompt_publish ON prompt(is_publish);


-- prompt version
CREATE TABLE prompt_version (
    id TEXT PRIMARY KEY,
    prompt_id TEXT NOT NULL,
    version TEXT NOT NULL,
    content TEXT NOT NULL,
    variables TEXT,
    is_publish INTEGER NOT NULL DEFAULT 0,
    change_log TEXT,
    created_by TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 同一 prompt 下版本唯一
CREATE UNIQUE INDEX uk_prompt_version
ON prompt_version(prompt_id, version);

-- 查询发布版本
CREATE INDEX idx_prompt_version_publish
ON prompt_version(prompt_id, is_publish);


CREATE TABLE categories (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    icon TEXT NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    url TEXT NOT NULL,
    created_by TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO categories (id, title, icon, count, url, created_by, username) VALUES
('1', '文案生成', 'file', 5, '/category/copywriting', 'system', 'system'),
('2', '代码助手', 'sparkles', 3, '/category/coding', 'system', 'system'),
('3', '翻译工具', 'file', 2, '/category/translation', 'system', 'system'),
('4', '数据分析', 'file', 4, '/category/analysis', 'system', 'system');


-- favorites (收藏夹)
CREATE TABLE IF NOT EXISTS favorites (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    prompt_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX uk_favorites_user_prompt ON favorites(user_id, prompt_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_prompt_id ON favorites(prompt_id);


-- recently_used (最近使用)
CREATE TABLE IF NOT EXISTS recently_used (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    prompt_id TEXT NOT NULL,
    used_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX uk_recently_used_user_prompt ON recently_used(user_id, prompt_id);
CREATE INDEX IF NOT EXISTS idx_recently_used_user_id ON recently_used(user_id);
CREATE INDEX IF NOT EXISTS idx_recently_used_used_at ON recently_used(user_id, used_at);
