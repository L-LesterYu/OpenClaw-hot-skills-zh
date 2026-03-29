# 通过聊天完成 AI 视频编辑

剪辑视频通常意味着打开时间线编辑器、拖拽素材、点击菜单。对于重复性操作——修剪片头、给一批视频添加字幕、调整 10 个视频的色调——这种手动循环会耗费数小时。

这个用例将视频编辑转化为对话。描述你想要的修改，放下文件，就能拿到成品。没有时间线，没有 GUI。

## 你可以做什么

- 用自然语言描述时间戳来修剪、裁切和合并片段
- 添加背景音乐并自动进行音频闪避
- 从语音生成并烧录字幕（支持 50+ 种语言）
- 色彩调色（"调暖一些"、"匹配第一个片段的风格"）
- 裁剪为竖屏，适配 TikTok/Reels/Shorts
- 批量处理多个文件，应用相同的编辑操作

## 所需技能

- [video-editor-ai](https://clawhub.ai/skills/video-editor-ai) ——基于聊天的视频编辑，支持背景音乐、字幕和导出
- [ai-subtitle-generator](https://clawhub.ai/skills/ai-subtitle-generator) ——自动字幕生成、字幕烧录、SRT 导出

## 设置方法

1. 安装技能：
```bash
clawhub install video-editor-ai
clawhub install ai-subtitle-generator
```

2. 将视频文件拖入聊天并描述你的编辑需求：
```text
将这个视频从 0:15 裁剪到 1:30，添加背景音乐（欢快的风格），并烧录英文字幕。
```

3. 批量处理时，描述编辑模式：
```text
我在 /videos/raw/ 下有 5 个视频片段。对每个视频执行以下操作：
- 裁剪为 9:16 竖屏
- 在底部添加自动生成的字幕
- 导出为 mp4
```

智能体会处理 API 调用、轮询完成状态，并将成品文件发送回聊天。

## 小贴士

- 明确指定时间戳和输出格式（如"导出为 1080p 的 mp4"）
- 字幕工作如源语言不是英语，请说明源语言
- 色彩调色使用描述性参考效果最好（如"温暖的日落色调"），而非技术数值
- 避免上传敏感画面（人脸、身份证件、机密屏幕），除非你已审查提供商的数据保留和隐私政策

---

**原文链接**：[English Version](https://github.com/hesamsheikh/awesome-openclaw-usecases/blob/main/usecases/ai-video-editing.md)
