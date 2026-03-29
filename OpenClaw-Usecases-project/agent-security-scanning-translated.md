# 智能体安全扫描 with ClawGuard

## 概述
使用 OpenClaw + [ClawGuard](https://github.com/NeuZhou/clawguard) 在安装前自动扫描技能和插件，保护你的 AI 智能体免受提示注入、隐私泄露和供应链攻击。

## 功能
- **安装前扫描**：在安装任何技能/插件之前进行安全扫描
- **提示注入检测**：25+ 种检测模式，支持多种语言（英语/中文/日语/韩语）
- **隐私信息遮蔽**：自动识别和处理电子邮件、社会保障号码、API 密钥、信用卡号——完全本地运行
- **意图-行为不一致检测**：捕获 AI 智能体说一套做一套的情况

## 痛点
每个使用第三方技能的开发者都面临相同的安全问题：

- **提示注入攻击**：恶意技能可能通过精心设计的提示词操纵智能体的行为
- **隐私信息泄露**：技能可能意外提取并泄露敏感的个人信息
- **供应链攻击**：看似无害的技能可能包含恶意代码或后门
- **行为不一致**：技能声称做一件事，实际却做另一件事

## 工作原理
1. **预扫描验证**：在安装任何新技能之前，ClawGuard 自动扫描其代码和配置
2. **模式匹配检测**：使用 25+ 种预定义的恶意模式检测潜在威胁
3. **隐私信息扫描**：自动扫描代码中的敏感信息模式
4. **风险评估**：基于检测结果生成安全评分和建议
5. **实时监控**：安装后持续监控技能行为，检测异常活动

## 使用方法

### 1. 安装 ClawGuard
```bash
# 克隆仓库
git clone https://github.com/NeuZhou/clawguard.git
cd clawguard

# 安装依赖
pip install -r requirements.txt

# 或者使用 pip 直接安装
pip install clawguard
```

### 2. 配置 OpenClaw
在 OpenClaw 配置中添加 ClawGuard 作为预安装检查：

```json
{
  "mcpServers": {
    "clawguard": {
      "command": "python",
      "args": ["-m", "clawguard.cli"],
      "env": {
        "CLAWGUARD_API_KEY": "your-api-key"
      }
    }
  }
}
```

### 3. 自动扫描
当你尝试安装新技能时，ClawGuard 会自动运行：

```bash
# 尝试安装技能
openclaw skills install some-skill

# ClawGuard 会自动扫描并报告结果
[ClawGuard] Scanning skill: some-skill
[ClawGuard] Security score: 85/100
[ClawGuard] Issues found:
  - Potential prompt injection in prompts/main.txt
  - Privacy concern: API key detected in config.py
[ClawGuard] Recommendation: Review before installation
```

### 4. 手动扫描
你也可以手动扫描特定技能：

```bash
# 扫描本地技能目录
clawguard scan /path/to/skill

# 扫描 GitHub 仓库
clawguard scan https://github.com/user/skill-repo

# 扫描已安装的技能
clawguard scan --installed
```

## 检测能力

### 提示注入检测
检测各种提示注入攻击模式：
- 角色扮演攻击
- 系统指令覆盖
- 输出注入
- 上下文污染
- 多轮对话劫持

### 隐私信息检测
自动识别和遮蔽：
- 电子邮件地址
- 电话号码
- 身份证号
- 信用卡号
- API 密钥
- 密码
- 个人身份信息

### 供应链安全
检测供应链攻击：
- 依赖项漏洞扫描
- 代码完整性检查
- 来源验证
- 数字签名验证

### 行为监控
安装后监控：
- 异常的 API 调用
- 敏感文件访问
- 网络连接异常
- 系统资源滥用

## 配置选项

### 扫描规则
```json
{
  "scan_rules": {
    "prompt_injection": {
      "enabled": true,
      "severity": "high",
      "patterns": ["role:", "ignore", "forget", "system"]
    },
    "privacy_leak": {
      "enabled": true,
      "severity": "critical",
      "patterns": ["email:", "phone:", "id:", "api_key"]
    },
    "supply_chain": {
      "enabled": true,
      "severity": "medium",
      "check_dependencies": true
    }
  }
}
```

### 白名单
```json
{
  "whitelist": [
    "github.com/official/openclaw-skills",
    "github.com/trusted-org/skills"
  ]
}
```

## 报告格式
扫描完成后，ClawGuard 生成详细的报告：

```json
{
  "skill_name": "example-skill",
  "version": "1.0.0",
  "security_score": 85,
  "scan_timestamp": "2024-01-15T10:30:00Z",
  "issues": [
    {
      "type": "prompt_injection",
      "severity": "high",
      "file": "prompts/main.txt",
      "line": 12,
      "description": "Potential system instruction override detected",
      "recommendation": "Review and sanitize the prompt"
    },
    {
      "type": "privacy_leak",
      "severity": "critical",
      "file": "config.py",
      "line": 8,
      "description": "API key found in source code",
      "recommendation": "Move to environment variables"
    }
  ],
  "summary": {
    "total_issues": 2,
    "high_severity": 1,
    "medium_severity": 0,
    "low_severity": 1,
    "passed_checks": 15
  }
}
```

## 安全最佳实践

### 1. 定期扫描
- 定期扫描已安装的技能
- 监控技能行为异常
- 保持检测规则更新

### 2. 权限最小化
- 只授予必要的权限
- 限制网络访问
- 沙盒执行环境

### 3. 代码审查
- 在安装前审查代码
- 检查依赖项安全
- 验证代码来源

### 4. 监控日志
- 记录所有扫描活动
- 监控异常行为
- 定期审查日志

## 故障排除

### 常见问题
1. **扫描误报**：调整检测规则阈值
2. **漏报**：更新检测模式库
3. **性能问题**：优化扫描配置

### 调试模式
```bash
# 启用详细日志
clawguard --verbose scan /path/to/skill

# 测试特定规则
clawguard --test-rule prompt_injection
```

## 相关链接
- [ClawGuard GitHub](https://github.com/NeuZhou/clawguard)
- [OpenClaw 文档](https://github.com/openclaw/openclaw)
- [AI 安全最佳实践](https://owasp.org/www-project-ai-security/)

---

**原文链接**：[English Version](https://github.com/AlexAnys/awesome-openclaw-usecases/blob/main/usecases/agent-security-scanning.md)