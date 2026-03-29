// Package humanizer provides AI writing trace removal functionality
package humanizer

import (
	"fmt"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/promptcatalog"
)

// BuildPrompt 构建给 Claude 的提示词
func BuildPrompt(req *HumanizeRequest) string {
	if req.Intensity == "" {
		req.Intensity = IntensityMedium
	}

	catalog, err := promptcatalog.DefaultCatalog()
	if err != nil {
		return buildPromptFallback(req)
	}

	intensitySpec, err := catalog.Get("humanizer", req.Intensity.String())
	if err != nil {
		intensitySpec, err = catalog.Get("humanizer", IntensityMedium.String())
		if err != nil {
			return buildPromptFallback(req)
		}
	}

	rendered, _, err := catalog.Render("humanizer", "base", map[string]string{
		"INTENSITY_TEXT":            intensitySpec.Template,
		"STYLE_BLOCK":               buildStyleBlock(req),
		"FOCUS_BLOCK":               buildFocusBlock(req),
		"SOURCE_BLOCK":              buildSourceBlock(req),
		"OUTPUT_FORMAT_BLOCK":       getOutputFormatTemplate(),
		"OUTPUT_REQUIREMENTS_BLOCK": buildOutputRequirements(req),
		"CONTENT":                   req.Content,
	})
	if err != nil {
		return buildPromptFallback(req)
	}

	return rendered
}

// getOutputFormatTemplate 输出格式模板
func getOutputFormatTemplate() string {
	return `

## 输出格式

请按以下格式输出：

# 人性化后的文本

[这里是重写后的文本，去除 AI 痕迹，注入人味]

# 修改说明

[简要说明做了哪些主要修改，如果需要详细对比请说明]

# 质量评分

| 维度 | 得分 | 说明 |
|------|------|------|
| 直接性 | x/10 | [说明] |
| 节奏 | x/10 | [说明] |
| 信任度 | x/10 | [说明] |
| 真实性 | x/10 | [说明] |
| 精炼度 | x/10 | [说明] |
| **总分** | **x/50** | [评级] |

## 最终要求

1. 输出完整的重写后文本
`
}

func buildPromptFallback(req *HumanizeRequest) string {
	var prompt strings.Builder
	prompt.WriteString("# Humanizer-zh: 去除 AI 写作痕迹")
	prompt.WriteString("\n\n## 处理强度\n\n")
	prompt.WriteString(req.Intensity.Description())
	prompt.WriteString(buildStyleBlock(req))
	prompt.WriteString(buildFocusBlock(req))
	prompt.WriteString(buildSourceBlock(req))
	prompt.WriteString(getOutputFormatTemplate())
	prompt.WriteString(buildOutputRequirements(req))
	prompt.WriteString(fmt.Sprintf("\n\n# 待处理文本\n\n%s", req.Content))
	return prompt.String()
}

func buildStyleBlock(req *HumanizeRequest) string {
	if !req.PreserveStyle || req.OriginalStyle == "" {
		return ""
	}
	return fmt.Sprintf("\n\n## 风格保护\n\n原文采用「%s」写作风格，请保留该风格的核心特征。\n\n**重要原则**：\n- 如果某种模式是该风格刻意为之（如使用破折号制造停顿），请保留\n- 只去除无意的 AI 痕迹\n- 保持风格的一致性\n", req.OriginalStyle)
}

func buildFocusBlock(req *HumanizeRequest) string {
	if len(req.FocusOn) == 0 {
		return ""
	}

	var prompt strings.Builder
	prompt.WriteString("\n\n## 重点处理模式\n\n")
	prompt.WriteString("请重点关注以下类型的模式：\n")
	for _, p := range req.FocusOn {
		switch p {
		case PatternContent:
			prompt.WriteString("- **内容模式**：过度强调、夸大意义、宣传语言、模糊归因\n")
		case PatternLanguage:
			prompt.WriteString("- **语言语法**：AI 词汇、否定排比、三段式、同义词循环\n")
		case PatternStyle:
			prompt.WriteString("- **风格模式**：破折号过度、粗体滥用、表情符号\n")
		case PatternFiller:
			prompt.WriteString("- **填充词回避**：填充短语、过度限定、通用结论\n")
		case PatternCollaboration:
			prompt.WriteString("- **协作痕迹**：对话式填充、知识截止免责声明\n")
		}
	}
	return prompt.String()
}

func buildSourceBlock(req *HumanizeRequest) string {
	if req.SourceHint == "" {
		return ""
	}
	return fmt.Sprintf("\n\n## 源信息\n\n文本来源: %s\n", req.SourceHint)
}

func buildOutputRequirements(req *HumanizeRequest) string {
	var prompt strings.Builder
	if req.ShowChanges {
		prompt.WriteString("2. 提供修改说明和主要变更点\n")
	}
	if req.IncludeScore {
		prompt.WriteString("3. 按 5 维度给出质量评分\n")
	}
	prompt.WriteString("4. 只返回上述格式的内容，不需要其他解释\n")
	return prompt.String()
}

// BuildAIRequest 构建 AI 转换请求（用于与 writer 模块兼容）
func BuildAIRequest(req *HumanizeRequest) *AIConvertRequest {
	return &AIConvertRequest{
		Content: req.Content,
		Prompt:  BuildPrompt(req),
		Settings: HumanizeSettings{
			Intensity:     req.Intensity,
			FocusOn:       req.FocusOn,
			PreserveStyle: req.PreserveStyle,
			OriginalStyle: req.OriginalStyle,
			ShowChanges:   req.ShowChanges,
			IncludeScore:  req.IncludeScore,
		},
	}
}

// ParseIntensity 从字符串解析强度
func ParseIntensity(s string) HumanizeIntensity {
	switch strings.ToLower(s) {
	case "gentle", "light", "温和", "轻度":
		return IntensityGentle
	case "aggressive", "heavy", "激进", "深度":
		return IntensityAggressive
	case "medium", "normal", "中等", "标准", "":
		return IntensityMedium
	default:
		return IntensityMedium
	}
}

// ParseFocusPattern 从字符串切片解析聚焦模式
func ParseFocusPattern(patterns []string) []FocusPattern {
	var result []FocusPattern
	patternMap := map[string]FocusPattern{
		"content":       PatternContent,
		"language":      PatternLanguage,
		"style":         PatternStyle,
		"filler":        PatternFiller,
		"collaboration": PatternCollaboration,
		// 中文别名
		"内容": PatternContent,
		"语言": PatternLanguage,
		"风格": PatternStyle,
		"填充": PatternFiller,
		"协作": PatternCollaboration,
	}

	for _, p := range patterns {
		lower := strings.ToLower(strings.TrimSpace(p))
		if fp, ok := patternMap[lower]; ok {
			result = append(result, fp)
		}
	}
	return result
}
