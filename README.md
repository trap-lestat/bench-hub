# Bench Hub 压测项目

## 项目简介
- Go 后端 API + Locust 压测 + Vue 管理后台
- 内置 JWT 认证、用户/脚本/任务/报告管理
- 支持 Prometheus 指标采集与 Docker 部署

## 技术栈
- 后端：Go + Gin + PostgreSQL
- 压测：Locust
- 前端：Vue 3 + Vite + Pinia + Vue Router
- 监控：Prometheus client

## 目录结构
- `cmd/server`：后端启动入口
- `internal`：业务代码（API/Service/Repo/Config/Middleware）
- `migrations`：数据库迁移 SQL
- `scripts`：脚本（回归、seed、压测）
- `locust`：Locust 脚本
- `web`：Vue 管理后台
- `docs`：计划与设计文档

## 后端启动
1. 配置数据库并执行迁移
   - 迁移 SQL 在 `migrations/`
   - 可使用 psql 手动执行：
     - `psql -h 127.0.0.1 -U postgres -d bench_hub -f migrations/0001_create_users.up.sql`
   - 可选自动迁移：设置 `AUTO_MIGRATE=true` + `MIGRATIONS_PATH`
2. 启动服务
   - `go run ./cmd/server`

## 环境变量（后端）
- `PORT`：服务端口（默认 8080）
- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASS`、`DB_NAME`、`DB_SSLMODE`
- `JWT_SECRET`、`JWT_ISSUER`
- `ACCESS_TOKEN_MINUTES`、`REFRESH_TOKEN_DAYS`
- `LOCUST_BIN`、`LOCUST_HOST`、`REPORTS_DIR`
- `MIGRATIONS_PATH`、`AUTO_MIGRATE`
- `RUNNER_URL`（使用独立 runner 容器时）

## 真实数据初始化
- 执行：
  - `DB_HOST=... DB_USER=... DB_PASS=... DB_NAME=... scripts/seed.sh`
- 默认会创建管理员账号：`admin / admin123`

## 管理员初始化
- 一键初始化/重置管理员账号（默认 `admin / admin123`）：
  - `scripts/init-admin.sh`
- 使用本地数据库连接（需要 `psql`）：
  - `DB_HOST=127.0.0.1 DB_USER=postgres DB_PASS=postgres DB_NAME=bench_hub scripts/init-admin.sh`
- 自定义账号/密码（bcrypt hash）：
  - `ADMIN_USERNAME=admin ADMIN_PASSWORD_HASH=... scripts/init-admin.sh`

## 前端启动
1. 进入前端目录：`cd web`
2. 安装依赖：`npm install`
3. 启动：`npm run dev`
4. 如后端不是本机 8080，设置：`VITE_API_PROXY_TARGET`

## 压测与报告
- 交互式运行：
  - `locust -f locust/locustfile.py --host http://localhost:8080`
- 自动生成 CSV/HTML 报告：
  - `HOST=http://localhost:8080 USERS=50 SPAWN_RATE=5 DURATION=5m scripts/run-loadtest.sh`
- 报告输出到 `reports/`
- 报告下载接口：`/api/v1/reports/{id}/download`

## 监控指标
- Prometheus 指标：`/metrics`

## Docker 部署
- `docker compose up --build`
- 管理后台：`http://localhost:5173`
 - 前端通过 Vite 反向代理访问后端 `/api`
 - Runner 服务执行 Locust，API 通过 `RUNNER_URL` 调用

## API 回归脚本
- `BASE_URL=http://localhost:8080 scripts/api-regression.sh`

## 文档
- 项目计划与设计：`docs/`
