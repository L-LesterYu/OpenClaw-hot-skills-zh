#!/usr/bin/env node

import { Readability } from "@mozilla/readability";
import { JSDOM } from "jsdom";
import TurndownService from "turndown";
import { gfm } from "turndown-plugin-gfm";

const url = process.argv[2];

if (!url) {
	console.log("用法: content.js <网址>");
	console.log("\n从网页提取可读内容并转为 Markdown 格式。");
	console.log("\n示例:");
	console.log("  content.js https://example.com/article");
	console.log("  content.js https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html");
	process.exit(1);
}

function htmlToMarkdown(html) {
	const turndown = new TurndownService({ headingStyle: "atx", codeBlockStyle: "fenced" });
	turndown.use(gfm);
	turndown.addRule("removeEmptyLinks", {
		filter: (node) => node.nodeName === "A" && !node.textContent?.trim(),
		replacement: () => "",
	});
	return turndown
		.turndown(html)
		.replace(/\[\\?\[\s*\\?\]\]\([^)]*\)/g, "")
		.replace(/ +/g, " ")
		.replace(/\s+,/g, ",")
		.replace(/\s+\./g, ".")
		.replace(/\n{3,}/g, "\n\n")
		.trim();
}

try {
	const response = await fetch(url, {
		headers: {
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
		},
		signal: AbortSignal.timeout(15000),
	});
	
	if (!response.ok) {
		console.error(`HTTP ${response.status}: ${response.statusText}`);
		process.exit(1);
	}
	
	const html = await response.text();
	const dom = new JSDOM(html, { url });
	const reader = new Readability(dom.window.document);
	const article = reader.parse();
	
	if (article && article.content) {
		if (article.title) {
			console.log(`# ${article.title}\n`);
		}
		console.log(htmlToMarkdown(article.content));
		process.exit(0);
	}
	
	// Fallback: try to extract main content
	const fallbackDoc = new JSDOM(html, { url });
	const body = fallbackDoc.window.document;
	body.querySelectorAll("script, style, noscript, nav, header, footer, aside").forEach(el => el.remove());
	
	const title = body.querySelector("title")?.textContent?.trim();
	const main = body.querySelector("main, article, [role='main'], .content, #content") || body.body;
	
	if (title) {
		console.log(`# ${title}\n`);
	}
	
	const text = main?.innerHTML || "";
	if (text.trim().length > 100) {
		console.log(htmlToMarkdown(text));
	} else {
		console.error("无法从此页面提取可读内容。");
		process.exit(1);
	}
} catch (e) {
	console.error(`错误: ${e.message}`);
	process.exit(1);
}
