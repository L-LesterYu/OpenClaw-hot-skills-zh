# ByteRover OpenClaw 技能

ByteRover CLI for OpenClaw - 上下文工程平台，为 AI 编码助手管理项目知识、决策和模式。

## 简介

ByteRover 是一个上下文工程平台，通过以下方式消除内存碎片化：
- 📚 **整理上下文** - 从项目中提取模式、决策和经验
- 🤖 **AI 助手集成** - 支持 19+ 种 AI 编码助手
- 🔄 **自动同步** - 上下文在团队成员之间自动同步
- 🚫 **无 markdown 杂乱** - 零手动上下文管理文件

本技能使用 Docker 在隔离环境中运行 ByteRover CLI，即使在无头服务器环境中也能正常工作。

## 为什么选择 ByteRover？

**防止内存碎片化：**
- 集中式知识库，而非分散的 markdown 文件
- 结构化上下文树（架构、测试、部署模式）
- 基于查询的检索，而非文件搜索

**团队协作：**
- 所有团队成员共享上下文
- 助手使用相同的知识库工作
- 不再有"知识孤岛"

**零维护：**
- 无需手动管理 markdown 文件
- 自动组织和标记
- 始终同步且最新

## 安装

### 前提条件

- Docker
- Docker Compose
- ByteRover 账户（含团队和空间）

### 设置

1. 将此技能目录复制到你的 OpenClaw 工作区：
   ```bash
   cp -r byterover-zh ~/.openclaw/workspace/skills/
   ```

2. 创建 API 配置文件：
   ```bash
   cat > ~/.clawdbot/byterover-config.json << 'EOF'
   {
     "apiKey": "你的_byterover_api_密钥",
     "team": "你的团队名称",
     "space": "你的空间名称"
   }
   EOF
   ```

3. 构建并启动 Docker 容器：
   ```bash
   cd ~/.openclaw/workspace/skills/byterover-zh
   docker-compose build
   docker-compose up -d
   ```

## 使用方法

本技能为常用操作提供了 Bash 辅助脚本：

### 检查状态
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/status.sh
```

### 查询上下文
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/query.sh "测试策略是什么？"
```

### 添加上下文
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/curate.sh "记录 API 认证模式"
```

### 同步更改
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/sync.sh
```

## Docker 命令

### 直接使用 Docker

**使用 API 密钥登录：**
```bash
cd ~/.openclaw/workspace/skills/byterover-zh
BRV_API_KEY=你的密钥 docker-compose exec byterover brv login --api-key $BRV_API_KEY
```

**初始化项目：**
```bash
docker-compose exec byterover brv init --headless --team "我的团队" --space "我的空间"
```

**查询知识：**
```bash
docker-compose exec byterover brv query "我们如何处理认证？" --headless
```

**添加上下文：**
```bash
docker-compose exec byterover brv curate "记录限流策略" --headless
```

**与远程同步：**
```bash
docker-compose exec byterover brv pull --headless  # 拉取最新
docker-compose exec byterover brv push --headless  # 推送本地更改
```

### 容器管理

**查看日志：**
```bash
cd ~/.openclaw/workspace/skills/byterover-zh
docker-compose logs -f byterover
```

**停止容器：**
```bash
docker-compose down
```

**重启容器：**
```bash
docker-compose restart
```

**进入容器 shell：**
```bash
docker-compose exec byterover bash
```

## 上下文树结构

ByteRover 将知识组织到结构化的领域中：

- **Structure（结构）** - 项目架构和设计模式
- **Database（数据库）** - 模式、迁移、关系
- **Backend（后端）** - API 端点、业务逻辑
- **Frontend（前端）** - UI 组件、状态管理
- **Testing（测试）** - 测试策略、测试数据、模式
- **Deployment（部署）** - 基础设施、CI/CD
- **Documentation（文档）** - 指南、API 文档

## 使用场景

### 为 AI 编码助手

**给 Claude Code 的示例提示：**
> 使用 brv query 检查这个项目使用了什么认证模式

**助手执行：**
```bash
brv query "使用了什么认证模式？" --headless --format json
```

### 为开发者

**检查项目架构：**
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/query.sh "项目架构是什么？"
```

**添加决策文档：**
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/curate.sh "我们选择 PostgreSQL 而非 MongoDB，因为有关系数据需求"
```

**查看测试策略：**
```bash
~/.openclaw/workspace/skills/byterover-zh/scripts/query.sh "我们的测试策略是什么？"
```

## 配置

配置存储在 `~/.clawdbot/byterover-config.json`：

```json
{
  "apiKey": "你的_api_密钥",
  "team": "你的团队名称",
  "space": "你的空间名称"
}
```

### 环境变量

- `BRV_API_KEY` - ByteRover API 密钥
- `BRV_TEAM` - 团队名称
- `BRV_SPACE` - 空间名称
- `BRV_HEADLESS` - 启用无头模式
- `BRV_FORMAT` - 输出格式（json/text）

## 无头模式

ByteRover CLI 支持无头模式用于自动化：
- 无交互式提示
- 机器可读的 JSON 输出
- CI/CD 友好
- 非常适合 OpenClaw 自动化

**支持的命令：**
- `brv init` --headless --team X --space Y
- `brv status` --headless
- `brv query` --headless
- `brv curate` --headless
- `brv push` --headless
- `brv pull` --headless

## 与 AI 助手集成

ByteRover 支持 19+ 种 AI 编码助手：

### 现代集成
- **Claude Code** - 技能文件
- **Cursor** - 技能文件
- **Windsurf** - MCP 工具

### 旧版集成
- **GitHub Copilot** - 基于规则
- **Tabnine** - MCP 工具
- **Continue** - MCP 工具

## Docker 架构

**为什么使用 Docker？**

ByteRover CLI 需要：
- Node.js 20+
- libsecret 用于凭证存储
- Python 用于原生依赖
- 桌面环境用于 OAuth

**Docker 解决方案：**
- 包含所有依赖的隔离环境
- 在无头服务器上工作
- 无需 OAuth（使用 API 密钥认证）
- 持久化凭证存储

**Dockerfile 包含：**
- Node.js 20 基础镜像
- 原生模块的构建工具
- libsecret 用于凭证存储
- 预安装的 ByteRover CLI

## 故障排除

**容器无法启动：**
```bash
# 检查端口是否可用
docker-compose ps
# 查看日志
docker-compose logs
```

**构建失败：**
```bash
# 不使用缓存重新构建
docker-compose build --no-cache
```

**登录失败：**
```bash
# 验证 API 密钥
cat ~/.clawdbot/byterover-config.json | jq '.apiKey'
# 检查团队/空间名称
```

**无上下文返回：**
```bash
# 初始化项目
docker-compose exec byterover brv init --headless --team "团队" --space "空间"
# 拉取最新上下文
docker-compose exec byterover brv pull --headless
```

## 系统要求

- Docker
- Docker Compose
- ByteRover 账户
- 活跃的团队和空间

## 许可证

MIT

## 作者

为 OpenClaw 创建，作者：Jun Suzuki

---

更多关于 ByteRover 的信息，请访问：https://docs.byterover.dev/
