-- 创建ESG数据库
CREATE DATABASE IF NOT EXISTS esg CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE esg;

-- 创建ESG文件表
CREATE TABLE IF NOT EXISTS esg_files (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(128) UNIQUE NOT NULL COMMENT 'IPFS CID',
    filename VARCHAR(255) NOT NULL COMMENT '文件名',
    desc_text TEXT COMMENT '文件描述',
    uploader VARCHAR(100) NOT NULL COMMENT '上传者',
    upload_at VARCHAR(50) NOT NULL COMMENT '上传时间',
    txid VARCHAR(255) COMMENT '区块链交易ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_cid (cid),
    INDEX idx_uploader (uploader),
    INDEX idx_upload_at (upload_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ESG文件信息表';

-- 创建DID表
CREATE TABLE IF NOT EXISTS dids (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL COMMENT '姓名',
    phone VARCHAR(32) COMMENT '电话',
    email VARCHAR(128) COMMENT '邮箱',
    password VARCHAR(128) NOT NULL COMMENT '密码',
    role VARCHAR(32) DEFAULT 'user' COMMENT '角色',
    age INT COMMENT '年龄',
    created_at VARCHAR(50) NOT NULL COMMENT '创建时间',
    did VARCHAR(128) UNIQUE NOT NULL COMMENT '数字身份标识',
    created_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_did (did),
    INDEX idx_name (name),
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数字身份表';

-- 创建Fabric资产表（用于存储链码数据）
CREATE TABLE IF NOT EXISTS fabric_assets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    asset_id VARCHAR(100) UNIQUE NOT NULL COMMENT '资产ID',
    color VARCHAR(50) COMMENT '颜色',
    size INT COMMENT '尺寸',
    owner VARCHAR(100) COMMENT '所有者',
    appraised_value INT COMMENT '评估价值',
    txid VARCHAR(255) COMMENT '交易ID',
    timestamp VARCHAR(100) COMMENT '时间戳',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_asset_id (asset_id),
    INDEX idx_owner (owner),
    INDEX idx_txid (txid),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Fabric资产表';

-- 插入测试数据
INSERT INTO esg_files (cid, filename, desc_text, uploader, upload_at) VALUES
('QmTest123', 'test_esg.pdf', '测试ESG文件', 'admin', '2025-08-13 17:00:00'),
('QmTest456', 'sample_report.pdf', '示例报告', 'user1', '2025-08-13 17:01:00');

INSERT INTO dids (name, phone, email, password, role, age, created_at, did) VALUES
('张三', '13800138000', 'zhangsan@example.com', 'password123', 'admin', 30, '2025-08-13 17:00:00', 'did:example:123456789'),
('李四', '13800138001', 'lisi@example.com', 'password456', 'user', 25, '2025-08-13 17:01:00', 'did:example:987654321');

-- 创建用户并授权
CREATE USER IF NOT EXISTS 'fabric'@'%' IDENTIFIED BY 'Test@123456';
GRANT ALL PRIVILEGES ON esg.* TO 'fabric'@'%';
FLUSH PRIVILEGES;
