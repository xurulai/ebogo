-- 创建用户表
CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, -- 用户ID，主键
    username VARCHAR(50) NOT NULL,              -- 用户名
    password VARCHAR(255) NOT NULL,             -- 密码
    email VARCHAR(100) NOT NULL,                -- 邮箱
    phone VARCHAR(20) UNIQUE,                   -- 手机号，唯一
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    PRIMARY KEY (id),                           -- 设置主键
    UNIQUE INDEX idx_username (username),       -- 用户名唯一索引
    UNIQUE INDEX idx_email (email),             -- 邮箱唯一索引
    UNIQUE INDEX idx_phone (phone)              -- 手机号唯一索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';