# go-plane

[English](./README.md)

[Plane](https://plane.so/) 与 [GitHub](https://docs.github.com/en/actions) 或 [Gitea Action](https://docs.gitea.com/usage/actions/overview) 集成的项目管理工具。

## 目标

本项目旨在将 Plane 的项目管理功能与 GitHub Actions 或 Gitea Actions 集成，使开发者能够在 CI/CD 过程中根据提交消息自动更新问题、添加评论、更改状态和分配用户。

## 安装

### 从源代码构建

```bash
# 克隆代码库
git clone https://github.com/GeekWorkCode/go-plane
cd go-plane

# 构建二进制文件
make build

# 安装到 $GOPATH/bin
make install
```

### 使用预编译二进制文件

从 [GitHub Releases](https://github.com/GeekWorkCode/go-plane/releases) 下载适合您平台的最新版本。

```bash
# Linux (amd64)
curl -L https://github.com/GeekWorkCode/go-plane/releases/download/vX.Y.Z/go-plane-vX.Y.Z-linux-amd64 -o go-plane
chmod +x go-plane
sudo mv go-plane /usr/local/bin/

# macOS (amd64)
curl -L https://github.com/GeekWorkCode/go-plane/releases/download/vX.Y.Z/go-plane-vX.Y.Z-darwin-amd64 -o go-plane
chmod +x go-plane
sudo mv go-plane /usr/local/bin/
```

### Docker

```bash
# 使用 GitHub Container Registry
docker pull ghcr.io/GeekWorkCode/go-plane:latest

# 运行
docker run --rm \
  -e PLANE_BASE_URL="https://plane.example.com/api/v1" \
  -e PLANE_TOKEN="your-token" \
  -e PLANE_WORKSPACE_SLUG="your-workspace" \
  -e PLANE_REF="PROJ-123 修复问题" \
  ghcr.io/GeekWorkCode/go-plane
```

## 使用方法

### 配置

您可以使用环境变量配置 `go-plane`：

| 环境变量             | 描述                                            | 示例                             |
| -------------------- | ----------------------------------------------- | -------------------------------- |
| PLANE_BASE_URL       | 您的 Plane 实例 API 的基础 URL                  | https://plane.example.com/api/v1 |
| PLANE_TOKEN          | 您的 Plane API 令牌                             | plane_api_123456789              |
| PLANE_WORKSPACE_SLUG | 您的 Plane 工作区标识                           | my-workspace                     |
| PLANE_REF            | 要分析的 git 提交消息或引用                     | "修复问题 PROJ-123: 更新 API"    |
| PLANE_TO_STATE       | 问题要转换的状态（可选）                        | "已完成"                         |
| PLANE_COMMENT        | 要添加到问题的评论（可选）                      | "在提交中修复: $GITHUB_SHA"      |
| PLANE_ASSIGNEE       | 要分配问题的用户名（可选）                      | "john.doe"                       |
| PLANE_MARKDOWN       | 设置为 "true" 以将评论格式化为 Markdown（可选） | true                             |
| PLANE_INSECURE       | 设置为 "true" 以跳过 SSL 验证（不推荐）         | false                            |
| PLANE_DEBUG          | 设置为 "true" 以启用调试输出                    | false                            |

### GitHub Actions 示例

```yaml
name: Plane 集成

on:
  push:
    branches: [ main ]

jobs:
  plane-integration:
    runs-on: ubuntu-latest
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: 更新 Plane 问题
        uses: docker://ghcr.io/GeekWorkCode/go-plane:latest
        env:
          PLANE_BASE_URL: ${{ secrets.PLANE_BASE_URL }}
          PLANE_TOKEN: ${{ secrets.PLANE_TOKEN }}
          PLANE_WORKSPACE_SLUG: ${{ secrets.PLANE_WORKSPACE_SLUG }}
          PLANE_REF: ${{ github.event.head_commit.message }}
          PLANE_TO_STATE: "已完成"
          PLANE_COMMENT: "在提交中解决: ${{ github.sha }}\n\n作者: ${{ github.actor }}"
          PLANE_MARKDOWN: "true"
```

### Gitea Actions 示例

```yaml
name: Plane 集成

on:
  push:
    branches: [ main ]

jobs:
  plane-integration:
    runs-on: ubuntu-latest
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: 更新 Plane 问题
        uses: docker://ghcr.io/GeekWorkCode/go-plane:latest
        env:
          PLANE_BASE_URL: ${{ secrets.PLANE_BASE_URL }}
          PLANE_TOKEN: ${{ secrets.PLANE_TOKEN }}
          PLANE_WORKSPACE_SLUG: ${{ secrets.PLANE_WORKSPACE_SLUG }}
          PLANE_REF: ${{ github.event.head_commit.message }}
          PLANE_TO_STATE: "已完成"
          PLANE_COMMENT: "在提交中解决: ${{ github.sha }}\n\n作者: ${{ github.actor }}"
          PLANE_MARKDOWN: "true"
```

## 常见问题

### 问题格式

工具会提取格式为 "PROJ-123" 的问题键，其中 "PROJ" 是项目标识符，"123" 是问题序列ID。问题键可以位于提交消息的任何位置，例如：

```
修复登录表单中的错误 PROJ-123
```

或者

```
PROJ-123 实现新功能
```

### API 连接问题

如果您遇到 SSL 证书问题，可以使用 `PLANE_INSECURE=true` 跳过验证（不推荐在生产环境中使用）。

确保您的 Plane API URL 格式正确，通常为 `https://your-plane-instance.com/api/v1`。

### 找不到要更新的问题

- 确保您的工作区标识 (`PLANE_WORKSPACE_SLUG`) 正确
- 检查提交消息中是否包含有效的问题引用
- 确认您的 API 令牌具有访问这些问题的权限

## 许可证

MIT

## 致谢

本项目的灵感来源于：
- [go-jira](https://github.com/appleboy/go-jira) 的结构和方法
- [plane-api-go](https://github.com/GeekWorkCode/plane-api-go) 的 Plane API 客户端 