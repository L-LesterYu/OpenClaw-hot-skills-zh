# 错误日志

此文件记录命令失败、异常和集成问题。

## 使用方法

当发生以下情况时记录错误：
- 命令返回非零退出码
- 出现异常或堆栈跟踪
- 意外的输出或行为
- 超时或连接失败

## 格式

```markdown
## [ERR-YYYYMMDD-XXX] skill_or_command_name

**Logged**: ISO-8601 时间戳
**Priority**: high
**Status**: pending
**Area**: frontend | backend | infra | tests | docs | config

### Summary
失败内容的简要描述

### Error
```
实际错误消息或输出
```

### Context
- 尝试的命令/操作
- 使用的输入或参数
- 相关的环境详情

### Suggested Fix
如果可识别，可能解决此问题的方法

### Metadata
- Reproducible: yes | no | unknown
- Related Files: path/to/file.ext
- See Also: ERR-20250110-001 (如果重复)

---
```

## 示例条目

## [ERR-20250224-001] npm_install

**Logged**: 2025-02-24T09:35:00Z
**Priority**: high
**Status**: pending
**Area**: config

### Summary
npm install 失败 - 项目使用 pnpm

### Error
```
npm ERR! code ERESOLVE
npm ERR! ERESOLVE could not resolve
```

### Context
- 尝试运行 `npm install`
- 项目有 pnpm-lock.yaml 文件
- 使用 pnpm 工作区

### Suggested Fix
使用 `pnpm install` 代替 `npm install`

### Metadata
- Reproducible: yes
- Related Files: package.json, pnpm-lock.yaml
- See Also: 

---
