# 架构设计

## 1. 总体架构
- 单体 API 服务（Go）对外提供 REST 接口
- 数据层使用 PostgreSQL
- 配置与日志独立模块
- 监控指标暴露 /metrics

## 2. 模块划分
- api: 路由与控制器
- service: 业务逻辑层
- repository: 数据访问层
- model: 数据模型与 DTO
- config: 配置加载与校验
- middleware: 认证、日志、错误处理
- observability: 日志、指标、追踪字段

## 3. 数据流
- 客户端 -> API -> service -> repository -> DB
- 中间件统一处理鉴权、日志、错误
- 监控系统采集 /metrics 指标

## 4. 关键依赖
- Web 框架：Gin（可替换为 Echo）
- 数据库驱动：pgx
- 配置：viper
- 认证：JWT
- 迁移：golang-migrate

## 5. 边界与扩展
- 先单体部署，后续可拆分 service
- 预留缓存层（Redis）接口
- 预留任务队列接口（如 Asynq）
