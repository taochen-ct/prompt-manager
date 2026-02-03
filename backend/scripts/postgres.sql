-- version
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    nickname VARCHAR(64) NOT NULL DEFAULT '',
    department VARCHAR(64) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

-- prompt
CREATE TABLE prompt (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    latest_version TEXT NOT NULL,
    is_publish BOOLEAN NOT NULL DEFAULT FALSE,
    created_by TEXT NOT NULL,
    username TEXT NOT NULL,
    category TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uk_prompt_path ON prompt(path);
CREATE INDEX idx_prompt_publish ON prompt(is_publish);

-- prompt version
CREATE TABLE prompt_version (
    id UUID PRIMARY KEY,
    prompt_id UUID NOT NULL,
    version TEXT NOT NULL,
    content TEXT NOT NULL,
    variables TEXT,
    is_publish BOOLEAN NOT NULL DEFAULT FALSE,
    change_log TEXT,
    created_by TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 同 prompt 版本唯一
CREATE UNIQUE INDEX uk_prompt_version
ON prompt_version(prompt_id, version);

-- 查询发布版本
CREATE INDEX idx_prompt_version_publish
ON prompt_version(prompt_id, is_publish);


CREATE TABLE prompt_categories (
    id VARCHAR(32) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    icon VARCHAR(50) NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO prompt_categories (id, title, icon, count, url) VALUES
('1', '文案生成', 'file', 5, '/category/copywriting'),
('2', '代码助手', 'sparkles', 3, '/category/coding'),
('3', '翻译工具', 'file', 2, '/category/translation'),
('4', '数据分析', 'file', 4, '/category/analysis');


-- favorites (收藏夹)
CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    prompt_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE INDEX uk_favorites_user_prompt ON favorites(user_id, prompt_id)
);
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_prompt_id ON favorites(prompt_id);


-- recently_used (最近使用)
CREATE TABLE IF NOT EXISTS recently_used (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    prompt_id UUID NOT NULL,
    used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE INDEX uk_recently_used_user_prompt ON recently_used(user_id, prompt_id)
);
CREATE INDEX IF NOT EXISTS idx_recently_used_user_id ON recently_used(user_id);
CREATE INDEX IF NOT EXISTS idx_recently_used_used_at ON recently_used(user_id, used_at);


