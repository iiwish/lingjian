-- 定时任务表
CREATE TABLE IF NOT EXISTS sys_scheduled_tasks (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    app_id      BIGINT UNSIGNED NOT NULL COMMENT '应用ID',
    name        VARCHAR(100) NOT NULL COMMENT '任务名称',
    type        VARCHAR(20) NOT NULL COMMENT '任务类型：sql/http',
    cron        VARCHAR(100) NOT NULL COMMENT 'cron表达式',
    content     TEXT NOT NULL COMMENT '任务内容（JSON格式）',
    timeout     INT NOT NULL DEFAULT 60 COMMENT '超时时间（秒）',
    retry_times INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    status      TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0禁用/1启用/2运行中',
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_app_name (app_id, name) COMMENT '应用ID和任务名称唯一索引',
    KEY idx_status (status) COMMENT '状态索引',
    KEY idx_app_status (app_id, status) COMMENT '应用ID和状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务表' COLLATE=utf8mb4_general_ci;

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS sys_task_logs (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    task_id    BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    status     TINYINT NOT NULL COMMENT '执行状态：0失败/1成功',
    result     TEXT COMMENT '执行结果',
    error      TEXT COMMENT '错误信息',
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    end_time   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
    PRIMARY KEY (id),
    KEY idx_task_time (task_id, start_time) COMMENT '任务ID和开始时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务执行日志表' COLLATE=utf8mb4_general_ci;

-- 元素触发器表
CREATE TABLE IF NOT EXISTS sys_element_triggers (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    app_id        BIGINT UNSIGNED NOT NULL COMMENT '应用ID',
    element_type  VARCHAR(50) NOT NULL COMMENT '元素类型：form/table/model',
    element_id    BIGINT UNSIGNED NOT NULL COMMENT '元素ID',
    trigger_point VARCHAR(20) NOT NULL COMMENT '触发点：before/after',
    type          VARCHAR(20) NOT NULL COMMENT '触发器类型：sql/http',
    content       TEXT NOT NULL COMMENT '触发器内容（JSON格式）',
    status        TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0禁用/1启用',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_element (app_id, element_type, element_id) COMMENT '应用ID、元素类型和元素ID索引',
    KEY idx_status (status) COMMENT '状态索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='元素触发器表' COLLATE=utf8mb4_general_ci;
