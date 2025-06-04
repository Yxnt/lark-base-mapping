# GitLab Webhook API 使用指南

## 概述

本项目新增了GitLab webhook API接口，用于处理GitLab的各种事件，特别是Merge Request事件。

## 配置

### 环境变量配置

在`.env`文件中添加以下GitLab相关配置：

```bash
# GitLab Webhook 配置
GITLAB_WEBHOOK_SECRET=your_gitlab_webhook_secret_token
GITLAB_BASE_URL=https://gitlab.com
```

### GitLab项目配置

1. 在GitLab项目中，进入 **Settings > Webhooks**
2. 添加新的Webhook：
   - **URL**: `https://your-domain.com/webhook/gitlab`
   - **Secret Token**: 与环境变量`GITLAB_WEBHOOK_SECRET`相同
   - **Trigger**: 选择要监听的事件，如：
     - ✅ Merge request events
     - ✅ Push events
     - ✅ Tag push events
     - ✅ Issues events

## API 端点

### POST /webhook/gitlab

处理GitLab webhook事件的主要端点。

**请求头：**
- `Content-Type`: `application/json`
- `X-Gitlab-Event`: 事件类型（如：`Merge Request Hook`）
- `X-Gitlab-Token`: webhook密钥（可选，如果配置了`GITLAB_WEBHOOK_SECRET`）
- `X-Gitlab-Instance`: GitLab实例URL
- `X-Gitlab-Event-UUID`: 事件唯一标识符

## 支持的事件类型

### 1. Merge Request Hook
处理Merge Request相关事件，包括：
- `open` - 创建MR
- `update` - 更新MR
- `merge` - 合并MR
- `close` - 关闭MR
- `reopen` - 重新打开MR

**数据存储：**
MR事件会自动保存到`gitlab_merge_requests`表中，包含以下信息：
- MR ID和IID
- 标题和描述
- 状态和操作
- 作者信息
- 项目信息
- 源分支和目标分支
- MR URL
- 完整事件数据（JSON格式）

### 2. Push Hook
处理代码推送事件（待实现具体逻辑）

### 3. Tag Push Hook
处理标签推送事件（待实现具体逻辑）

### 4. Issues Hook
处理Issue相关事件（待实现具体逻辑）

## 响应格式

### 成功响应
```json
{
  "status": "success",
  "message": "Merge request event processed",
  "event": {
    "action": "open",
    "mr_id": 123,
    "title": "Fix bug in authentication",
    "state": "opened",
    "project": "my-project"
  }
}
```

### 错误响应
```json
{
  "status": "error",
  "message": "Invalid GitLab webhook token"
}
```

## 安全性

### Webhook密钥验证
- 配置`GITLAB_WEBHOOK_SECRET`环境变量来启用密钥验证
- GitLab会在`X-Gitlab-Token`头中发送密钥
- 如果密钥不匹配，请求会被拒绝（401 Unauthorized）

### 请求验证
- 验证`Content-Type`必须为`application/json`
- 验证必须包含`X-Gitlab-Event`头
- 记录所有webhook请求的详细信息用于调试

## 数据库

### gitlab_merge_requests 表结构

| 字段 | 类型 | 描述 |
|------|------|------|
| mr_id | Number | GitLab MR ID |
| mr_iid | Number | 项目内MR序号 |
| title | Text | MR标题 |
| description | Text | MR描述 |
| state | Text | MR状态 |
| action | Text | 触发的操作 |
| author_name | Text | 作者姓名 |
| author_username | Text | 作者用户名 |
| project_id | Number | 项目ID |
| project_name | Text | 项目名称 |
| source_branch | Text | 源分支 |
| target_branch | Text | 目标分支 |
| url | URL | MR链接 |
| event_data | JSON | 完整事件数据 |

## 扩展功能

你可以在相应的处理函数中添加自定义业务逻辑：

### handleMergeRequestEvent
- 发送通知到飞书
- 触发CI/CD流程
- 自动代码审查
- 更新项目管理工具

### handlePushEvent
- 构建和部署
- 代码质量检查
- 安全扫描

### handleIssuesEvent
- 自动分配
- 通知相关人员
- 更新工作流

## 日志和监控

所有webhook事件都会记录详细日志，包括：
- 事件类型和来源
- 处理结果
- 错误信息（如果有）
- 性能指标

查看日志来调试webhook集成：
```bash
# 查看应用日志
docker logs your-container-name

# 或者本地运行时查看控制台输出
```

## 测试

### 本地测试
1. 使用ngrok或类似工具暴露本地端口
2. 在GitLab中配置webhook指向ngrok URL
3. 在项目中创建MR或执行其他操作来触发webhook

### 手动测试
使用curl发送测试请求：
```bash
curl -X POST http://localhost:8080/webhook/gitlab \
  -H "Content-Type: application/json" \
  -H "X-Gitlab-Event: Merge Request Hook" \
  -H "X-Gitlab-Token: your_secret_token" \
  -d '{"object_kind":"merge_request","event_type":"merge_request",...}'
```

## 故障排查

### 常见问题

1. **401 Unauthorized**
   - 检查`GITLAB_WEBHOOK_SECRET`配置
   - 确认GitLab中的Secret Token设置正确

2. **400 Bad Request**
   - 检查Content-Type是否为`application/json`
   - 确认包含必要的请求头

3. **数据库错误**
   - 确认已运行数据库迁移
   - 检查`gitlab_merge_requests`表是否存在

4. **事件未处理**
   - 检查`X-Gitlab-Event`头的值
   - 确认事件类型在支持列表中

### 调试步骤
1. 检查应用日志
2. 验证环境变量配置
3. 测试webhook连接性
4. 检查GitLab webhook配置
5. 验证数据库状态 