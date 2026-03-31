# 能签单的 AI 员工

你的 AI 智能体可以处理调研、撰写提案、发送邮件——但当客户说"把合同发给我"……你还得手动介入。现在不用了。

这个用例将 OpenClaw 变成成交高手。它能以 Markdown 格式起草合同、保密协议（NDA）或提案，并发送电子签名——全程自动化。双方签署后，均会收到一份 SHA-256 认证的安全 PDF。就这么简单。

## 痛点

每一个 AI 销售工作流都存在同样的缺口：智能体什么都能做，*唯独*无法完成纸质签署。Stripe 处理支付，AgentMail 处理邮件，但签署呢？仍然需要手动操作——直到现在。

## 工作原理

1. 你告诉 OpenClaw 一笔交易、客户或协议的详情
2. 它以 Markdown 格式起草文档（NDA、合同、提案、SOW——任意类型）
3. 通过 Signbee 发送文档进行电子签名
4. 发送方签署 → 收件方收到邮件 → 收件方签署
5. 双方均收到带有 SHA-256 认证的防篡改签名 PDF

整个过程自动运行。你醒来时就能在收件箱里找到已签署的合同。

## 所需技能

[Signbee MCP Server](https://github.com/signbee/mcp) — 通过 npm 安装：

```json
{
  "mcpServers": {
    "signbee": {
      "command": "npx",
      "args": ["-y", "signbee-mcp"],
      "env": {
        "SIGNBEE_API_KEY": "your-api-key-from-signb.ee"
      }
    }
  }
}
```

在 [signb.ee](https://signb.ee) 获取免费 API 密钥（每月 5 份文档免费）。

## 设置

1. 安装 Signbee MCP Server（配置如上）
2. 指示 OpenClaw：

```txt
你是我的 AI 成交助手。当我告诉你一笔交易时，你需要：

1. 以规范的 Markdown 格式起草适当的文档（NDA、合同、提案或 SOW）
2. 包含所有相关条款、日期和当事方信息
3. 使用 send_document 工具发送进行电子签名
4. 汇报文档状态

起草时应专业且简洁。使用规范的标题、编号条款和清晰的语言。始终包含生效日期、当事方名称和终止条款。

如果我提供 PDF，请使用 send_document_pdf 并附上 URL。
```

## 示例提示词

**促成自由职业交易：**
> "给 acme.com 的 Sarah（sarah@acme.com）发送一份 NDA。我们即将开始咨询合作。我叫 Michael Beckett，邮箱 michael@company.com。"

**以 PDF 形式发送提案：**
> "我有一份提案 PDF，地址是 https://example.com/proposal.pdf ——请发送给 James（james@client.com）进行签署。"

**自主处理交易：**
> "为 DataCorp 起草一份为期 6 个月的服务协议，每月 $5,000（联系人：lisa@datacorp.com）。包含付款条款、知识产权归属和保密条款，然后发送签署。"

## 独特之处

大多数"AI 员工"解决方案止步于对话。这一个能真正签单。你的 AI 从线索 → 提案 → 签署合同完全自动化。结合邮件（AgentMail）、支付（Stripe）和日程安排（Cal.com），你将获得一台全自动的销售机器。

## 相关链接

- [Signbee](https://signb.ee) — 面向 AI 智能体的文档签署 API
- [signbee-mcp on npm](https://www.npmjs.com/package/signbee-mcp) — MCP 服务器包
- [GitHub](https://github.com/signbee/mcp) — 源代码
