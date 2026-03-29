// Package humanizer provides AI writing trace removal functionality
package humanizer

import "github.com/geekjourneyx/md2wechat-skill/internal/action"

// HumanizeIntensity 去痕强度
type HumanizeIntensity string

const (
	IntensityGentle     HumanizeIntensity = "gentle"     // 温和：只处理明显的
	IntensityMedium     HumanizeIntensity = "medium"     // 中等：默认
	IntensityAggressive HumanizeIntensity = "aggressive" // 激进：深度处理
)

// String returns the string representation
func (h HumanizeIntensity) String() string {
	return string(h)
}

// Description returns the description of intensity level
func (h HumanizeIntensity) Description() string {
	switch h {
	case IntensityGentle:
		return "温和处理，只修改最明显、最确定的问题"
	case IntensityMedium:
		return "平衡处理，保留合理的表达，去除明显 AI 痕迹"
	case IntensityAggressive:
		return "深度审查，最大化去除 AI 痕迹，大幅改写"
	default:
		return "中等强度"
	}
}

// FocusPattern 去痕聚焦模式（24种模式的分类）
type FocusPattern string

const (
	PatternContent       FocusPattern = "content"       // 内容模式
	PatternLanguage      FocusPattern = "language"      // 语言语法
	PatternStyle         FocusPattern = "style"         // 风格模式
	PatternFiller        FocusPattern = "filler"        // 填充词回避
	PatternCollaboration FocusPattern = "collaboration" // 协作交流痕迹
)

// AllFocusPatterns 返回所有可用的聚焦模式
func AllFocusPatterns() []FocusPattern {
	return []FocusPattern{
		PatternContent,
		PatternLanguage,
		PatternStyle,
		PatternFiller,
		PatternCollaboration,
	}
}

// HumanizeRequest 去痕请求
type HumanizeRequest struct {
	// 输入
	Content string `json:"content"` // 待处理文本

	// 处理控制
	Intensity HumanizeIntensity `json:"intensity,omitempty"` // 处理强度
	FocusOn   []FocusPattern    `json:"focus_on,omitempty"`  // 重点处理的模式分类（为空则全部）

	// 行为控制
	PreserveStyle bool `json:"preserve_style,omitempty"` // 保持原有风格特征（风格优先）
	ShowChanges   bool `json:"show_changes,omitempty"`   // 返回修改对比
	IncludeScore  bool `json:"include_score,omitempty"`  // 返回质量评分

	// 源信息（用于更好的处理）
	SourceHint    string `json:"source_hint,omitempty"`    // "ai-generated" / "human-written" / "unknown"
	OriginalStyle string `json:"original_style,omitempty"` // 如果使用了写作风格，传入风格名
}

// HumanizeResult 去痕结果
type HumanizeResult struct {
	Success   bool          `json:"success"`
	Status    action.Status `json:"status,omitempty"`
	Action    string        `json:"action,omitempty"`
	Retryable bool          `json:"retryable,omitempty"`

	// 输出
	Content string `json:"content"` // 处理后的文本
	Prompt  string `json:"prompt,omitempty"`

	// 变更信息（可选）
	Changes []Change `json:"changes,omitempty"` // 修改详情
	Score   *Score   `json:"score,omitempty"`   // 质量评分
	Report  string   `json:"report,omitempty"`  // 处理报告（自然语言描述）

	// 错误信息
	Error string `json:"error,omitempty"`
}

const (
	HumanizeActionCompleted = action.ActionHumanize
	HumanizeActionAIRequest = action.ActionHumanize
)

// RequiresAI 是否需要外部 AI 执行
func (r *HumanizeResult) RequiresAI() bool {
	if r == nil {
		return false
	}
	if r.Status != "" {
		return r.Status == action.StatusActionRequired
	}
	if r.Prompt != "" {
		return true
	}
	return false
}

// HasChanges 是否有修改记录
func (r *HumanizeResult) HasChanges() bool {
	return len(r.Changes) > 0
}

// ChangeCount 修改数量
func (r *HumanizeResult) ChangeCount() int {
	return len(r.Changes)
}

// Change 单次修改记录
type Change struct {
	Type     string `json:"type"`     // 模式类型: "filler_phrase" / "ai_vocabulary" / ...
	Original string `json:"original"` // 原文
	Revised  string `json:"revised"`  // 修改后
	Position int    `json:"position"` // 在文中的大致位置（字符偏移）
	Reason   string `json:"reason"`   // 修改原因
}

// ChangeType 修改类型常量
const (
	ChangeTypeFillerPhrase      = "filler_phrase"       // 填充短语
	ChangeTypeAIVocabulary      = "ai_vocabulary"       // AI 词汇
	ChangeTypeFormulaic         = "formulaic_structure" // 公式化结构
	ChangeTypeOveremphasis      = "overemphasis"        // 过度强调
	ChangeTypeIngiAnalysis      = "ing_analysis"        // -ing 肤浅分析
	ChangeTypeVagueAttribution  = "vague_attribution"   // 模糊归因
	ChangeTypePromotional       = "promotional"         // 宣传性语言
	ChangeTypeDashOveruse       = "dash_overuse"        // 破折号过度
	ChangeTypeGenericConclusion = "generic_conclusion"  // 通用积极结论
	ChangeTypeOther             = "other"               // 其他
)

// Score 质量评分（基于 humanizer-zh 的 5 维度评估）
type Score struct {
	Total        int `json:"total"`        // 总分 /50
	Directness   int `json:"directness"`   // 直接性 /10
	Rhythm       int `json:"rhythm"`       // 节奏 /10
	Trust        int `json:"trust"`        // 信任度 /10
	Authenticity int `json:"authenticity"` // 真实性 /10
	Conciseness  int `json:"conciseness"`  // 精炼度 /10
}

// Rating 返回评分等级
func (s *Score) Rating() string {
	if s == nil {
		return "未评分"
	}
	switch {
	case s.Total >= 45:
		return "优秀 - 已去除 AI 痕迹"
	case s.Total >= 35:
		return "良好 - 仍有改进空间"
	case s.Total >= 25:
		return "一般 - 需要进一步修订"
	default:
		return "较差 - 建议重新处理"
	}
}

// AIConvertRequest AI 转换请求（用于传递给 Claude）
type AIConvertRequest struct {
	Content  string           // 原始内容
	Prompt   string           // 完整的提示词
	Settings HumanizeSettings // 处理设置
}

// HumanizeSettings 去痕设置
type HumanizeSettings struct {
	Intensity     HumanizeIntensity // 处理强度
	FocusOn       []FocusPattern    // 聚焦模式
	PreserveStyle bool              // 保持风格
	OriginalStyle string            // 原始风格名
	ShowChanges   bool              // 显示变更
	IncludeScore  bool              // 包含评分
}

// AIConvertResult AI 转换结果
type AIConvertResult struct {
	HTML    string // 生成的文本（复用字段名，实际是纯文本）
	Success bool
	Error   string
}
