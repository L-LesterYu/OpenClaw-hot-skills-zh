# AI Employee That Closes Deals

Your AI agent handles research, writes proposals, sends emails—but when the client says "send me the contract"... you still have to manually intervene. Not anymore.

This use case turns OpenClaw into a deal closer. It drafts contracts, NDAs, or proposals in Markdown and sends them for electronic signatures—all automatically. After both parties sign, both receive a SHA-256 certified secure PDF. That's it.

## The Pain Point

Every AI sales workflow has the same gap: the agent can do everything, *except* get signatures on paper. Stripe handles payments, AgentMail handles emails, but signing? That's still manual—until now.

## How It Works

1. You tell OpenClaw about a deal, client, or agreement
2. It drafts the document in Markdown (NDA, contract, proposal, SOW—anything)
3. Sends the document for electronic signature via Signbee
4. Sender signs → recipient gets email → recipient signs
5. Both parties receive tamper-proof signed PDF with SHA-256 certification

The whole process runs automatically. You wake up to find signed contracts in your inbox.

## Required Skills

[Signbee MCP Server](https://github.com/signbee/mcp) — Install via npm:

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

Get a free API key at [signb.ee](https://signb.ee) (5 documents per month free).

## Setup

1. Install Signbee MCP Server (config as above)
2. Instruct OpenClaw:

```txt
You are my AI deal closer. When I tell you about a deal, you need to:

1. Draft an appropriate document in clean Markdown (NDA, contract, proposal, or SOW)
2. Include all relevant terms, dates, and party information
3. Use the send_document tool to send for electronic signature
4. Report document status

When drafting, be professional but concise. Use proper headings, numbered clauses, and clear language. Always include effective dates, party names, and termination clauses.

If I give you a PDF, use send_document_pdf with the URL.
```

## Example Prompts

**Close a freelance deal:**
> "Send an NDA to Sarah at acme.com (sarah@acme.com). We're starting a consulting partnership. My name is Michael Beckett, email michael@company.com."

**Send proposal as PDF:**
> "I have a proposal PDF at https://example.com/proposal.pdf — please send it to James at james@client.com for signing."

**Handle deal autonomously:**
> "Draft a 6-month service agreement for DataCorp at $5,000/month (contact: lisa@datacorp.com). Include payment terms, IP ownership, and confidentiality clauses, then send for signing."

## What Makes It Different

Most "AI employee" solutions stop at conversation. This one closes deals. Your AI goes from lead → proposal → signed contract completely automatically. Combined with email (AgentMail), payments (Stripe), and scheduling (Cal.com), you get a fully automated sales machine.

## Related Links

- [Signbee](https://signb.ee) — Document signing API for AI agents
- [signbee-mcp on npm](https://www.npmjs.com/package/signbee-mcp) — MCP server package
- [GitHub](https://github.com/signbee/mcp) — Source code