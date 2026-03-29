// Package writer provides assisted writing functionality with customizable creator styles
package writer

import (
	"fmt"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
)

// Generator 文章生成器接口
type Generator interface {
	// Generate 生成文章
	Generate(req *GenerateRequest) *GenerateResult
	// GenerateTitles 生成标题
	GenerateTitles(style *WriterStyle, content string, count int) []string
	// ExtractQuotes 提取金句
	ExtractQuotes(article string, style *WriterStyle) []string
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	Style       *WriterStyle
	UserInput   string
	InputType   InputType
	Title       string
	Length      Length
	ArticleType ArticleType
}

// GenerateResult 生成结果
type GenerateResult struct {
	Article   string
	Title     string
	Quotes    []string
	Status    action.Status `json:"status,omitempty"`
	Action    string        `json:"action,omitempty"`
	Retryable bool          `json:"retryable,omitempty"`
	Success   bool
	Error     string
	Prompt    string       // 用于调试，显示使用的提示词
	Style     *WriterStyle // 使用的风格（用于完成请求时提取金句）
}

// articleGenerator 文章生成器实现
type articleGenerator struct {
	// 可以在这里添加依赖，如 AI 客户端
}

// NewGenerator 创建文章生成器
func NewGenerator() Generator {
	return &articleGenerator{}
}

// Generate 生成文章
func (g *articleGenerator) Generate(req *GenerateRequest) *GenerateResult {
	if req.Style == nil {
		return &GenerateResult{
			Status:    action.StatusFailed,
			Action:    action.ActionWrite,
			Retryable: false,
			Success:   false,
			Error:     "风格未指定",
		}
	}

	if req.UserInput == "" {
		return &GenerateResult{
			Status:    action.StatusFailed,
			Action:    action.ActionWrite,
			Retryable: false,
			Success:   false,
			Error:     "用户输入不能为空",
		}
	}

	// 构建完整的 AI 提示词
	prompt := g.buildPrompt(req)

	// 构建结果
	result := &GenerateResult{
		Prompt:    prompt,
		Status:    action.StatusActionRequired,
		Action:    action.ActionWrite,
		Retryable: false,
		Success:   true,
		Style:     req.Style,
	}

	// 注意：实际的 AI 调用在外部（Claude）完成
	// 这里返回特殊标记，告诉调用者需要使用 AI
	result.Article = ""
	result.Error = "AI_MODE_REQUEST:" + prompt

	return result
}

// buildPrompt 构建完整的 AI 提示词
func (g *articleGenerator) buildPrompt(req *GenerateRequest) string {
	style := req.Style
	var prompt strings.Builder

	// 添加写作提示词
	prompt.WriteString(style.WritingPrompt)
	prompt.WriteString("\n\n")

	// 添加核心信念（如果有）
	if len(style.CoreBeliefs) > 0 {
		prompt.WriteString("## 核心写作 DNA\n")
		for i, belief := range style.CoreBeliefs {
			prompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, belief))
		}
		prompt.WriteString("\n")
	}

	// 添加输入类型说明
	prompt.WriteString("## 用户输入\n")
	prompt.WriteString(fmt.Sprintf("输入类型: %s\n", req.InputType.String()))

	if req.Title != "" {
		prompt.WriteString(fmt.Sprintf("标题: %s\n", req.Title))
	}

	prompt.WriteString(fmt.Sprintf("文章类型: %s\n", req.ArticleType.String()))
	prompt.WriteString(fmt.Sprintf("期望长度: %s\n", req.Length.String()))

	prompt.WriteString("\n## 用户内容\n")
	prompt.WriteString(req.UserInput)

	prompt.WriteString("\n\n---\n\n")
	prompt.WriteString("请根据以上要求，生成符合该风格的文章。")
	prompt.WriteString("直接输出文章内容，不需要其他说明。")

	return prompt.String()
}

// GenerateTitles 生成标题
func (g *articleGenerator) GenerateTitles(style *WriterStyle, content string, count int) []string {
	titles := []string{}

	// 使用风格中定义的标题公式
	if len(style.TitleFormulas) > 0 {
		for i, formula := range style.TitleFormulas {
			if len(titles) >= count {
				break
			}

			// 如果有示例，使用示例
			if len(formula.Examples) > 0 {
				for _, ex := range formula.Examples {
					if len(titles) >= count {
						break
					}
					titles = append(titles, ex)
				}
			} else if formula.Template != "" {
				// 使用模板生成（简化版）
				titles = append(titles, formula.Template)
			} else {
				titles = append(titles, fmt.Sprintf("[%s] 标题 %d", formula.Type, i+1))
			}
		}
	}

	// 如果没有足够的标题，生成通用标题
	for len(titles) < count {
		titles = append(titles, fmt.Sprintf("标题 %d", len(titles)+1))
	}

	return titles
}

// ExtractQuotes 从文章中提取金句
func (g *articleGenerator) ExtractQuotes(article string, style *WriterStyle) []string {
	quotes := []string{}

	// 如果有预定义的金句模板，使用它们
	if len(style.QuoteTemplates) > 0 {
		quotes = append(quotes, style.QuoteTemplates...)
	}

	// 如果没有预定义模板，从文章中提取
	if len(quotes) == 0 {
		quotes = g.extractQuotesFromArticle(article, 5)
	}

	return quotes
}

