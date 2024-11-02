-- 数据表配置
CREATE TABLE IF NOT EXISTS config_tables (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    mysql_table_name VARCHAR(100) NOT NULL,  -- 对应的MySQL表名
    fields JSON NOT NULL,
    indexes JSON,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_app_code (application_id, code),
    UNIQUE KEY uk_app_mysql_table (application_id, mysql_table_name),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 维度配置
CREATE TABLE IF NOT EXISTS config_dimensions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT 'time:时间维度 enum:枚举维度 range:范围维度',
    mysql_table_name VARCHAR(100) NOT NULL,  -- 对应的MySQL表名
    configuration JSON NOT NULL,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_app_code (application_id, code),
    UNIQUE KEY uk_app_mysql_table (application_id, mysql_table_name),
FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 数据模型配置
CREATE TABLE IF NOT EXISTS config_data_models (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    table_id BIGINT UNSIGNED NOT NULL,
    fields JSON NOT NULL,
    dimensions JSON,
    metrics JSON NOT NULL,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_app_code (application_id, code),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (table_id) REFERENCES config_tables(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 表单配置
CREATE TABLE IF NOT EXISTS config_forms (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL COMMENT 'create:新建表单 edit:编辑表单 view:查看表单',
    table_id BIGINT UNSIGNED NOT NULL,
    layout JSON NOT NULL,
    fields JSON NOT NULL,
    rules JSON,
    events JSON,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_app_code (application_id, code),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (table_id) REFERENCES config_tables(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 菜单配置
CREATE TABLE IF NOT EXISTS config_menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    icon VARCHAR(50),
    path VARCHAR(200),
    component VARCHAR(200),
    sort INT NOT NULL DEFAULT 0,
    status TINYINT NOT NULL DEFAULT 1 COMMENT '0:禁用 1:启用',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_app_code (application_id, code),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 配置版本记录
CREATE TABLE IF NOT EXISTS config_versions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    config_type VARCHAR(20) NOT NULL COMMENT 'table:数据表 dimension:维度 model:数据模型 form:表单 menu:菜单',
    config_id BIGINT UNSIGNED NOT NULL,
    version INT NOT NULL,
    content JSON NOT NULL,
    comment TEXT,
    creator_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_config_version (config_type, config_id, version),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
