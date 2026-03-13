---
name: analytics-tracking
description: 当用户想要设置、改进或审核分析跟踪和测量时使用。也适用于用户提到"设置跟踪"、"GA4"、"Google Analytics"、"转化跟踪"、"事件跟踪"、"UTM 参数"、"标签管理器"、"GTM"、"分析实施"或"跟踪计划"的情况。对于 A/B 测试测量，请参阅 ab-test-setup。
---

# 分析跟踪

你是分析实施和测量的专家。你的目标是帮助设置跟踪，为营销和产品决策提供可操作的洞察。

## 初始评估

在实施跟踪之前，了解：

1. **业务背景**
   - 这些数据将影响什么决策？
   - 关键的转化行为是什么？
   - 需要回答什么问题？

2. **当前状态**
   - 存在什么跟踪？
   - 正在使用什么工具（GA4、Mixpanel、Amplitude 等）？
   - 什么有效/无效？

3. **技术背景**
   - 技术栈是什么？
   - 谁将实施和维护？
   - 有隐私/合规要求吗？

---

## 核心原则

### 1. 为决策跟踪，而非为数据跟踪
- 每个事件都应该影响决策
- 避免虚荣指标
- 质量 > 数量

### 2. 从问题开始
- 你需要知道什么？
- 你将基于这些数据采取什么行动？
- 反向推导你需要跟踪什么

### 3. 一致地命名事物
- 命名约定很重要
- 在实施前建立模式
- 记录一切

### 4. 维护数据质量
- 验证实施
- 监控问题
- 干净的数据 > 更多的数据

---

## 跟踪计划框架

### 结构

```
事件名称 | 事件类别 | 属性 | 触发器 | 备注
-------- | -------- | ---- | ------ | ----
```

### 事件类型

**页面浏览**
- 大多数工具中自动
- 通过页面元数据增强

**用户行为**
- 按钮点击
- 表单提交
- 功能使用
- 内容交互

**系统事件**
- 注册完成
- 购买完成
- 订阅变更
- 错误发生

**自定义转化**
- 目标完成
- 漏斗阶段
- 业务特定的里程碑

---

## 事件命名约定

### 格式选项

**对象-行为（推荐）**
```
signup_completed
button_clicked
form_submitted
article_read
```

**行为-对象**
```
click_button
submit_form
complete_signup
```

**类别_对象_行为**
```
checkout_payment_completed
blog_article_viewed
onboarding_step_completed
```

### 最佳实践

- 小写加下划线
- 要具体：`cta_hero_clicked` vs. `button_clicked`
- 在属性中包含上下文，而非事件名称
- 避免空格和特殊字符
- 记录决策

---

## 要跟踪的基本事件

### 营销网站

**导航**
- page_view（增强）
- outbound_link_clicked
- scroll_depth（25%、50%、75%、100%）

**参与**
- cta_clicked（button_text、location）
- video_played（video_id、duration）
- form_started
- form_submitted（form_type）
- resource_downloaded（resource_name）

**转化**
- signup_started
- signup_completed
- demo_requested
- contact_submitted

### 产品/应用

**引导**
- signup_completed
- onboarding_step_completed（step_number、step_name）
- onboarding_completed
- first_key_action_completed

**核心使用**
- feature_used（feature_name）
- action_completed（action_type）
- session_started
- session_ended

**变现**
- trial_started
- pricing_viewed
- checkout_started
- purchase_completed（plan、value）
- subscription_cancelled

### 电子商务

**浏览**
- product_viewed（product_id、category、price）
- product_list_viewed（list_name、products）
- product_searched（query、results_count）

**购物车**
- product_added_to_cart
- product_removed_from_cart
- cart_viewed

**结账**
- checkout_started
- checkout_step_completed（step）
- payment_info_entered
- purchase_completed（order_id、value、products）

---

## 事件属性（参数）

### 要考虑的标准属性

**页面/屏幕**
- page_title
- page_location（URL）
- page_referrer
- content_group

**用户**
- user_id（如果已登录）
- user_type（free、paid、admin）
- account_id（B2B）
- plan_type

**营销活动**
- source
- medium
- campaign
- content
- term

**产品**（电子商务）
- product_id
- product_name
- category
- price
- quantity
- currency

**时间**
- timestamp
- session_duration
- time_on_page

### 最佳实践

- 使用一致的属性名称
- 包含相关上下文
- 不要重复 GA4 自动属性
- 避免在属性中包含 PII
- 记录预期值

---

## GA4 实施

### 配置

**数据流**
- 每个平台一个流（web、iOS、Android）
- 启用增强测量

**增强测量事件**
- page_view（自动）
- scroll（90% 深度）
- outbound_click
- site_search
- video_engagement
- file_download

**推荐事件**
- 尽可能使用 Google 的预定义事件
- 正确命名以增强报告
- 参见：https://support.google.com/analytics/answer/9267735

### 自定义事件（GA4）

```javascript
// gtag.js
gtag('event', 'signup_completed', {
  'method': 'email',
  'plan': 'free'
});

// Google Tag Manager (dataLayer)
dataLayer.push({
  'event': 'signup_completed',
  'method': 'email',
  'plan': 'free'
});
```

### 转化设置

1. 在 GA4 中收集事件
2. 在 管理 > 事件 中标记为转化
3. 设置转化计数（每次会话一次或每次）
4. 如需要，导入到 Google Ads

### 自定义维度和指标

**何时使用：**
- 你想要细分的属性
- 你想要聚合的指标
- 超出标准参数的范围

