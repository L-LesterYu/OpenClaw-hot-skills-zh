// Package writer provides assisted writing functionality with customizable creator styles
package writer

import (
	"context"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
)

// InputType 输入类型
type InputType string

const (
	InputTypeIdea     InputType = "idea"     // 观点/想法
	InputTypeFragment InputType = "fragment" // 内容片段
	InputTypeOutline  InputType = "outline"  // 大纲
	InputTypeTitle    InputType = "title"    // 标题扩展
)

// String 返回 InputType 的字符串表示
func (t InputType) String() string {
	return string(t)
}

// ArticleType 文章类型
type ArticleType string

const (
	ArticleTypeEssay      ArticleType = "essay"      // 散文
	ArticleTypeCommentary ArticleType = "commentary" // 评论
	ArticleTypeStory      ArticleType = "story"      // 故事
	ArticleTypeTutorial   ArticleType = "tutorial"   // 教程
	ArticleTypeReview     ArticleType = "review"     // 评论/评测
	ArticleType随笔         ArticleType = "suibi"      // 随笔
)

// String 返回 ArticleType 的字符串表示
func (t ArticleType) String() string {
	return string(t)
}

// Length 文章长度
type Length string

const (
	LengthShort  Length = "short"  // 800-1200字
	LengthMedium Length = "medium" // 1500-2500字
	LengthLong   Length = "long"   // 3000-5000字
)

// String 返回 Length 的字符串表示
func (l Length) String() string {
	return string(l)
}

// WriteRequest 写作请求
type WriteRequest struct {
	// 输入内容
	Input     string    // 用户输入（观点或内容）
	InputType InputType // 输入类型

	// 风格设置
	StyleName string // 作家风格名称

	// 文章设置
	ArticleType ArticleType // 文章类型
	Length      Length      // 期望长度

	// 可选内容
	Title        string            // 文章标题（可选）
	Context      map[string]string // 上下文信息
	CustomPrompt string            // 自定义提示词（可选）
}

// RefineRequest 润色请求
type RefineRequest struct {
	Content   string // 原始内容
	StyleName string // 作家风格名称
	Feedback  string // 用户反馈（可选）
}

// RefineResult 润色结果
type RefineResult struct {
	Refined     string        // 润色后的内容
	Changes     []string      // 变更说明
	BeforeAfter string        // 对比（可选）
	Status      action.Status `json:"status,omitempty"`
	Action      string        `json:"action,omitempty"`
	Retryable   bool          `json:"retryable,omitempty"`
	Prompt      string        `json:"prompt,omitempty"`
	Success     bool
	Error       string
}

// GenerateCoverRequest 生成封面请求
type GenerateCoverRequest struct {
	ArticleTitle   string
	ArticleContent string
	StyleName      string
}

// GenerateCoverResult 生成封面结果
type GenerateCoverResult struct {
	Prompt      string // 生成的提示词
	Explanation string // 隐喻解释
	ImageURL    string // 生成的图片 URL
	MediaID     string // 微信素材 ID
	MetaData    CoverMetaData
	Success     bool
	Error       string
}

// WriterStyle 写作风格定义
type WriterStyle struct {
	Name        string `yaml:"name"`         // 风格名称
	EnglishName string `yaml:"english_name"` // 英文标识（用于命令行）
	Category    string `yaml:"category"`     // 分类
	Description string `yaml:"description"`  // 描述
	Version     string `yaml:"version"`      // 版本

	// 核心写作 DNA
	CoreBeliefs []string `yaml:"core_beliefs,omitempty"` // 核心信念

	// 写作风格定义
	WritingStyle *WritingStyleDef `yaml:"writing_style,omitempty"`

	// 段落规则
	ParagraphRules *ParagraphRules `yaml:"paragraph_rules,omitempty"`

	// 格式化规范
	Formatting *FormattingDef `yaml:"formatting,omitempty"`

	// 标点节奏
	Punctuation *PunctuationDef `yaml:"punctuation_rhythm,omitempty"`

	// AI 提示词
	WritingPrompt string `yaml:"writing_prompt,omitempty"`

	// 标题公式库
	TitleFormulas []TitleFormula `yaml:"title_formulas,omitempty"`

	// 金句模板
	QuoteTemplates []string `yaml:"quote_templates,omitempty"`

	// 封面相关
	CoverPrompt      string   `yaml:"cover_prompt,omitempty"`
	CoverStyle       string   `yaml:"cover_style,omitempty"`
	CoverMood        string   `yaml:"cover_mood,omitempty"`
	CoverColorScheme []string `yaml:"cover_color_scheme,omitempty"`
}

// WritingStyleDef 写作风格定义
type WritingStyleDef struct {
	Tone        string `yaml:"tone"`        // 语气
	Voice       string `yaml:"voice"`       // 声音/口吻
	Perspective string `yaml:"perspective"` // 视角
}

