# API 设计与规范

## 1. 统一响应格式
- 成功：
  - code: 0
  - message: "ok"
  - data: 任意对象
- 失败：
  - code: 非 0
  - message: 错误说明
  - data: null

示例：
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "123",
    "name": "demo"
  }
}

## 2. 错误码规范（示例）
- 1000: 参数错误
- 1001: 未授权
- 1002: 禁止访问
- 1003: 资源不存在
- 2000: 业务校验失败
- 3000: 依赖服务失败
- 9000: 系统内部错误

## 3. 鉴权方案
- 登录成功返回 access_token（JWT）与 refresh_token
- 访问 API 时携带 Authorization: Bearer <token>
- refresh_token 用于刷新 access_token

## 4. 资源模型（示例）
- User: id, username, password_hash, created_at
- Item: id, name, description, status, created_at, updated_at

## 5. REST 接口清单
### 5.1 认证
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh
- POST /api/v1/auth/logout

### 5.2 用户
- GET /api/v1/users
- GET /api/v1/users/{id}
- POST /api/v1/users
- PUT /api/v1/users/{id}
- DELETE /api/v1/users/{id}

### 5.3 资源（Item）
- GET /api/v1/items
- GET /api/v1/items/{id}
- POST /api/v1/items
- PUT /api/v1/items/{id}
- DELETE /api/v1/items/{id}

## 6. 通用查询参数
- page, page_size
- sort_by, order
- filter（可扩展）

## 7. 版本策略
- URL 版本：/api/v1
- 重大变更时升级到 /api/v2
