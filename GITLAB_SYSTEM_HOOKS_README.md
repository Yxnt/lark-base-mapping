# GitLab System Hooks 集成指南

## 概述

本项目现已支持完整的 GitLab System Hooks 事件处理。System Hooks 是 GitLab 管理员级别的 webhooks，能够监听整个 GitLab 实例的系统级事件，如项目创建、用户管理、组变更等。

## 支持的事件类型

### 传统格式事件

#### 项目相关事件
- `project_create` - 项目创建
- `project_destroy` - 项目删除  
- `project_rename` - 项目重命名
- `project_transfer` - 项目转移
- `project_update` - 项目更新

#### 用户相关事件  
- `user_create` - 用户创建
- `user_destroy` - 用户删除
- `user_rename` - 用户重命名
- `user_failed_login` - 用户登录失败（被封禁用户）

#### 组相关事件
- `group_create` - 组创建
- `group_destroy` - 组删除
- `group_rename` - 组重命名

#### 访问请求事件
- `user_access_request_revoked_for_group` - 撤销组访问请求
- `user_access_request_revoked_for_project` - 撤销项目访问请求
- `user_access_request_to_group` - 组访问请求
- `user_access_request_to_project` - 项目访问请求
- `user_add_to_group` - 用户加入组
- `user_add_to_team` - 用户加入团队
- `user_remove_from_group` - 用户移出组
- `user_remove_from_team` - 用户移出团队
- `user_update_for_group` - 组用户信息更新
- `user_update_for_team` - 团队用户信息更新

#### 密钥事件
- `key_create` - SSH密钥创建
- `key_destroy` - SSH密钥删除

#### 仓库事件
- `repository_update` - 仓库更新（推送代码、标签等）

### 新格式事件

#### 成员审批事件
- `gitlab_subscription_member_approval` (action: `enqueue`) - 成员升级申请排队
- `gitlab_subscription_member_approvals` (action: `approve`/`deny`) - 成员升级审批

## 数据库表结构

系统会自动创建以下数据库表来存储不同类型的事件：

### 1. gitlab_project_system_events
存储项目相关的系统事件
- `event_name` - 事件名称
- `project_id` - 项目ID
- `project_name` - 项目名称
- `path_with_namespace` - 项目完整路径
- `project_visibility` - 项目可见性
- `owner_name` / `owner_email` - 项目所有者信息
- `old_path_with_namespace` - 旧路径（重命名/转移时）

### 2. gitlab_user_system_events  
存储用户相关的系统事件
- `event_name` - 事件名称
- `user_id` - 用户ID
- `user_name` - 用户名
- `user_email` - 用户邮箱
- `user_username` - 用户名
- `old_username` - 旧用户名（重命名时）

### 3. gitlab_group_system_events
存储组相关的系统事件
- `event_name` - 事件名称
- `group_id` - 组ID
- `group_name` - 组名称
- `path_with_namespace` - 组完整路径
- `old_path_with_namespace` - 旧路径（重命名时）

### 4. gitlab_access_request_events
存储访问请求相关事件
- `event_name` - 事件名称
- `user_id` - 用户ID
- 组相关字段：`group_id`, `group_name`, `group_access`
- 项目相关字段：`project_id`, `project_name`, `project_access`

### 5. gitlab_key_events
存储SSH密钥相关事件
- `event_name` - 事件名称
- `user_id` - 用户ID
- `key_id` - 密钥ID

### 6. gitlab_repository_update_events
存储仓库更新事件
- `event_name` - 事件名称
- `user_id` - 用户ID
- `project_id` - 项目ID
- `refs` - 引用信息（JSON格式）
- `changes` - 变更信息（JSON格式）

### 7. gitlab_member_approval_events
存储成员审批事件（新格式）
- `object_kind` - 对象类型
- `action` - 操作类型
- `user_id` - 用户ID
- `status` - 状态
- `new_access_level` / `old_access_level` - 访问级别

## 配置说明

### 环境变量
确保设置以下环境变量：
```bash
GITLAB_WEBHOOK_SECRET=your_secret_token
GITLAB_BASE_URL=https://your-gitlab-instance.com
```

