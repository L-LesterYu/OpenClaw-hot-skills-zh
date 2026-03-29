// Package humanizer provides AI writing trace removal functionality
package humanizer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
)

// Humanizer 去痕处理器
type Humanizer struct{}

// NewHumanizer 创建去痕处理器
func NewHumanizer() *Humanizer {
	return &Humanizer{}
}

// Humanize 执行去痕处理
// 这个方法构建 AI 请求，由外部（Claude）执行实际处理
func (h *Humanizer) Humanize(req *HumanizeRequest) *HumanizeResult {
	// 验证输入
	if strings.TrimSpace(req.Content) == "" {
		return &HumanizeResult{
			Status:    action.StatusFailed,
			Action:    HumanizeActionAIRequest,
			Retryable: false,
			Success:   false,
			Error:     "输入内容为空",
		}
	}

	// 设置默认强度
	if req.Intensity == "" {
		req.Intensity = IntensityMedium
	}

	return h.buildAIRequestResult(req)
}

// HumanizeWithResult 直接处理并返回结果（用于 CLI 独立使用）
func (h *Humanizer) HumanizeWithResult(req *HumanizeRequest) *HumanizeResult {
	// 验证输入
	if strings.TrimSpace(req.Content) == "" {
		return &HumanizeResult{
			Status:    action.StatusFailed,
			Action:    HumanizeActionAIRequest,
			Retryable: false,
			Success:   false,
			Error:     "输入内容为空",
		}
	}

	// 设置默认强度
	if req.Intensity == "" {
		req.Intensity = IntensityMedium
	}

	return h.buildAIRequestResult(req)
}

// BuildAIRequestForAI 构建供 AI 使用的请求
// 返回提示词，由 Claude 执行实际处理
func (h *Humanizer) BuildAIRequestForAI(req *HumanizeRequest) string {
	return BuildPrompt(req)
}

// ParseAIResponse 解析 AI 返回的结果
func (h *Humanizer) ParseAIResponse(aiResponse string, req *HumanizeRequest) *HumanizeResult {
	result := &HumanizeResult{
		Success:   true,
		Status:    action.StatusCompleted,
		Action:    HumanizeActionCompleted,
		Retryable: false,
	}

	// 尝试解析结构化输出
	parsed := h.parseStructuredResponse(aiResponse)
	if parsed != nil {
		result.Content = parsed.Content
		result.Report = parsed.Report
		result.Changes = parsed.Changes
		result.Score = parsed.Score
		result.Success = true
		result.Action = HumanizeActionCompleted
		result.Status = action.StatusCompleted
		result.Retryable = false
		return result
	}

	// 如果无法解析结构化输出，尝试提取文本
	content := h.extractContent(aiResponse)
	if content != "" {
		result.Content = content
		result.Success = true
		result.Action = HumanizeActionCompleted
		result.Status = action.StatusCompleted
		result.Retryable = false
		return result
	}

	// 如果都失败，返回原文
	result.Success = false
	result.Content = req.Content
	result.Error = "无法解析 AI 返回结果，已返回原始文本"
	result.Status = action.StatusFailed
	result.Action = HumanizeActionCompleted
	result.Retryable = false
	return result
}

func (h *Humanizer) buildAIRequestResult(req *HumanizeRequest) *HumanizeResult {
	return &HumanizeResult{
		Success:   true,
		Status:    action.StatusActionRequired,
		Action:    HumanizeActionAIRequest,
		Retryable: false,
		Prompt:    BuildPrompt(req),
	}
}

// parseStructuredResponse 解析结构化响应
func (h *Humanizer) parseStructuredResponse(response string) *HumanizeResult {
	// 查找各个部分
	content := h.extractSection(response, "# 人性化后的文本", "# 修改说明", "# 处理结果")
	if content == "" {
		content = h.extractSection(response, "# 人性化后的文本", "# 修改说明", "")
	}

	changes := h.extractChanges(response)
	score := h.extractScore(response)
	report := h.extractSection(response, "# 修改说明", "# 质量评分", "# 处理结果")
	if report == "" {
		report = h.extractSection(response, "# 修改说明", "# 质量评分", "")
	}

	if content == "" {
		return nil
	}

	return &HumanizeResult{
		Success: true,
		Content: content,
		Changes: changes,
		Score:   score,
		Report:  strings.TrimSpace(report),
	}
}

// extractContent 提取内容部分
func (h *Humanizer) extractContent(response string) string {
	// 移除 markdown 代码块标记
	content := response

	// 移除 ```markdown 和 ``` 标记
	re := regexp.MustCompile("```(?:markdown)?\n?")
	content = re.ReplaceAllString(content, "")

	// 移除标题部分
	lines := strings.Split(content, "\n")
	var resultLines []string
	inContent := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# 人性化后的文本") ||
			strings.HasPrefix(line, "# 处理结果") ||
			strings.HasPrefix(line, "# Result") {
			inContent = true
			continue
		}
		if strings.HasPrefix(line, "# 修改说明") ||
			strings.HasPrefix(line, "# 质量评分") ||
			strings.HasPrefix(line, "# Changes") ||
			strings.HasPrefix(line, "# Score") {
			break
		}
		if inContent || (!strings.HasPrefix(line, "#") && line != "") {
			if !strings.HasPrefix(line, "#") {
				resultLines = append(resultLines, line)
			}
		}
	}

	return strings.TrimSpace(strings.Join(resultLines, "\n"))
}

