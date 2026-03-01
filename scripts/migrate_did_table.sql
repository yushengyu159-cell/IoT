-- DID表结构迁移脚本
-- 为支持前端注册功能，扩展DID表结构

-- 1. 检查表是否存在
CREATE TABLE IF NOT EXISTS `dids` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(64) DEFAULT NULL,
    `phone` varchar(32) DEFAULT NULL,
    `email` varchar(128) NOT NULL,
    `password` varchar(128) NOT NULL,
    `role` varchar(32) DEFAULT NULL,
    `age` int DEFAULT NULL,
    `created_at` varchar(32) DEFAULT NULL,
    `did` varchar(128) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_did` (`did`),
    KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2. 添加新字段（使用兼容语法）
-- 全名字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'full_name') = 0,
    'ALTER TABLE `dids` ADD COLUMN `full_name` varchar(64) DEFAULT NULL COMMENT "用户全名" AFTER `role`',
    'SELECT "full_name column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 建筑名称字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'building_name') = 0,
    'ALTER TABLE `dids` ADD COLUMN `building_name` varchar(128) DEFAULT NULL COMMENT "建筑名称" AFTER `full_name`',
    'SELECT "building_name column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 建筑地址字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'building_addr') = 0,
    'ALTER TABLE `dids` ADD COLUMN `building_addr` varchar(256) DEFAULT NULL COMMENT "建筑地址" AFTER `building_name`',
    'SELECT "building_addr column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 建筑类型字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'building_type') = 0,
    'ALTER TABLE `dids` ADD COLUMN `building_type` varchar(64) DEFAULT NULL COMMENT "建筑类型" AFTER `building_addr`',
    'SELECT "building_type column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 物业名称字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'property_name') = 0,
    'ALTER TABLE `dids` ADD COLUMN `property_name` varchar(128) DEFAULT NULL COMMENT "物业名称" AFTER `building_type`',
    'SELECT "property_name column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 职业字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'occupation') = 0,
    'ALTER TABLE `dids` ADD COLUMN `occupation` varchar(64) DEFAULT NULL COMMENT "职业" AFTER `property_name`',
    'SELECT "occupation column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 机构名称字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'institution') = 0,
    'ALTER TABLE `dids` ADD COLUMN `institution` varchar(128) DEFAULT NULL COMMENT "机构名称" AFTER `occupation`',
    'SELECT "institution column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 状态字段
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND COLUMN_NAME = 'status') = 0,
    'ALTER TABLE `dids` ADD COLUMN `status` varchar(32) DEFAULT "pending" COMMENT "注册状态：pending, verified, completed" AFTER `institution`',
    'SELECT "status column already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3. 添加索引优化查询性能（使用兼容语法）
-- 角色索引
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND INDEX_NAME = 'idx_role') = 0,
    'CREATE INDEX `idx_role` ON `dids` (`role`)',
    'SELECT "idx_role index already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 状态索引
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND INDEX_NAME = 'idx_status') = 0,
    'CREATE INDEX `idx_status` ON `dids` (`status`)',
    'SELECT "idx_status index already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 建筑名称索引（用于搜索）
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND INDEX_NAME = 'idx_building_name') = 0,
    'CREATE INDEX `idx_building_name` ON `dids` (`building_name`)',
    'SELECT "idx_building_name index already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 机构名称索引（用于搜索）
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dids' AND INDEX_NAME = 'idx_institution') = 0,
    'CREATE INDEX `idx_institution` ON `dids` (`institution`)',
    'SELECT "idx_institution index already exists"'
));
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 4. 更新现有记录的默认值
UPDATE `dids` SET `status` = 'completed' WHERE `status` IS NULL;

-- 5. 显示表结构
DESCRIBE `dids`;

-- 6. 显示索引信息
SHOW INDEX FROM `dids`;