### GitLab 配置
1. 以管理员身份登录 GitLab
2. 进入 **Admin Area** > **System Hooks**
3. 点击 **Add new webhook**
4. 配置以下参数：
   - **URL**: `https://your-domain.com/webhook/gitlab`
   - **Secret Token**: 与环境变量 `GITLAB_WEBHOOK_SECRET` 相同
   - **Trigger**: 勾选所有需要的事件类型
   - **Enable SSL verification**: 根据需要启用

## 事件处理流程

1. **事件接收**: System Hook 事件通过 `POST /webhook/gitlab` 端点接收
2. **安全验证**: 中间件验证 webhook 密钥和请求头
3. **事件分发**: 根据 `X-Gitlab-Event: System Hook` 头识别为系统事件
4. **类型识别**: 
   - 新格式事件：根据 `object_kind` 字段分发
   - 传统格式事件：根据 `event_name` 字段分发
5. **数据处理**: 解析事件数据并保存到对应的数据库表
6. **响应返回**: 返回处理结果给 GitLab

## 事件处理示例

### 项目创建事件
```json
{
  "event_name": "project_create",
  "project_id": 74,
  "name": "StoreCloud",
  "path_with_namespace": "jsmith/storecloud",
  "project_visibility": "private",
  "owner_name": "John Smith",
  "owner_email": "johnsmith@example.com"
}
```

### 用户创建事件
```json
{
  "event_name": "user_create",
  "user_id": 41,
  "user_name": "John Smith",
  "user_email": "johnsmith@example.com",
  "user_username": "johnsmith"
}
```

### 仓库更新事件
```json
{
  "event_name": "repository_update",
  "user_id": 1,
  "project_id": 1,
  "changes": [
    {
      "before": "8205ea8d81ce0c6b90fbe8280d118cc9fdad6130",
      "after": "4045ea7a3df38697b3730a20fb73c8bed8a3e69e",
      "ref": "refs/heads/master"
    }
  ]
}
```

## 扩展和自定义

### 添加业务逻辑
在各个事件处理函数中，您可以添加自定义的业务逻辑：

```go
// 在 handleProjectSystemEvent 中添加
if eventName == "project_create" {
    // 项目创建时的自定义逻辑
    // 例如：发送飞书通知、创建相关资源等
}
```

### 集成飞书通知
利用现有的飞书集成功能，可以在系统事件发生时发送通知：

```go
// 引用现有的飞书中间件
import "gitlab.yogorobot.com/sre/lark-base-mapping/middlewares"

// 在事件处理中发送通知
// 具体实现可参考现有的 lark.go 代码
```

### 添加新的事件类型
如果 GitLab 添加了新的 System Hook 事件类型：

1. 在 `handleTraditionalSystemEvent` 函数中添加新的 case
2. 创建对应的数据结构
3. 实现处理函数
4. 创建对应的数据库表（如需要）

## 监控和日志

所有 System Hook 事件都会记录详细的日志：

- 事件接收日志
- 处理过程日志  
- 数据库保存结果
- 错误信息（如有）

可以通过应用日志监控 System Hook 事件的处理情况。

## 故障排查

### 常见问题

1. **事件未触发**
   - 检查 GitLab System Hook 配置
   - 验证 webhook URL 可访问性
   - 确认勾选了相应的事件类型

2. **签名验证失败**
   - 检查 `GITLAB_WEBHOOK_SECRET` 环境变量
   - 确认 GitLab 配置的 Secret Token 一致

3. **数据库保存失败**
   - 检查数据库表是否已创建
   - 运行迁移：`your-app migrate`
   - 查看应用日志获取详细错误信息

### 调试模式
启用详细日志以便调试：
```bash
LOG_LEVEL=debug your-app serve
```

## 最佳实践

1. **安全性**
   - 始终使用 HTTPS
   - 设置复杂的 webhook 密钥
   - 定期轮换密钥

2. **性能**
   - 异步处理重负载操作
   - 合理设置数据库索引
   - 监控处理性能

3. **可靠性**  
   - 实现幂等性处理
   - 添加重试机制
   - 记录完整的审计日志

4. **扩展性**
   - 模块化事件处理逻辑
   - 使用消息队列处理大量事件
   - 考虑事件的优先级处理

## 相关文档

- [GitLab System Hooks 官方文档](https://docs.gitlab.com/administration/system_hooks/)
- [项目现有的 GitLab Webhook 集成指南](./GITLAB_WEBHOOK_README.md)
- [飞书集成文档](./LARK_README.md) 