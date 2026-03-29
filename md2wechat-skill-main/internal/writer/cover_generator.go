// Package writer provides assisted writing functionality with customizable creator styles
package writer

import (
	"fmt"
	"regexp"
	"strings"
)

// CoverGenerator 封面生成器
type CoverGenerator struct {
	styleManager *StyleManager
}

// NewCoverGenerator 创建封面生成器
func NewCoverGenerator(styleManager *StyleManager) *CoverGenerator {
	return &CoverGenerator{
		styleManager: styleManager,
	}
}

// GeneratePrompt 生成封面提示词
func (cg *CoverGenerator) GeneratePrompt(req *GenerateCoverRequest) (*GenerateCoverResult, error) {
	// 获取风格
	style, err := cg.styleManager.GetStyle(req.StyleName)
	if err != nil {
		return &GenerateCoverResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// 构建文章内容
	content := req.ArticleContent
	if req.ArticleTitle != "" {
		content = fmt.Sprintf("标题：%s\n\n内容：%s", req.ArticleTitle, req.ArticleContent)
	}

	// 使用风格的封面提示词模板
	prompt := style.CoverPrompt

	// 替换占位符
	if strings.Contains(prompt, "{article_content}") {
		prompt = strings.ReplaceAll(prompt, "{article_content}", content)
	} else {
		prompt = prompt + "\n\n# 文章内容\n" + content
	}

	return &GenerateCoverResult{
		Prompt:   prompt,
		MetaData: cg.analyzeContent(req),
		Success:  true,
	}, nil
}

// analyzeContent 分析文章内容（简化版）
func (cg *CoverGenerator) analyzeContent(req *GenerateCoverRequest) CoverMetaData {
	content := req.ArticleContent
	if req.ArticleTitle != "" {
		content = req.ArticleTitle + " " + content
	}

	return CoverMetaData{
		CoreTheme:    cg.extractTheme(content),
		CoreView:     cg.extractView(content),
		Mood:         cg.determineMood(content),
		VisualAnchor: cg.findVisualAnchor(content),
	}
}

// extractTheme 提取核心主题（简化版）
func (cg *CoverGenerator) extractTheme(content string) string {
	// 移除标点符号（保留字母、数字、空格、中文字符）
	chineseRange := "[\\x{4e00}-\\x{9fff}]"
	re := regexp.MustCompile(`[^\w\s` + chineseRange + `]+`)
	content = re.ReplaceAllString(content, "")

	// 取前几个关键词
	words := strings.Fields(content)
	if len(words) > 0 {
		return words[0]
	}
	return "主题"
}

// extractView 提取核心观点（简化版）
func (cg *CoverGenerator) extractView(content string) string {
	// 查找可能的观点句
	patterns := []string{
		`我认为([^。！？\n]{5,30})`,
		`([^。！？\n]{5,30})，这是`,
		`([^。！？\n]{5,30})的本质`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// 取第一句话
	sentences := regexp.MustCompile(`[。！？\n]`).Split(content, 2)
	if len(sentences) > 0 && len(strings.TrimSpace(sentences[0])) > 5 {
		return strings.TrimSpace(sentences[0])
	}

	// 截取前50字
	if len(content) > 50 {
		return content[:50] + "..."
	}
	return content
}

// determineMood 确定情绪基调
func (cg *CoverGenerator) determineMood(content string) string {
	moods := map[string][]string{
		"inspirational": {"激励", "启发", "成长", "突破", "成功", "梦想"},
		"mysterious":    {"秘密", "隐藏", "未知", "探索", "谜团"},
		"protective":    {"保护", "防御", "安全", "避免", "防范"},
		"breakthrough":  {"突破", "改变", "转型", "升级", "革新"},
		"contemplative": {"思考", "反思", "观察", "理解", "洞察"},
		"rebellious":    {"反叛", "挑战", "质疑", "不循规蹈矩", "打破"},
	}

	content = strings.ToLower(content)

	for mood, keywords := range moods {
		for _, keyword := range keywords {
			if strings.Contains(content, keyword) {
				return mood
			}
		}
	}

	// 默认情绪
	return "contemplative"
}

// findVisualAnchor 找视觉锚点
func (cg *CoverGenerator) findVisualAnchor(content string) string {
	// 查找可能的具体物体 - 使用直接构建的正则表达式
	chineseRange := "[\\x{4e00}-\\x{9fff}]"

	patterns := []string{
		`([一个一二三四五六七八九十\d]+只?` + chineseRange + `{2,6})`, // 数量+物体
		chineseRange + `{2,6}人`, // 人物
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return "元素"
}

// ExplainMetaphor 解释隐喻
func (cg *CoverGenerator) ExplainMetaphor(meta CoverMetaData) string {
	return fmt.Sprintf("封面通过「%s」作为视觉锚点，表达「%s」的核心观点",
		meta.VisualAnchor, meta.CoreView)
}

// GetCoverStyleInfo 获取封面风格信息
func (cg *CoverGenerator) GetCoverStyleInfo(styleName string) (*CoverStyleInfo, error) {
	style, err := cg.styleManager.GetStyle(styleName)
	if err != nil {
		return nil, err
	}

	return &CoverStyleInfo{
		Style:       style.CoverStyle,
		Mood:        style.CoverMood,
		ColorScheme: style.CoverColorScheme,
		AspectRatio: "16:9",
		Orientation: "horizontal",
	}, nil
}

// CoverStyleInfo 封面风格信息
type CoverStyleInfo struct {
	Style       string
	Mood        string
	ColorScheme []string
	AspectRatio string
	Orientation string
}

// FormatCoverStyleInfo 格式化封面风格信息
func (info *CoverStyleInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("风格: %s\n", info.Style))
	sb.WriteString(fmt.Sprintf("情绪: %s\n", info.Mood))
	sb.WriteString(fmt.Sprintf("比例: %s (%s)\n", info.AspectRatio, info.Orientation))

	if len(info.ColorScheme) > 0 {
		sb.WriteString("配色: " + strings.Join(info.ColorScheme, ", "))
	}

	return sb.String()
}

// GenerateCoverPromptWithStyle 使用指定风格生成封面提示词
func (cg *CoverGenerator) GenerateCoverPromptWithStyle(style *WriterStyle, articleTitle, articleContent string) string {
	content := articleContent
	if articleTitle != "" {
		content = fmt.Sprintf("标题：%s\n\n%s", articleTitle, articleContent)
	}

	prompt := style.CoverPrompt
	if strings.Contains(prompt, "{article_content}") {
		prompt = strings.ReplaceAll(prompt, "{article_content}", content)
	} else {
		prompt = prompt + "\n\n# 文章内容\n" + content
	}

	return prompt
}

// ValidateCoverRequest 验证封面生成请求
func (cg *CoverGenerator) ValidateCoverRequest(req *GenerateCoverRequest) error {
	if req.ArticleContent == "" {
		return &WriterError{
			Code:    ErrCodeInvalidInput,
			Message: "文章内容不能为空",
		}
	}

	if req.StyleName == "" {
		req.StyleName = DefaultStyleName
	}

	return nil
}

// IsCoverRequest 检查结果是否是封面请求
func IsCoverRequest(result *GenerateCoverResult) bool {
	return result != nil && result.Error == "" && result.Success
}

// BuildDefaultCoverPrompt 构建默认封面提示词
func (cg *CoverGenerator) BuildDefaultCoverPrompt(articleTitle, articleContent string) string {
	result, err := cg.GeneratePrompt(&GenerateCoverRequest{
		ArticleTitle:   articleTitle,
		ArticleContent: articleContent,
		StyleName:      DefaultStyleName,
	})

	if err != nil {
		// 返回基础提示词
		return fmt.Sprintf("Generate a 16:9 horizontal article cover about: %s\n\n%s", articleTitle, articleContent)
	}

	return result.Prompt
}

// ExtractCoverRequest 提取封面请求（用于 AI 模式）
func ExtractCoverRequest(result *GenerateCoverResult) string {
	if result.Success && result.Prompt != "" {
		return result.Prompt
	}
	return ""
}

// CompleteCoverRequest 完成封面请求（图片生成后调用）
func CompleteCoverRequest(result *GenerateCoverResult, imageURL, mediaID string) *GenerateCoverResult {
	if result == nil {
		return &GenerateCoverResult{
			Success: false,
			Error:   "结果为空",
		}
	}

	result.ImageURL = imageURL
	result.MediaID = mediaID

	return result
}

// GetCoverPromptTemplate 获取封面提示词模板
func (cg *CoverGenerator) GetCoverPromptTemplate(styleName string) (string, error) {
	style, err := cg.styleManager.GetStyle(styleName)
	if err != nil {
		return "", err
	}

	return style.CoverPrompt, nil
}

// FormatCoverResult 格式化封面生成结果
func FormatCoverResult(result *GenerateCoverResult) string {
	var sb strings.Builder

	if result.Success {
		sb.WriteString("✅ 封面提示词已生成\n\n")

		if result.Explanation != "" {
			sb.WriteString("📖 隐喻说明: ")
			sb.WriteString(result.Explanation)
			sb.WriteString("\n\n")
		}

		sb.WriteString("🎨 提示词:\n")
		sb.WriteString(result.Prompt)

		if result.ImageURL != "" {
			sb.WriteString("\n\n📸 图片: ")
			sb.WriteString(result.ImageURL)
		}

		if result.MediaID != "" {
			sb.WriteString("\n\n📦 微信素材ID: ")
			sb.WriteString(result.MediaID)
		}
	} else {
		sb.WriteString("❌ 封面生成失败: ")
		sb.WriteString(result.Error)
	}

	return sb.String()
}
