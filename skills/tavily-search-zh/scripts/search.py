#!/usr/bin/env python3
"""
OpenClaw Tavily 网页搜索脚本
使用 Tavily API 执行网页搜索
"""

import sys
import json
import urllib.request
import urllib.parse
import re
import os

def get_api_key():
    """从 OpenClaw 配置文件中提取 Tavily API 密钥"""
    try:
        config_path = os.path.expanduser("~/.openclaw/openclaw.json")
        with open(config_path, 'r', encoding='utf-8') as f:
            config = json.load(f)
            # 首先尝试 skills.entries.web-search.apiKey
            api_key = config.get('skills', {}).get('entries', {}).get('web-search', {}).get('apiKey', '')
            if api_key:
                return api_key
            # 备用方案：尝试 mcpServers.tavily.url
            tavily_url = config.get('mcpServers', {}).get('tavily', {}).get('url', '')
            match = re.search(r'tavilyApiKey=([^&]+)', tavily_url)
            if match:
                return match.group(1)
    except Exception as e:
        print(f"读取配置文件出错: {e}", file=sys.stderr)
    return None

def search(query, num_results=5):
    """执行 Tavily 搜索"""
    api_key = get_api_key()
    if not api_key:
        return {"错误": "未找到 Tavily API 密钥。请在 ~/.openclaw/openclaw.json 的 skills.entries.web-search.apiKey 中添加密钥"}
    
    url = "https://api.tavily.com/search"
    headers = {"Content-Type": "application/json"}
    data = {
        "api_key": api_key,
        "query": query,
        "search_depth": "basic",
        "include_answer": False,
        "max_results": num_results
    }
    
    try:
        req = urllib.request.Request(
            url,
            data=json.dumps(data).encode('utf-8'),
            headers=headers,
            method='POST'
        )
        with urllib.request.urlopen(req, timeout=30) as response:
            result = json.loads(response.read().decode('utf-8'))
            return result
    except Exception as e:
        return {"错误": f"搜索失败: {str(e)}"}

def format_results(results):
    """格式化搜索结果为中文友好格式"""
    if "错误" in results:
        return f"❌ {results['错误']}"
    
    if "results" not in results:
        return "未找到搜索结果"
    
    output = []
    output.append(f"🔍 找到 {len(results['results'])} 条结果：\n")
    
    for i, item in enumerate(results['results'], 1):
        output.append(f"【{i}】{item.get('title', '无标题')}")
        output.append(f"    链接: {item.get('url', '')}")
        output.append(f"    摘要: {item.get('content', '无内容')[:200]}...")
        output.append("")
    
    return "\n".join(output)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("用法: python3 search.py <搜索查询>", file=sys.stderr)
        print("示例: python3 search.py 今天的热点新闻", file=sys.stderr)
        sys.exit(1)
    
    query = " ".join(sys.argv[1:])
    print(f"🔍 正在搜索: {query}\n")
    
    results = search(query)
    
    # 输出格式化的结果
    print(format_results(results))
    
    # 同时输出原始 JSON（便于程序处理）
    # print("\n--- 原始 JSON ---")
    # print(json.dumps(results, indent=2, ensure_ascii=False))