// ParagraphRules 段落规则
type ParagraphRules struct {
	MaxLinesPerParagraph      int  `yaml:"max_lines_per_paragraph,omitempty"`      // 每段最大行数
	ImportantPointsStandalone bool `yaml:"important_points_standalone,omitempty"`  // 重要观点独立成段
	BlankLineBeforeTransition bool `yaml:"blank_line_before_transition,omitempty"` // 转折前空行
}

// FormattingDef 格式化定义
type FormattingDef struct {
	BoldFor       string `yaml:"bold_for,omitempty"`        // 粗体用于
	ItalicFor     string `yaml:"italic_for,omitempty"`      // 斜体用于
	QuoteMarksFor string `yaml:"quote_marks_for,omitempty"` // 引号用于
}

// PunctuationDef 标点节奏定义
type PunctuationDef struct {
	UsePeriodsForBreathing bool `yaml:"use_periods_for_breathing,omitempty"`    // 用句号制造呼吸感
	UseQuestionsReflection bool `yaml:"use_questions_for_reflection,omitempty"` // 用问号引发思考
	UseEmDashesEmphasis    bool `yaml:"use_em_dashes_for_emphasis,omitempty"`   // 用破折号强调
}

// TitleFormula 标题公式
type TitleFormula struct {
	Type     string   `yaml:"type"`               // 公式类型
	Template string   `yaml:"template,omitempty"` // 模板
	Examples []string `yaml:"examples,omitempty"` // 示例
}

// CoverMetaData 封面元数据
type CoverMetaData struct {
	CoreTheme    string // 核心主题
	CoreView     string // 核心观点
	Mood         string // 情绪基调
	VisualAnchor string // 视觉锚点
}

// StyleListResult 风格列表结果
type StyleListResult struct {
	Styles  []StyleSummary
	Success bool
	Error   string
}

// StyleSummary 风格摘要
type StyleSummary struct {
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	CoverStyle  string `json:"cover_style,omitempty"`
}

// AIGenerationRequest AI 生成请求（用于传递给 Claude）
type AIGenerationRequest struct {
	Context     context.Context
	Style       *WriterStyle
	UserInput   string
	InputType   InputType
	ArticleType ArticleType
	Length      Length
	Title       string
}

// AIGenerationResult AI 生成结果
type AIGenerationResult struct {
	Article    string
	Title      string
	Quotes     []string
	PromptUsed string
	Success    bool
	Error      string
}

// WriterError 写作错误
type WriterError struct {
	Code    string
	Message string
	Hint    string // 操作提示
	Err     error
}

func (e *WriterError) Error() string {
	msg := e.Code + ": " + e.Message
	if e.Hint != "" {
		msg += "\n提示: " + e.Hint
	}
	if e.Err != nil {
		msg += "\n" + e.Err.Error()
	}
	return msg
}

func (e *WriterError) Unwrap() error {
	return e.Err
}

// 错误代码常量
const (
	ErrCodeStyleNotFound      = "STYLE_NOT_FOUND"
	ErrCodeStyleLoadFailed    = "STYLE_LOAD_FAILED"
	ErrCodeInvalidInput       = "INVALID_INPUT"
	ErrCodeGenerationFailed   = "GENERATION_FAILED"
	ErrCodeCoverGenFailed     = "COVER_GEN_FAILED"
	ErrCodeInvalidInputType   = "INVALID_INPUT_TYPE"
	ErrCodeInvalidArticleType = "INVALID_ARTICLE_TYPE"
)

// NewStyleNotFoundError 创建风格未找到错误
func NewStyleNotFoundError(name string) *WriterError {
	return &WriterError{
		Code:    ErrCodeStyleNotFound,
		Message: "风格未找到: " + name,
		Hint:    "使用 `md2wechat styles` 查看可用风格，或在 writers/ 目录添加自定义风格",
	}
}

// NewInvalidInputError 创建无效输入错误
func NewInvalidInputError(reason string) *WriterError {
	return &WriterError{
		Code:    ErrCodeInvalidInput,
		Message: "输入无效: " + reason,
		Hint:    "请提供有效的观点、内容片段或大纲",
	}
}

// NewGenerationFailedError 创建生成失败错误
func NewGenerationFailedError(err error) *WriterError {
	return &WriterError{
		Code:    ErrCodeGenerationFailed,
		Message: "文章生成失败",
		Err:     err,
	}
}

// DefaultStyleName 默认风格名称
const DefaultStyleName = "dan-koe"

// InputTypes 所有支持的输入类型
var InputTypes = []InputType{
	InputTypeIdea,
	InputTypeFragment,
	InputTypeOutline,
	InputTypeTitle,
}

// ArticleTypes 所有支持的文章类型
var ArticleTypes = []ArticleType{
	ArticleTypeEssay,
	ArticleTypeCommentary,
	ArticleTypeStory,
	ArticleTypeTutorial,
	ArticleTypeReview,
	ArticleType随笔,
}

// Lengths 所有支持的长度
var Lengths = []Length{
	LengthShort,
	LengthMedium,
	LengthLong,
}
