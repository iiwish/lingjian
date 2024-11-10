-- 创建默认管理员账号
INSERT INTO sys_users (username, password, email, status, created_at, updated_at)
VALUES ('admin', '4336fe705174a646e5ffd19d5976421088ec707c30e036ffc32f99190af113ce', 'admin@lingjian.com', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 创建默认角色
INSERT INTO sys_roles (name, code, status, created_at, updated_at)
VALUES ('超级管理员', 'super_admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 创建基本权限
INSERT INTO sys_permissions (name, code, type, path, method, status, created_at, updated_at)
VALUES 
    ('用户管理', 'user_manage', 'menu', '/users', '*', 1, NOW(), NOW()),
    ('角色管理', 'role_manage', 'menu', '/roles', '*', 1, NOW(), NOW()),
    ('权限管理', 'permission_manage', 'menu', '/permissions', '*', 1, NOW(), NOW()),
    ('应用管理', 'app_manage', 'menu', '/apps', '*', 1, NOW(), NOW()),
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
    ('应用删除', 'app_delete', 'api', '/api/v1/apps/:id', 'DELETE', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 为超级管理员角色分配所有权限
INSERT INTO sys_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM sys_roles r, sys_permissions p
WHERE r.code = 'super_admin'
ON DUPLICATE KEY UPDATE role_id = role_id;

-- 为管理员用户分配超级管理员角色
INSERT INTO sys_user_roles (user_id, role_id)
SELECT u.id, r.id
FROM sys_users u, sys_roles r
WHERE u.username = 'admin' AND r.code = 'super_admin'
ON DUPLICATE KEY UPDATE user_id = user_id;
