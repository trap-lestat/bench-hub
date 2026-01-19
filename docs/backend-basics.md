# 后端工程化基础

## 1. 项目结构（建议）
- cmd/server: 启动入口
- internal/api: 路由与控制器
- internal/service: 业务逻辑
- internal/repository: 数据访问
- internal/model: 领域模型与 DTO
- internal/middleware: 认证、日志、错误处理
- internal/config: 配置加载
- internal/observability: 日志、指标、追踪字段
- migrations: 数据库迁移

## 2. 配置管理
- 支持多环境配置（dev/staging/prod）
- 使用 viper 读取 env 与配置文件
- 密钥类通过环境变量注入

## 3. 日志与追踪
- 日志字段统一：time、level、trace_id、request_id
- 入口生成 request_id 并向下传递
- 可选接入 OpenTelemetry

## 4. 错误处理
- 统一错误结构（code/message）
- 中间件捕获 panic 并返回 500
- 业务错误与系统错误区分

## 5. 健康检查与指标
- /health 返回服务状态
- /metrics 暴露 Prometheus 指标

## 6. 编码规范
- 必要注释、函数职责单一
- 避免跨层调用与循环依赖
- 关键业务逻辑必须有单元测试
