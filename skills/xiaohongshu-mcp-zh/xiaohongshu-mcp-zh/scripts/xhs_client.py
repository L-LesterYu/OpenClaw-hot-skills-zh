#!/usr/bin/env python3
"""
小红书 MCP 客户端 - 用于 xiaohongshu-mcp HTTP API 的 Python 客户端。

用法：
    python xhs_client.py <命令> [选项]

命令：
    status              检查登录状态
    search <关键词>     按关键词搜索笔记
    detail <feed_id> <xsec_token>   获取笔记详情
    feeds               获取推荐信息流列表
    publish <标题> <内容> <图片>  发布笔记

示例：
    python xhs_client.py status
    python xhs_client.py search "咖啡推荐"
    python xhs_client.py detail "abc123" "token456"
    python xhs_client.py feeds
"""

import argparse
import json
import sys
import requests

BASE_URL = "http://localhost:18060"
TIMEOUT = 60


def check_status():
    """检查登录状态。"""
    try:
        resp = requests.get(f"{BASE_URL}/api/v1/login/status", timeout=TIMEOUT)
        data = resp.json()
        if data.get("success"):
            login_info = data.get("data", {})
            if login_info.get("is_logged_in"):
                print(f"✅ 已登录为：{login_info.get('username', '未知')}")
            else:
                print("❌ 未登录，请先运行登录工具。")
        else:
            print(f"❌ 错误：{data.get('error', '未知错误')}")
        return data
    except requests.exceptions.ConnectionError:
        print("❌ 无法连接到 MCP 服务器，请确认 xiaohongshu-mcp 正在 localhost:18060 上运行")
        sys.exit(1)


def search_notes(keyword, sort_by="综合", note_type="不限", publish_time="不限"):
    """按关键词搜索笔记，支持可选筛选条件。"""
    try:
        payload = {
            "keyword": keyword,
            "filters": {
                "sort_by": sort_by,
                "note_type": note_type,
                "publish_time": publish_time
            }
        }
        resp = requests.post(
            f"{BASE_URL}/api/v1/feeds/search",
            json=payload,
            timeout=TIMEOUT
        )
        data = resp.json()
        
        if data.get("success"):
            feeds = data.get("data", {}).get("feeds", [])
            print(f"🔍 找到 {len(feeds)} 条关于「{keyword}」的笔记：\n")
            
            for i, feed in enumerate(feeds, 1):
                note_card = feed.get("noteCard", {})
                user = note_card.get("user", {})
                interact = note_card.get("interactInfo", {})
                
                print(f"[{i}] {note_card.get('displayTitle', '无标题')}")
                print(f"    作者：{user.get('nickname', '未知')}")
                print(f"    点赞：{interact.get('likedCount', '0')} | 收藏：{interact.get('collectedCount', '0')} | 评论：{interact.get('commentCount', '0')}")
                print(f"    feed_id：{feed.get('id')}")
                print(f"    xsec_token：{feed.get('xsecToken')}")
                print()
        else:
            print(f"❌ 搜索失败：{data.get('error', '未知错误')}")
        
        return data
    except requests.exceptions.ConnectionError:
        print("❌ 无法连接到 MCP 服务器。")
        sys.exit(1)


def get_note_detail(feed_id, xsec_token, load_comments=False):
    """获取指定笔记的详细信息。"""
    try:
        payload = {
            "feed_id": feed_id,
            "xsec_token": xsec_token,
            "load_all_comments": load_comments
        }
        resp = requests.post(
            f"{BASE_URL}/api/v1/feeds/detail",
            json=payload,
            timeout=TIMEOUT
        )
        data = resp.json()
        
        if data.get("success"):
            note_data = data.get("data", {}).get("data", {})
            note = note_data.get("note", {})
            comments = note_data.get("comments", {})
            
            print(f"📝 笔记详情：\n")
            print(f"标题：{note.get('title', '无标题')}")
            print(f"作者：{note.get('user', {}).get('nickname', '未知')}")
            print(f"地区：{note.get('ipLocation', '未知')}")
            print(f"\n内容：\n{note.get('desc', '暂无内容')}\n")
            
            interact = note.get("interactInfo", {})
            print(f"点赞：{interact.get('likedCount', '0')} | 收藏：{interact.get('collectedCount', '0')} | 评论：{interact.get('commentCount', '0')}")
            
            comment_list = comments.get("list", [])
            if comment_list:
                print(f"\n💬 热门评论（共 {len(comment_list)} 条）：")
                for c in comment_list[:5]:
                    user_info = c.get("userInfo", {})
                    print(f"  - {user_info.get('nickname', '匿名用户')}：{c.get('content', '')}")
        else:
            print(f"❌ 获取详情失败：{data.get('error', '未知错误')}")
        
        return data
    except requests.exceptions.ConnectionError:
        print("❌ 无法连接到 MCP 服务器。")
        sys.exit(1)


