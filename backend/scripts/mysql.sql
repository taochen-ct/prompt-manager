-- user
CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username   VARCHAR(64)     NOT NULL,
    nickname   VARCHAR(64)     NOT NULL DEFAULT '',
    department VARCHAR(64)     NOT NULL DEFAULT '',
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_deleted TINYINT(1)      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_username (username),
    KEY idx_users_is_deleted (is_deleted)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- prompt
CREATE TABLE prompt
(
    id             CHAR(36) PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    path           VARCHAR(512) NOT NULL,
    latest_version VARCHAR(64)  NOT NULL,
    is_publish     TINYINT(1)   NOT NULL DEFAULT 0,
    created_by     VARCHAR(64)  NOT NULL,
    username       VARCHAR(64)  NOT NULL,
    category       VARCHAR(64)  NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_prompt_path (path),
    KEY idx_prompt_publish (is_publish)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- prompt version
CREATE TABLE prompt_version
(
    id         CHAR(36) PRIMARY KEY,
    prompt_id  CHAR(36)    NOT NULL,
    version    VARCHAR(64) NOT NULL,
    content    LONGTEXT    NOT NULL,
    variables  LONGTEXT,
    is_publish TINYINT(1)  NOT NULL DEFAULT 0,
    change_log TEXT,
    created_by VARCHAR(64) NOT NULL,
    username   VARCHAR(64) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_prompt_version (prompt_id, version),
    KEY idx_prompt_version_publish (prompt_id, is_publish)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- category
CREATE TABLE prompt_categories
(
    id         VARCHAR(32)  NOT NULL COMMENT '分类ID',
    title      VARCHAR(100) NOT NULL COMMENT '分类标题',
    icon       VARCHAR(50)  NOT NULL COMMENT '图标名称',
    count      INT          NOT NULL DEFAULT 0 COMMENT '数量',
    url        VARCHAR(255) NOT NULL COMMENT '访问路径',
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='提示词分类表';

INSERT INTO prompt_categories (id, title, icon, count, url) VALUES
('1', '文案生成', 'file', 5, '/category/copywriting'),
('2', '代码助手', 'sparkles', 3, '/category/coding'),
('3', '翻译工具', 'file', 2, '/category/translation'),
('4', '数据分析', 'file', 4, '/category/analysis');


-- favorites (收藏夹)
CREATE TABLE IF NOT EXISTS favorites
(
    id         CHAR(36)    NOT NULL PRIMARY KEY,
    user_id    BIGINT      NOT NULL COMMENT '用户ID',
    prompt_id  CHAR(36)    NOT NULL COMMENT '提示词ID',
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_favorites_user_id (user_id),
    INDEX idx_favorites_prompt_id (prompt_id),
    UNIQUE KEY uk_favorites_user_prompt (user_id, prompt_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='收藏夹表';


-- recently_used (最近使用)
CREATE TABLE IF NOT EXISTS recently_used
(
    id         CHAR(36)    NOT NULL PRIMARY KEY,
    user_id    BIGINT      NOT NULL COMMENT '用户ID',
    prompt_id  CHAR(36)    NOT NULL COMMENT '提示词ID',
    used_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '使用时间',
    INDEX idx_recently_used_user_id (user_id),
    INDEX idx_recently_used_used_at (user_id, used_at),
    UNIQUE KEY uk_recently_used_user_prompt (user_id, prompt_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='最近使用表';