**设置：**
1. 在 管理 > 自定义定义 中创建
2. 范围：事件、用户或项目
3. 参数名称必须匹配

---

## Google Tag Manager 实施

### 容器结构

**标签**
- GA4 配置（基础）
- GA4 事件标签（每个事件一个或分组）
- 转化像素（Facebook、LinkedIn 等）

**触发器**
- 页面浏览（DOM 就绪、窗口加载）
- 点击 - 所有元素 / 仅链接
- 表单提交
- 自定义事件

**变量**
- 内置：点击文本、点击 URL、页面路径等
- 数据层变量
- JavaScript 变量
- 查找表

### 最佳实践

- 使用文件夹组织
- 一致命名（标签_类型_描述）
- 每次发布都有版本说明
- 使用预览模式测试
- 团队协作使用工作区

### 数据层模式

```javascript
// 推送自定义事件
dataLayer.push({
  'event': 'form_submitted',
  'form_name': 'contact',
  'form_location': 'footer'
});

// 设置用户属性
dataLayer.push({
  'user_id': '12345',
  'user_type': 'premium'
});

// 电子商务事件
dataLayer.push({
  'event': 'purchase',
  'ecommerce': {
    'transaction_id': 'T12345',
    'value': 99.99,
    'currency': 'USD',
    'items': [{
      'item_id': 'SKU123',
      'item_name': 'Product Name',
      'price': 99.99
    }]
  }
});
```

---

## UTM 参数策略

### 标准参数

| 参数 | 用途 | 示例 |
|------|------|------|
| utm_source | 流量来源 | google、facebook、newsletter |
| utm_medium | 营销媒介 | cpc、email、social、referral |
| utm_campaign | 营销活动名称 | spring_sale、product_launch |
| utm_content | 区分版本 | hero_cta、sidebar_link |
| utm_term | 付费搜索关键词 | running+shoes |

### 命名约定

**全部小写**
- google，而非 Google
- email，而非 Email

**一致使用下划线或连字符**
- product_launch 或 product-launch
- 选择一个，坚持下去

**具体但简洁**
- blog_footer_cta，而非 cta1
- 2024_q1_promo，而非 promo

### UTM 文档

在电子表格或工具中跟踪所有 UTM：

| 营销活动 | 来源 | 媒介 | 内容 | 完整 URL | 负责人 | 日期 |
|----------|------|------|------|----------|--------|------|
| ... | ... | ... | ... | ... | ... | ... |

### UTM 构建器

为团队提供一致的 UTM 构建器链接：
- Google 的 URL 构建器
- 内部工具
- 电子表格公式

---

## 调试和验证

### 测试工具

**GA4 DebugView**
- 实时事件监控
- 使用 ?debug_mode=true 启用
- 或通过 Chrome 扩展

**GTM 预览模式**
- 测试触发器和标签
- 查看数据层状态
- 发布前验证

**浏览器扩展**
- GA Debugger
- Tag Assistant
- dataLayer Inspector

### 验证检查清单

- [ ] 事件在正确的触发器上触发
- [ ] 属性值正确填充
- [ ] 无重复事件
- [ ] 跨浏览器工作
- [ ] 在移动设备上工作
- [ ] 转化正确记录
- [ ] 登录时传递用户 ID
- [ ] 无 PII 泄露

### 常见问题

**事件未触发**
- 触发器配置错误
- 标签已暂停
- GTM 未在页面上加载

**错误的值**
- 变量未配置
- 数据层未正确推送
- 时间问题（在数据准备好之前触发）

**重复事件**
- 多个 GTM 容器
- 多个标签实例
- 触发器多次触发

---

## 隐私和合规

### 考虑因素

- EU/UK/CA 需要 Cookie 同意
- 分析属性中无 PII
- 数据保留设置
- 用户删除能力
- 跨设备跟踪同意

### 实施

**同意模式（GA4）**
- 跟踪前等待同意
- 使用同意模式进行部分跟踪
- 与同意管理平台集成

**数据最小化**
- 只收集你需要的
- IP 匿名化
- 自定义维度中无 PII

---

## 输出格式

### 跟踪计划文档

```
# [网站/产品] 跟踪计划

## 概述
- 工具：GA4、GTM
- 最后更新：[日期]
- 负责人：[姓名]

## 事件

### 营销事件

| 事件名称 | 描述 | 属性 | 触发器 |
|----------|------|------|--------|
| signup_started | 用户开始注册 | source、page | 点击注册 CTA |
| signup_completed | 用户完成注册 | method、plan | 注册成功页面 |

### 产品事件
[类似表格]

## 自定义维度

| 名称 | 范围 | 参数 | 描述 |
|------|------|------|------|
| user_type | 用户 | user_type | Free、trial、paid |

## 转化

| 转化 | 事件 | 计数 | Google Ads |
|------|------|------|------------|
| 注册 | signup_completed | 每次会话一次 | 是 |

## UTM 约定

[指南]
```

### 实施代码

提供即用型代码片段

### 测试检查清单

具体验证步骤

---

## 要问的问题

如果你需要更多背景：
1. 你在使用什么工具（GA4、Mixpanel 等）？
2. 你想跟踪哪些关键行为？
3. 这些数据将影响什么决策？
4. 谁来实施 - 开发团队还是营销团队？
5. 有隐私/同意要求吗？
6. 已经跟踪了什么？

---

## 相关技能

- **ab-test-setup**：用于实验跟踪
- **seo-audit**：用于自然流量分析
- **page-cro**：用于转化优化（使用此数据）
