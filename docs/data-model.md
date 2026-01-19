# 数据模型与存储设计

## 1. 数据库选型
- PostgreSQL 15+（开源、稳定、索引与事务支持完善）

## 2. 表结构设计（示例）

### users
- id: uuid (PK)
- username: varchar(64) unique not null
- password_hash: varchar(128) not null
- created_at: timestamp not null default now()

索引：
- unique(username)

### items
- id: uuid (PK)
- name: varchar(128) not null
- description: text
- status: varchar(32) not null
- created_at: timestamp not null default now()
- updated_at: timestamp not null default now()

索引：
- idx_items_status (status)
- idx_items_created_at (created_at)

## 3. 迁移策略
- 使用 golang-migrate
- 迁移文件按时间戳命名
- 每次变更必须包含 up/down

## 4. 初始化与样例数据
- 初始化管理员用户（配置化）
- 准备少量 items 样例数据

## 5. 约束与扩展
- 保留 tenant_id 字段扩展多租户（可选）
- 预留软删除字段 deleted_at（可选）