// extractSection 提取指定章节内容
func (h *Humanizer) extractSection(response, sectionStart, sectionEnd1, sectionEnd2 string) string {
	lines := strings.Split(response, "\n")
	var resultLines []string
	inSection := false

	for _, line := range lines {
		if strings.HasPrefix(line, sectionStart) {
			inSection = true
			continue
		}
		if inSection {
			if (sectionEnd1 != "" && strings.HasPrefix(line, sectionEnd1)) ||
				(sectionEnd2 != "" && strings.HasPrefix(line, sectionEnd2)) ||
				strings.HasPrefix(line, "# ") {
				break
			}
			resultLines = append(resultLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(resultLines, "\n"))
}

// extractChanges 提取修改说明
func (h *Humanizer) extractChanges(response string) []Change {
	// 简化实现：从修改说明中提取信息
	section := h.extractSection(response, "# 修改说明", "# 质量评分", "# 处理结果")
	if section == "" {
		return nil
	}

	// 尝试解析 JSON 格式的修改记录
	var changes []Change
	if err := json.Unmarshal([]byte(section), &changes); err == nil {
		return changes
	}

	// 如果不是 JSON，返回一个汇总记录
	return []Change{
		{
			Type:     "summary",
			Original: "见修改说明",
			Revised:  section,
			Reason:   "综合处理",
		},
	}
}

// extractScore 提取质量评分
func (h *Humanizer) extractScore(response string) *Score {
	section := h.extractSection(response, "# 质量评分", "", "")
	if section == "" {
		return nil
	}

	// 尝试解析表格格式的评分
	lines := strings.Split(section, "\n")
	score := &Score{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "|") {
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				dimension := strings.TrimSpace(parts[1])
				valueStr := strings.TrimSpace(parts[2])

				// 提取数字
				re := regexp.MustCompile(`(\d+)`)
				matches := re.FindStringSubmatch(valueStr)
				if len(matches) > 1 {
					value, err := strconv.Atoi(matches[1])
					if err != nil {
						continue
					}

					switch strings.ToLower(dimension) {
					case "直接性":
						score.Directness = value
					case "节奏":
						score.Rhythm = value
					case "信任度":
						score.Trust = value
					case "真实性":
						score.Authenticity = value
					case "精炼度":
						score.Conciseness = value
					case "总分":
						score.Total = value
					}
				}
			}
		}
	}

	// 如果没有解析到总分，计算总分
	if score.Total == 0 {
		score.Total = score.Directness + score.Rhythm + score.Trust +
			score.Authenticity + score.Conciseness
	}

	if score.Total == 0 {
		return nil
	}

	return score
}

// GetSummary 获取处理摘要
func (h *Humanizer) GetSummary(result *HumanizeResult) string {
	if !result.Success {
		return fmt.Sprintf("[X] 处理失败: %s", result.Error)
	}

	var sb strings.Builder
	sb.WriteString("[OK] 处理完成\n\n")

	if result.Report != "" {
		sb.WriteString(fmt.Sprintf("[修改说明]\n%s\n\n", result.Report))
	}

	if result.HasChanges() {
		sb.WriteString(fmt.Sprintf("[修改] 修改了 %d 处\n\n", result.ChangeCount()))
	}

	if result.Score != nil {
		sb.WriteString(fmt.Sprintf("[评分] 质量评分: %d/50 - %s\n\n", result.Score.Total, result.Score.Rating()))
		sb.WriteString("| 维度 | 得分 |\n")
		sb.WriteString("|------|------|\n")
		sb.WriteString(fmt.Sprintf("| 直接性 | %d/10 |\n", result.Score.Directness))
		sb.WriteString(fmt.Sprintf("| 节奏 | %d/10 |\n", result.Score.Rhythm))
		sb.WriteString(fmt.Sprintf("| 信任度 | %d/10 |\n", result.Score.Trust))
		sb.WriteString(fmt.Sprintf("| 真实性 | %d/10 |\n", result.Score.Authenticity))
		sb.WriteString(fmt.Sprintf("| 精炼度 | %d/10 |\n", result.Score.Conciseness))
	}

	return sb.String()
}

// BuildConvertRequest 构建转换请求（兼容 converter 模块接口）
func (h *Humanizer) BuildConvertRequest(content string, settings map[string]interface{}) *AIConvertRequest {
	req := &HumanizeRequest{
		Content: content,
	}

	// 解析设置
	if intensity, ok := settings["intensity"].(string); ok {
		req.Intensity = ParseIntensity(intensity)
	}
	if preserveStyle, ok := settings["preserve_style"].(bool); ok {
		req.PreserveStyle = preserveStyle
	}
	if originalStyle, ok := settings["original_style"].(string); ok {
		req.OriginalStyle = originalStyle
	}
	if showChanges, ok := settings["show_changes"].(bool); ok {
		req.ShowChanges = showChanges
	}
	if includeScore, ok := settings["include_score"].(bool); ok {
		req.IncludeScore = includeScore
	}
	if focusOn, ok := settings["focus_on"].([]string); ok {
		req.FocusOn = ParseFocusPattern(focusOn)
	}

	return &AIConvertRequest{
		Content: content,
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
