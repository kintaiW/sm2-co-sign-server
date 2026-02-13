-- SM2 协同签名服务数据库初始化脚本
-- 数据库: SQLite3

-- 启用 WAL 模式以提高并发性能
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,              -- 用户ID (UUID)
    username TEXT UNIQUE NOT NULL,    -- 用户名
    password_hash TEXT NOT NULL,      -- 密码哈希 (SM3, hex编码)
    public_key TEXT NOT NULL,         -- 协同公钥 Pa (Base64)
    status INTEGER DEFAULT 1,         -- 状态: 1=启用, 0=禁用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 密钥分量表
CREATE TABLE IF NOT EXISTS keys (
    id TEXT PRIMARY KEY,              -- 密钥ID (UUID)
    user_id TEXT NOT NULL,            -- 所属用户ID
    d2 TEXT NOT NULL,                 -- 服务端私钥分量 D2 (Base64, 加密存储)
    d2_inv TEXT NOT NULL,             -- D2 的逆 (Base64, 加密存储)
    public_key TEXT NOT NULL,         -- 协同公钥 Pa (Base64)
    hmac_key TEXT,                    -- HMAC 密钥 (Base64, 用于验证)
    status INTEGER DEFAULT 1,         -- 状态: 1=启用, 0=禁用
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 会话表
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,              -- 会话ID (Token, hex)
    user_id TEXT NOT NULL,            -- 用户ID
    expires_at DATETIME NOT NULL,     -- 过期时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,              -- 日志ID (UUID)
    user_id TEXT,                     -- 用户ID
    action TEXT NOT NULL,             -- 操作类型: register, login, logout, sign, decrypt, etc.
    detail TEXT,                      -- 操作详情 (JSON)
    ip_address TEXT,                  -- 客户端IP
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_keys_user_id ON keys(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- 清理过期会话的触发器
CREATE TRIGGER IF NOT EXISTS cleanup_expired_sessions
AFTER INSERT ON sessions
BEGIN
    DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
END;
