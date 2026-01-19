# 管理后台（Vue）设计

## 1. 页面与信息架构
- 登录页
- 首页仪表盘（概览）
- 用户管理（列表/详情/创建/编辑）
- 资源管理（列表/详情/创建/编辑）
- 系统设置

## 2. 路由设计（示例）
- /login
- /dashboard
- /users
- /users/:id
- /items
- /items/:id
- /settings

## 3. 权限与角色
- 角色：admin、operator（可扩展）
- 登录后保存 access_token
- 路由守卫：未登录跳转 /login

## 4. API 契约对齐
- 列表分页：page/page_size
- 排序：sort_by/order
- 过滤：filter
- 统一响应格式与错误码处理

## 5. UI 与组件规范
- 基础组件：表格、分页、表单、弹窗
- 表单校验：必填/长度/格式
- 列表默认加载态与空状态

## 6. 前端工程化
- 状态管理：Pinia
- 网络请求：Axios
- 统一错误提示与拦截器
