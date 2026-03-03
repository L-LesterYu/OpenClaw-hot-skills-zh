---
name: youtube-watcher
description: 获取并读取 YouTube 视频的字幕文本。当您需要总结视频、回答关于视频内容的问题或从中提取信息时使用。
author: michael gathara
version: 1.0.0
triggers:
  - "观看 youtube"
  - "总结视频"
  - "视频字幕"
  - "youtube 摘要"
  - "分析视频"
metadata: {"bot":{"emoji":"📺","requires":{"bins":["yt-dlp"]},"install":[{"id":"brew","kind":"brew","formula":"yt-dlp","bins":["yt-dlp"],"label":"Install yt-dlp (brew)"},{"id":"pip","kind":"pip","package":"yt-dlp","bins":["yt-dlp"],"label":"Install yt-dlp (pip)"}]}}
---

# YouTube 观看器

从 YouTube 视频获取字幕文本，以实现摘要、问答和内容提取功能。

## 使用方法

### 获取字幕

检索视频的文字字幕。

```bash
python3 {baseDir}/scripts/get_transcript.py "https://www.youtube.com/watch?v=VIDEO_ID"
```

## 使用示例

**总结视频：**

1. 获取字幕：
   ```bash
   python3 {baseDir}/scripts/get_transcript.py "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
   ```
2. 阅读输出并为用户总结。

**查找特定信息：**

1. 获取字幕。
2. 在文本中搜索关键词或根据内容回答用户的问题。

## 注意事项

- 需要安装 `yt-dlp` 并在 PATH 中可用。
- 适用于有隐藏字幕（CC）或自动生成字幕的视频。
- 如果视频没有字幕，脚本将失败并显示错误消息。