def get_feeds():
    """获取推荐信息流列表。"""
    try:
        resp = requests.get(f"{BASE_URL}/api/v1/feeds/list", timeout=TIMEOUT)
        data = resp.json()
        
        if data.get("success"):
            feeds = data.get("data", {}).get("feeds", [])
            print(f"📋 推荐信息流（共 {len(feeds)} 条笔记）：\n")
            
            for i, feed in enumerate(feeds, 1):
                note_card = feed.get("noteCard", {})
                user = note_card.get("user", {})
                interact = note_card.get("interactInfo", {})
                
                print(f"[{i}] {note_card.get('displayTitle', '无标题')}")
                print(f"    作者：{user.get('nickname', '未知')}")
                print(f"    点赞：{interact.get('likedCount', '0')}")
                print()
        else:
            print(f"❌ 获取信息流失败：{data.get('error', '未知错误')}")
        
        return data
    except requests.exceptions.ConnectionError:
        print("❌ 无法连接到 MCP 服务器。")
        sys.exit(1)


def publish_note(title, content, images, tags=None):
    """发布新笔记。"""
    try:
        payload = {
            "title": title,
            "content": content,
            "images": images if isinstance(images, list) else [images]
        }
        if tags:
            payload["tags"] = tags if isinstance(tags, list) else [tags]
        
        resp = requests.post(
            f"{BASE_URL}/api/v1/publish",
            json=payload,
            timeout=120
        )
        data = resp.json()
        
        if data.get("success"):
            print(f"✅ 笔记发布成功！")
            print(f"   帖子 ID：{data.get('data', {}).get('post_id', '未知')}")
        else:
            print(f"❌ 发布失败：{data.get('error', '未知错误')}")
        
        return data
    except requests.exceptions.ConnectionError:
        print("❌ 无法连接到 MCP 服务器。")
        sys.exit(1)


def main():
    parser = argparse.ArgumentParser(
        description="小红书 MCP 客户端",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    subparsers = parser.add_subparsers(dest="command", help="可用命令")
    
    # status 命令
    subparsers.add_parser("status", help="检查登录状态")
    
    # search 命令
    search_parser = subparsers.add_parser("search", help="搜索笔记")
    search_parser.add_argument("keyword", help="搜索关键词")
    search_parser.add_argument("--sort", default="综合", 
                               choices=["综合", "最新", "最多点赞", "最多评论", "最多收藏"],
                               help="排序方式")
    search_parser.add_argument("--type", default="不限",
                               choices=["不限", "视频", "图文"],
                               help="笔记类型")
    search_parser.add_argument("--time", default="不限",
                               choices=["不限", "一天内", "一周内", "半年内"],
                               help="发布时间")
    search_parser.add_argument("--json", action="store_true", help="输出原始 JSON 数据")
    
    # detail 命令
    detail_parser = subparsers.add_parser("detail", help="获取笔记详情")
    detail_parser.add_argument("feed_id", help="笔记 ID")
    detail_parser.add_argument("xsec_token", help="安全令牌")
    detail_parser.add_argument("--comments", action="store_true", help="加载所有评论")
    detail_parser.add_argument("--json", action="store_true", help="输出原始 JSON 数据")
    
    # feeds 命令
    feeds_parser = subparsers.add_parser("feeds", help="获取推荐信息流")
    feeds_parser.add_argument("--json", action="store_true", help="输出原始 JSON 数据")
    
    # publish 命令
    publish_parser = subparsers.add_parser("publish", help="发布笔记")
    publish_parser.add_argument("title", help="笔记标题")
    publish_parser.add_argument("content", help="笔记内容")
    publish_parser.add_argument("images", help="图片链接（逗号分隔）")
    publish_parser.add_argument("--tags", help="标签（逗号分隔）")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        sys.exit(1)
    
    if args.command == "status":
        result = check_status()
    elif args.command == "search":
        result = search_notes(args.keyword, args.sort, args.type, args.time)
        if hasattr(args, 'json') and args.json:
            print(json.dumps(result, ensure_ascii=False, indent=2))
    elif args.command == "detail":
        result = get_note_detail(args.feed_id, args.xsec_token, args.comments)
        if hasattr(args, 'json') and args.json:
            print(json.dumps(result, ensure_ascii=False, indent=2))
    elif args.command == "feeds":
        result = get_feeds()
        if hasattr(args, 'json') and args.json:
            print(json.dumps(result, ensure_ascii=False, indent=2))
    elif args.command == "publish":
        images = args.images.split(",")
        tags = args.tags.split(",") if args.tags else None
        result = publish_note(args.title, args.content, images, tags)


if __name__ == "__main__":
    main()
