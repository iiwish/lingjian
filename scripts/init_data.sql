USE lingjian;

-- 插入默认租户
INSERT INTO tenants (name, code, description, status, admin_email, admin_phone)
VALUES ('默认租户', 'default', '系统默认租户', 1, 'admin@example.com', '13800000000');

-- 获取租户ID
SET @tenant_id = LAST_INSERT_ID();

-- 插入超级管理员用户（密码：admin123，salt：随机生成）
INSERT INTO users (tenant_id, username, password, salt, nickname, email, status)
VALUES (@tenant_id, 'admin', 
        '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918', -- admin123的SHA256
        'K8Cp2X', 'Administrator', 'admin@example.com', 1);

-- 获取用户ID
SET @admin_id = LAST_INSERT_ID();

-- 插入基础角色
INSERT INTO roles (tenant_id, name, code, description, status, is_system)
VALUES 
(@tenant_id, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1, 1),
(@tenant_id, '租户管理员', 'tenant_admin', '租户级别管理员，管理租户内所有应用和用户', 1, 1),
(@tenant_id, '应用管理员', 'app_admin', '应用级别管理员，管理指定应用的所有功能', 1, 1),
(@tenant_id, '普通用户', 'user', '普通用户，根据分配的权限使用系统', 1, 1);

-- 获取超级管理员角色ID
SET @super_admin_role_id = (SELECT id FROM roles WHERE tenant_id = @tenant_id AND code = 'super_admin');

-- 关联超级管理员用户和角色
INSERT INTO user_roles (tenant_id, user_id, role_id)
VALUES (@tenant_id, @admin_id, @super_admin_role_id);

-- 插入基础权限
INSERT INTO permissions (tenant_id, name, code, type, path, method, component, icon, sort, status, is_system, description)
VALUES
-- 系统管理
(@tenant_id, '系统管理', 'system', 'menu', '/system', NULL, 'Layout', 'setting', 1, 1, 1, '系统管理模块'),
(@tenant_id, '用户管理', 'system:user', 'menu', '/system/user', NULL, '/system/user/index', 'user', 1, 1, 1, '用户管理'),
(@tenant_id, '角色管理', 'system:role', 'menu', '/system/role', NULL, '/system/role/index', 'peoples', 2, 1, 1, '角色管理'),
(@tenant_id, '权限管理', 'system:permission', 'menu', '/system/permission', NULL, '/system/permission/index', 'lock', 3, 1, 1, '权限管理'),

-- 租户管理
(@tenant_id, '租户管理', 'tenant', 'menu', '/tenant', NULL, 'Layout', 'apartment', 2, 1, 1, '租户管理模块'),
(@tenant_id, '租户列表', 'tenant:list', 'menu', '/tenant/list', NULL, '/tenant/list/index', 'list', 1, 1, 1, '租户列表'),
(@tenant_id, '租户配置', 'tenant:config', 'menu', '/tenant/config', NULL, '/tenant/config/index', 'tool', 2, 1, 1, '租户配置'),

-- 应用管理
(@tenant_id, '应用管理', 'application', 'menu', '/application', NULL, 'Layout', 'component', 3, 1, 1, '应用管理模块'),
(@tenant_id, '应用列表', 'application:list', 'menu', '/application/list', NULL, '/application/list/index', 'list', 1, 1, 1, '应用列表'),
(@tenant_id, '应用配置', 'application:config', 'menu', '/application/config', NULL, '/application/config/index', 'edit', 2, 1, 1, '应用配置'),
(@tenant_id, '模板管理', 'application:template', 'menu', '/application/template', NULL, '/application/template/index', 'guide', 3, 1, 1, '应用模板管理');

-- 获取所有权限ID
SET @permission_ids = (SELECT GROUP_CONCAT(id) FROM permissions WHERE tenant_id = @tenant_id);

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (tenant_id, role_id, permission_id)
SELECT @tenant_id, @super_admin_role_id, id
FROM permissions
WHERE tenant_id = @tenant_id;

-- 创建示例应用
INSERT INTO applications (tenant_id, name, code, description, status, type, icon, is_template)
VALUES
(@tenant_id, '示例应用', 'demo', '用于演示的示例应用', 1, 'normal', 'example', 0);
