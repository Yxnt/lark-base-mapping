# Lark Base Mapping

一个基于 PocketBase 的飞书多维表格映射服务。

## 项目概述

这是一个Go语言项目，使用PocketBase作为后端数据库，提供飞书多维表格的映射功能。

## Docker 构建

项目包含 `Dockerfile`，支持容器化部署。

### 本地构建

```bash
# 构建镜像
docker build -t lark-base-mapping .

# 运行容器
docker run -p 8080:8080 lark-base-mapping
```

## GitHub Container Registry

项目已配置 GitHub Actions CI/CD 流水线，会自动构建 Docker 镜像并推送到 GitHub Container Registry。

### 自动构建触发条件

CI 流水线会在以下情况下自动触发：

- **推送到主分支** (`main` 或 `master`)：构建并推送 `latest` 标签
- **创建版本标签** (`v*`)：构建并推送对应版本标签
- **Pull Request**：仅构建，不推送

### 生成的镜像标签

- `ghcr.io/{owner}/{repo}:latest` - 主分支最新版本
- `ghcr.io/{owner}/{repo}:v1.0.0` - 版本标签
- `ghcr.io/{owner}/{repo}:main-{sha}` - 分支+提交哈希

### 使用已构建的镜像

```bash
# 拉取最新镜像
docker pull ghcr.io/{owner}/lark-base-mapping:latest

# 运行容器
docker run -p 8080:8080 ghcr.io/{owner}/lark-base-mapping:latest
```

### 多架构支持

CI 流水线支持构建多架构镜像：
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM64)

## 开发

### 环境要求

- Go 1.23+
- Docker (可选)

### 本地运行

```bash
# 安装依赖
go mod download

# 运行项目
go run main.go
```

### 环境变量

项目使用 `.env` 文件管理环境配置，请根据需要创建并配置。

## CI/CD 配置

### GitHub Actions 工作流

位置：`.github/workflows/ci.yaml`

主要功能：
- 自动检出代码
- 设置 Docker Buildx
- 登录 GitHub Container Registry
- 构建多架构 Docker 镜像
- 推送到仓库

### 权限要求

CI 流水线需要以下权限：
- `contents: read` - 读取仓库内容
- `packages: write` - 推送到 GitHub Container Registry
- `id-token: write` - 身份验证

这些权限已在工作流配置中声明，无需额外设置。

## 部署

### 使用 Docker Compose

```yaml
version: '3.8'
services:
  lark-base-mapping:
    image: ghcr.io/{owner}/lark-base-mapping:latest
    ports:
      - "8080:8080"
    environment:
      - ENV=production
    volumes:
      - ./pb_data:/pb_data
```

### Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lark-base-mapping
spec:
  replicas: 2
  selector:
    matchLabels:
      app: lark-base-mapping
  template:
    metadata:
      labels:
        app: lark-base-mapping
    spec:
      containers:
      - name: lark-base-mapping
        image: ghcr.io/{owner}/lark-base-mapping:latest
        ports:
        - containerPort: 8080
```

## 许可证

本项目采用 [MIT 许可证](LICENSE)。 