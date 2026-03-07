# 国内搜索引擎深度搜索指南

## 🔍 百度深度搜索

### 1.1 基础高级搜索操作符

| 操作符 | 功能 | 示例 |
|--------|------|------|
| `""` | 精确匹配 | `"人工智能"` |
| `-` | 排除关键词 | `python -培训` |
| `|` | 或运算 | `机器学习|深度学习` |
| `site:` | 站内搜索 | `site:csdn.net python` |
| `filetype:` | 文件类型 | `filetype:pdf 技术报告` |
| `intitle:` | 标题包含 | `intitle:教程` |
| `inurl:` | URL包含 | `inurl:blog` |

### 1.2 百度搜索示例

```javascript
// 搜索CSDN上的Python教程
web_fetch({"url": "https://www.baidu.com/s?wd=site:csdn.net+python%E6%95%99%E7%A8%8B"})

// 搜索PDF格式技术文档
web_fetch({"url": "https://www.baidu.com/s?wd=filetype:pdf+%E6%9C%BA%E5%99%A8%E5%AD%A6%E4%B9%A0"})

// 精确匹配搜索
web_fetch({"url": "https://www.baidu.com/s?wd=%22%E6%B7%B1%E5%BA%A6%E5%AD%A6%E4%B9%A0%22"})
```

---

## 🔎 必应中国深度搜索

### 2.1 必应特色功能

| 功能 | 说明 | 示例 |
|------|------|------|
| 中文优先 | 中文内容优先展示 | 默认 |
| 国际版切换 | `ensearch=1` | 获取国际结果 |
| 图片搜索 | 高质量图片 | 支持筛选 |
| 学术搜索 | 必应学术 | 学术论文 |

### 2.2 必应搜索示例

```javascript
// 必应中文搜索
web_fetch({"url": "https://cn.bing.com/search?q=%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD&ensearch=0"})

// 必应国际搜索
web_fetch({"url": "https://cn.bing.com/search?q=artificial+intelligence&ensearch=1"})

// 站内搜索
web_fetch({"url": "https://cn.bing.com/search?q=site:github.com+react"})
```

---

## 📱 社交媒体搜索

### 3.1 微信公众号搜索

```javascript
// 搜索微信公众号文章
web_fetch({"url": "https://wx.sogou.com/weixin?type=2&query=%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD"})
```

### 3.2 头条搜索

```javascript
// 搜索今日头条内容
web_fetch({"url": "https://so.toutiao.com/search?keyword=%E7%A7%91%E6%8A%80%E6%96%B0%E9%97%BB"})
```

### 3.3 搜狗搜索

```javascript
// 搜狗综合搜索
web_fetch({"url": "https://sogou.com/web?query=%E6%9C%BA%E5%99%A8%E5%AD%A6%E4%B9%A0"})
```

---

## 📊 专业领域搜索

### 4.1 集思录（投资理财）

```javascript
// 搜索投资相关内容
web_fetch({"url": "https://www.jisilu.cn/explore/?keyword=%E5%8F%AF%E8%BD%AC%E5%80%BA"})
```

### 4.2 360搜索

```javascript
// 360综合搜索
web_fetch({"url": "https://www.so.com/s?q=%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD"})
```

---

## 🛠️ 搜索技巧汇总

### URL编码

```javascript
// 中文关键词需要URL编码
function encodeKeyword(keyword) {
  return encodeURIComponent(keyword);
}

// 示例
const keyword = "人工智能";
const encoded = encodeKeyword(keyword); // "%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD"
```

### 多引擎对比搜索

```javascript
// 同一关键词多引擎搜索
const keyword = "机器学习教程";
const encoded = encodeURIComponent(keyword);

const searches = {
  baidu: `https://www.baidu.com/s?wd=${encoded}`,
  bing_cn: `https://cn.bing.com/search?q=${encoded}&ensearch=0`,
  sogou: `https://sogou.com/web?query=${encoded}`,
  so360: `https://www.so.com/s?q=${encoded}`
};
```

---

## 📚 参考资料

- [百度搜索帮助](https://help.baidu.com/question)
- [必应搜索高级语法](https://support.microsoft.com/zh-cn/search)
