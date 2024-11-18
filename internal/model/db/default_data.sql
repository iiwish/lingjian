-- 创建默认管理员账号
INSERT INTO sys_users (username, nickname, password, email, status, created_at, updated_at)
VALUES ('admin', '超级管理员', '4336fe705174a646e5ffd19d5976421088ec707c30e036ffc32f99190af113ce', 'admin@lingjian.com', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 创建默认角色
INSERT INTO sys_roles (name, code, status, created_at, updated_at)
VALUES 
    ('超级管理员', 'super_admin', 1, NOW(), NOW()),
    ('应用管理员', 'app_admin', 1, NOW(), NOW()),
    ('普通用户', 'normal_user', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 创建基本系统权限
INSERT INTO sys_permissions (name, code, type, path, method, status, created_at, updated_at)
VALUES 
    -- 系统管理菜单
    ('系统管理', 'system_manage', 'menu', '/system', '*', 1, NOW(), NOW()),
    ('用户管理', 'user_manage', 'menu', '/system/users', '*', 1, NOW(), NOW()),
    ('角色管理', 'role_manage', 'menu', '/system/roles', '*', 1, NOW(), NOW()),
    ('权限管理', 'permission_manage', 'menu', '/system/permissions', '*', 1, NOW(), NOW()),
    ('应用管理', 'app_manage', 'menu', '/system/apps', '*', 1, NOW(), NOW()),

    -- 系统管理API权限
    ('用户查询', 'user_query', 'api', '/api/v1/users', 'GET', 1, NOW(), NOW()),
    ('用户创建', 'user_create', 'api', '/api/v1/users', 'POST', 1, NOW(), NOW()),
    ('用户更新', 'user_update', 'api', '/api/v1/users/:id', 'PUT', 1, NOW(), NOW()),
    ('用户删除', 'user_delete', 'api', '/api/v1/users/:id', 'DELETE', 1, NOW(), NOW()),
    ('角色查询', 'role_query', 'api', '/api/v1/roles', 'GET', 1, NOW(), NOW()),
    ('角色创建', 'role_create', 'api', '/api/v1/roles', 'POST', 1, NOW(), NOW()),
    ('角色更新', 'role_update', 'api', '/api/v1/roles/:id', 'PUT', 1, NOW(), NOW()),
    ('角色删除', 'role_delete', 'api', '/api/v1/roles/:id', 'DELETE', 1, NOW(), NOW()),
    ('权限查询', 'permission_query', 'api', '/api/v1/permissions', 'GET', 1, NOW(), NOW()),
    ('权限创建', 'permission_create', 'api', '/api/v1/permissions', 'POST', 1, NOW(), NOW()),
    ('权限更新', 'permission_update', 'api', '/api/v1/permissions/:id', 'PUT', 1, NOW(), NOW()),
    ('权限删除', 'permission_delete', 'api', '/api/v1/permissions/:id', 'DELETE', 1, NOW(), NOW()),
    ('应用查询', 'app_query', 'api', '/api/v1/apps', 'GET', 1, NOW(), NOW()),
    ('应用创建', 'app_create', 'api', '/api/v1/apps', 'POST', 1, NOW(), NOW()),
    ('应用更新', 'app_update', 'api', '/api/v1/apps/:id', 'PUT', 1, NOW(), NOW()),
    ('应用删除', 'app_delete', 'api', '/api/v1/apps/:id', 'DELETE', 1, NOW(), NOW()),

    -- 配置中心菜单
    ('配置中心', 'config_center', 'menu', '/config', '*', 1, NOW(), NOW()),
    ('维度配置', 'dimension_manage', 'menu', '/config/dimensions', '*', 1, NOW(), NOW()),
    ('表单配置', 'form_manage', 'menu', '/config/forms', '*', 1, NOW(), NOW()),
    ('菜单配置', 'menu_manage', 'menu', '/config/menus', '*', 1, NOW(), NOW()),
    ('模型配置', 'model_manage', 'menu', '/config/models', '*', 1, NOW(), NOW()),
    ('数据表配置', 'table_manage', 'menu', '/config/tables', '*', 1, NOW(), NOW()),

    -- 维度配置API权限
    ('维度查询', 'dimension_query', 'api', '/api/v1/config/dimensions', 'GET', 1, NOW(), NOW()),
    ('维度创建', 'dimension_create', 'api', '/api/v1/config/dimensions', 'POST', 1, NOW(), NOW()),
    ('维度更新', 'dimension_update', 'api', '/api/v1/config/dimensions/:dim_id', 'PUT', 1, NOW(), NOW()),
    ('维度删除', 'dimension_delete', 'api', '/api/v1/config/dimensions/:dim_id', 'DELETE', 1, NOW(), NOW()),
    ('维度明细查询', 'dimension_item_query', 'api', '/api/v1/config/dimensions/:dim_id/items', 'GET', 1, NOW(), NOW()),
    ('维度明细创建', 'dimension_item_create', 'api', '/api/v1/config/dimensions/:dim_id/items', 'POST', 1, NOW(), NOW()),
    ('维度明细更新', 'dimension_item_update', 'api', '/api/v1/config/dimensions/:dim_id/items/:id', 'PUT', 1, NOW(), NOW()),
    ('维度明细删除', 'dimension_item_delete', 'api', '/api/v1/config/dimensions/:dim_id/items/:id', 'DELETE', 1, NOW(), NOW()),

    -- 表单配置API权限
    ('表单查询', 'form_query', 'api', '/api/v1/config/forms', 'GET', 1, NOW(), NOW()),
    ('表单创建', 'form_create', 'api', '/api/v1/config/forms', 'POST', 1, NOW(), NOW()),
    ('表单更新', 'form_update', 'api', '/api/v1/config/forms/:id', 'PUT', 1, NOW(), NOW()),
    ('表单删除', 'form_delete', 'api', '/api/v1/config/forms/:id', 'DELETE', 1, NOW(), NOW()),

    -- 菜单配置API权限
    ('菜单查询', 'menu_query', 'api', '/api/v1/config/menus', 'GET', 1, NOW(), NOW()),
    ('菜单创建', 'menu_create', 'api', '/api/v1/config/menus', 'POST', 1, NOW(), NOW()),
    ('菜单更新', 'menu_update', 'api', '/api/v1/config/menus/:id', 'PUT', 1, NOW(), NOW()),
    ('菜单删除', 'menu_delete', 'api', '/api/v1/config/menus/:id', 'DELETE', 1, NOW(), NOW()),

    -- 模型配置API权限
    ('模型查询', 'model_query', 'api', '/api/v1/config/models', 'GET', 1, NOW(), NOW()),
    ('模型创建', 'model_create', 'api', '/api/v1/config/models', 'POST', 1, NOW(), NOW()),
    ('模型更新', 'model_update', 'api', '/api/v1/config/models/:id', 'PUT', 1, NOW(), NOW()),
    ('模型删除', 'model_delete', 'api', '/api/v1/config/models/:id', 'DELETE', 1, NOW(), NOW()),

    -- 数据表配置API权限
    ('数据表查询', 'table_query', 'api', '/api/v1/config/tables', 'GET', 1, NOW(), NOW()),
    ('数据表创建', 'table_create', 'api', '/api/v1/config/tables', 'POST', 1, NOW(), NOW()),
    ('数据表更新', 'table_update', 'api', '/api/v1/config/tables/:table_id', 'PUT', 1, NOW(), NOW()),
    ('数据表删除', 'table_delete', 'api', '/api/v1/config/tables/:table_id', 'DELETE', 1, NOW(), NOW()),
    ('数据表明细查询', 'table_item_query', 'api', '/api/v1/config/tables/:table_id/items', 'GET', 1, NOW(), NOW()),
    ('数据表明细创建', 'table_item_create', 'api', '/api/v1/config/tables/:table_id/items', 'POST', 1, NOW(), NOW()),
    ('数据表明细更新', 'table_item_update', 'api', '/api/v1/config/tables/:table_id/items/:id', 'PUT', 1, NOW(), NOW()),
    ('数据表明细删除', 'table_item_delete', 'api', '/api/v1/config/tables/:table_id/items/:id', 'DELETE', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 创建默认菜单
-- INSERT INTO sys_config_menus (name, path, parent_id, sort_order, icon, status, created_at, updated_at)
-- VALUES 
--     -- 系统管理菜单
--     (1, '系统管理', '/system', 0, 1, 'setting', 1, NOW(), NOW()),
--     (2, '用户管理', '/system/users', 1, 1, 'user', 1, NOW(), NOW()),
--     (3, '角色管理', '/system/roles', 1, 2, 'team', 1, NOW(), NOW()),
--     (4, '权限管理', '/system/permissions', 1, 3, 'safety', 1, NOW(), NOW()),
--     (5, '应用管理', '/system/apps', 1, 4, 'appstore', 1, NOW(), NOW()),

--     -- 配置中心菜单
--     (6, '配置中心', '/config', 0, 2, 'tool', 1, NOW(), NOW()),
--     (7, '维度配置', '/config/dimensions', 6, 1, 'apartment', 1, NOW(), NOW()),
--     (8, '表单配置', '/config/forms', 6, 2, 'form', 1, NOW(), NOW()),
--     (9, '菜单配置', '/config/menus', 6, 3, 'menu', 1, NOW(), NOW()),
--     (10, '模型配置', '/config/models', 6, 4, 'database', 1, NOW(), NOW()),
--     (11, '数据表配置', '/config/tables', 6, 5, 'table', 1, NOW(), NOW())
-- ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 为超级管理员角色分配所有权限
INSERT INTO sys_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM sys_roles r, sys_permissions p
WHERE r.code = 'super_admin'
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 为应用管理员分配应用内权限
INSERT INTO sys_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM sys_roles r, sys_permissions p
WHERE r.code = 'app_admin' 
AND p.code IN (
    'dimension_query', 'dimension_create', 'dimension_update', 'dimension_delete',
    'dimension_item_query', 'dimension_item_create', 'dimension_item_update', 'dimension_item_delete',
    'form_query', 'form_create', 'form_update', 'form_delete',
    'menu_query', 'menu_create', 'menu_update', 'menu_delete',
    'model_query', 'model_create', 'model_update', 'model_delete',
    'table_query', 'table_create', 'table_update', 'table_delete',
    'table_item_query', 'table_item_create', 'table_item_update', 'table_item_delete'
)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 为普通用户分配基本查询权限
INSERT INTO sys_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM sys_roles r, sys_permissions p
WHERE r.code = 'normal_user'
AND p.code IN (
    'dimension_query', 'dimension_item_query',
    'form_query',
    'menu_query',
    'model_query',
    'table_query', 'table_item_query'
)
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 为管理员用户分配超级管理员角色
INSERT INTO sys_user_roles (user_id, role_id)
SELECT u.id, r.id
FROM sys_users u, sys_roles r
WHERE u.username = 'admin' AND r.code = 'super_admin'
ON DUPLICATE KEY UPDATE user_id = user_id;
