# CI/CD 与部署

## 1. 流水线阶段
- Lint（golangci-lint）
- Test（单元/集成）
- Build（二进制或 Docker 镜像）
- Deploy（dev/staging/prod）

## 2. 环境划分
- dev：本地或开发环境
- staging：预发环境
- prod：生产环境

## 3. 部署方式
- Docker + docker-compose
- 生产可改为 Kubernetes

## 4. 回滚策略
- 保留上一个稳定版本镜像
- 发布失败自动回滚
- 数据库迁移先备份再执行

## 5. 版本与变更记录
- 使用语义化版本
- 发布说明记录变更