// extractQuotesFromArticle 从文章中提取金句
func (g *articleGenerator) extractQuotesFromArticle(article string, count int) []string {
	quotes := []string{}

	// 按行分割
	lines := strings.Split(article, "\n")

	// 查找斜体内容（可能是金句）
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否是斜体标记的内容
		if strings.HasPrefix(line, "*") && strings.HasSuffix(line, "*") {
			content := strings.Trim(line, "*")
			if len(content) > 10 && len(content) < 100 {
				quotes = append(quotes, content)
			}
		}

		// 检查是否是粗体标记的内容
		if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") {
			content := strings.TrimPrefix(strings.TrimSuffix(line, "**"), "**")
			if len(content) > 10 && len(content) < 100 {
				quotes = append(quotes, content)
			}
		}

		if len(quotes) >= count {
			break
		}
	}

	return quotes
}

// IsAIRequest 检查结果是否是 AI 请求
func IsAIRequest(result *GenerateResult) bool {
	if result == nil {
		return false
	}
	if result.Status != "" {
		return result.Status == action.StatusActionRequired
	}
	if result.Prompt != "" {
		return true
	}
	return result.Error != "" && strings.HasPrefix(result.Error, "AI_MODE_REQUEST:")
}

// ExtractAIRequest 从结果中提取 AI 请求
func ExtractAIRequest(result *GenerateResult) string {
	if result == nil {
		return ""
	}
	if result.Status != "" {
		if result.Status == action.StatusActionRequired {
			return result.Prompt
		}
		return ""
	}
	if result.Prompt != "" {
		return result.Prompt
	}
	if strings.HasPrefix(result.Error, "AI_MODE_REQUEST:") {
		return strings.TrimPrefix(result.Error, "AI_MODE_REQUEST:")
	}
	return ""
}

// CompleteAIRequest 完成 AI 请求（AI 返回结果后调用）
func CompleteAIRequest(article string, result *GenerateResult) *GenerateResult {
	if result == nil {
		return &GenerateResult{
			Status:    action.StatusFailed,
			Action:    action.ActionWrite,
			Retryable: false,
			Success:   false,
			Error:     "结果为空",
		}
	}

	result.Article = article
	result.Error = ""
	result.Status = action.StatusCompleted
	result.Action = action.ActionWrite
	result.Retryable = false
	result.Success = true

	// 提取金句
	if result.Style != nil {
		result.Quotes = (&articleGenerator{}).ExtractQuotes(article, result.Style)
	}

	return result
}

// EstimateWordCount 估算字数
func EstimateWordCount(content string) int {
	// 简单估算：去除空格后按字符数计算
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\t", "")
	return len(content)
}

// ValidateInput 验证输入
func ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return NewInvalidInputError("输入内容为空")
	}

	if len(input) < 10 {
		return NewInvalidInputError("输入内容太短，请提供更多细节")
	}

	return nil
}

// BuildPromptForAI 为 AI 构建提示词（供外部使用）
func BuildPromptForAI(style *WriterStyle, userInput string, inputType InputType, articleType ArticleType) string {
	req := &GenerateRequest{
		Style:       style,
		UserInput:   userInput,
		InputType:   inputType,
		ArticleType: articleType,
		Length:      LengthMedium,
	}

	gen := NewGenerator()
	result := gen.Generate(req)

	if IsAIRequest(result) {
		return ExtractAIRequest(result)
	}

	return result.Prompt
}

// FormatInputType 格式化输入类型描述
func FormatInputType(inputType InputType) string {
	switch inputType {
	case InputTypeIdea:
		return "观点/想法"
	case InputTypeFragment:
		return "内容片段"
	case InputTypeOutline:
		return "大纲"
	case InputTypeTitle:
		return "标题扩展"
	default:
		return "其他"
	}
}

// FormatArticleType 格式化文章类型描述
func FormatArticleType(articleType ArticleType) string {
	switch articleType {
	case ArticleTypeEssay:
		return "散文"
	case ArticleTypeCommentary:
		return "评论"
	case ArticleTypeStory:
		return "故事"
	case ArticleTypeTutorial:
		return "教程"
	case ArticleTypeReview:
		return "评测"
	case ArticleType随笔:
		return "随笔"
	default:
		return "文章"
	}
}

// FormatLength 格式化长度描述
func FormatLength(length Length) string {
	switch length {
	case LengthShort:
		return "短文 (800-1200字)"
	case LengthMedium:
		return "中文 (1500-2500字)"
	case LengthLong:
		return "长文 (3000-5000字)"
	default:
		return "中等长度"
	}
}

// GetInputTypeFromString 从字符串获取输入类型
func GetInputTypeFromString(s string) InputType {
	switch strings.ToLower(s) {
	case "idea", "观点", "想法":
		return InputTypeIdea
	case "fragment", "片段", "内容":
		return InputTypeFragment
	case "outline", "大纲":
		return InputTypeOutline
	case "title", "标题":
		return InputTypeTitle
	default:
		return InputTypeIdea
	}
}

// GetArticleTypeFromString 从字符串获取文章类型
func GetArticleTypeFromString(s string) ArticleType {
	switch strings.ToLower(s) {
	case "essay", "散文":
		return ArticleTypeEssay
	case "commentary", "评论":
		return ArticleTypeCommentary
	case "story", "故事":
		return ArticleTypeStory
	case "tutorial", "教程":
		return ArticleTypeTutorial
	case "review", "评测":
		return ArticleTypeReview
	case "suibi", "随笔":
		return ArticleType随笔
	default:
		return ArticleTypeEssay
	}
}

// GetLengthFromString 从字符串获取长度
func GetLengthFromString(s string) Length {
	switch strings.ToLower(s) {
	case "short", "短", "短文":
		return LengthShort
	case "long", "长", "长文":
		return LengthLong
	default:
		return LengthMedium
	}
}
