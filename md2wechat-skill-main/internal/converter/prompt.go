package converter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// PromptTemplate Prompt 模板结构
type PromptTemplate struct {
	Name        string            // 模板名称
	Description string            // 模板描述
	Template    string            // 模板内容
	Variables   []string          // 支持的变量列表
	Metadata    map[string]string // 元数据
}

// PromptVariable Prompt 变量
type PromptVariable struct {
	Name         string // 变量名
	Description  string // 变量描述
	DefaultValue string // 默认值
	Required     bool   // 是否必填
}

// PromptBuilder Prompt 构建器
type PromptBuilder struct {
	templates map[string]*PromptTemplate
	variables map[string]*PromptVariable
}

// NewPromptBuilder 创建 Prompt 构建器
func NewPromptBuilder() *PromptBuilder {
	pb := &PromptBuilder{
		templates: make(map[string]*PromptTemplate),
		variables: make(map[string]*PromptVariable),
	}
	pb.initBuiltinVariables()
	return pb
}

// initBuiltinVariables 初始化内置变量
func (pb *PromptBuilder) initBuiltinVariables() {
	pb.variables["{{MARKDOWN}}"] = &PromptVariable{
		Name:         "MARKDOWN",
		Description:  "Markdown 内容",
		DefaultValue: "",
		Required:     true,
	}
	pb.variables["{{THEME_NAME}}"] = &PromptVariable{
		Name:         "THEME_NAME",
		Description:  "主题名称",
		DefaultValue: "default",
		Required:     false,
	}
	pb.variables["{{TITLE}}"] = &PromptVariable{
		Name:         "TITLE",
		Description:  "文章标题",
		DefaultValue: "未命名文章",
		Required:     false,
	}
	pb.variables["{{FONT_SIZE}}"] = &PromptVariable{
		Name:         "FONT_SIZE",
		Description:  "字体大小",
		DefaultValue: "16px",
		Required:     false,
	}
	pb.variables["{{LINE_HEIGHT}}"] = &PromptVariable{
		Name:         "LINE_HEIGHT",
		Description:  "行高",
		DefaultValue: "1.75",
		Required:     false,
	}
	pb.variables["{{PRIMARY_COLOR}}"] = &PromptVariable{
		Name:         "PRIMARY_COLOR",
		Description:  "主色调",
		DefaultValue: "#4a413d",
		Required:     false,
	}
	pb.variables["{{SECONDARY_COLOR}}"] = &PromptVariable{
		Name:         "SECONDARY_COLOR",
		Description:  "副强调色",
		DefaultValue: "#c06b4d",
		Required:     false,
	}
	pb.variables["{{BACKGROUND_COLOR}}"] = &PromptVariable{
		Name:         "BACKGROUND_COLOR",
		Description:  "背景色",
		DefaultValue: "#faf9f5",
		Required:     false,
	}
	pb.variables["{{ACCENT_COLOR}}"] = &PromptVariable{
		Name:         "ACCENT_COLOR",
		Description:  "强调色",
		DefaultValue: "#d97758",
		Required:     false,
	}
}

// AddTemplate 添加模板
func (pb *PromptBuilder) AddTemplate(tpl *PromptTemplate) error {
	if tpl.Name == "" {
		return fmt.Errorf("template name is required")
	}
	if tpl.Template == "" {
		return fmt.Errorf("template content is required")
	}

	// 提取模板中的变量
	vars := pb.extractVariables(tpl.Template)
	tpl.Variables = vars

	pb.templates[tpl.Name] = tpl
	return nil
}

// extractVariables 从模板中提取变量
func (pb *PromptBuilder) extractVariables(content string) []string {
	re := regexp.MustCompile(`\{\{([A-Z_]+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	varMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			varMap[match[1]] = true
		}
	}

	var vars []string
	for v := range varMap {
		vars = append(vars, "{{"+v+"}}")
	}
	return vars
}

// BuildPrompt 构建完整的 Prompt
func (pb *PromptBuilder) BuildPrompt(templateName string, vars map[string]string) (string, error) {
	tpl, ok := pb.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	result := tpl.Template

	// 替换变量
	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		if !strings.Contains(result, placeholder) {
			// 尝试带大括号的格式
			placeholder = key
		}
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// 替换未提供的变量为默认值
	for _, varRef := range tpl.Variables {
		if strings.Contains(result, varRef) {
			// 提取变量名
			varName := strings.Trim(varRef, "{}")
			if _, provided := vars[varName]; !provided {
				if variable, ok := pb.variables[varRef]; ok {
					result = strings.ReplaceAll(result, varRef, variable.DefaultValue)
				}
			}
		}
	}

	return result, nil
}

// BuildPromptFromTheme 从主题构建 Prompt
func (pb *PromptBuilder) BuildPromptFromTheme(theme *Theme, markdown string, vars map[string]string) (string, error) {
	if theme.Type != "ai" {
		return "", fmt.Errorf("theme '%s' is not an AI theme", theme.Name)
	}

	// 设置默认变量
	if vars == nil {
		vars = make(map[string]string)
	}
	if _, ok := vars["MARKDOWN"]; !ok {
		vars["MARKDOWN"] = markdown
	}
	if _, ok := vars["THEME_NAME"]; !ok {
		vars["THEME_NAME"] = theme.Name
	}

	// 使用主题的 Prompt 作为模板
	prompt := theme.Prompt

	// 替换 {{MARKDOWN}} 变量
	if strings.Contains(prompt, "{{MARKDOWN}}") {
		prompt = strings.ReplaceAll(prompt, "{{MARKDOWN}}", markdown)
	} else {
		// 如果没有占位符，追加到末尾
		prompt = prompt + "\n\n```\n" + markdown + "\n```"
	}

	// 替换其他变量
	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}

	return prompt, nil
}

// ValidateTemplate 验证模板
func (pb *PromptBuilder) ValidateTemplate(templateName string) error {
	tpl, ok := pb.templates[templateName]
	if !ok {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// 检查必填变量
	for _, varRef := range tpl.Variables {
		if variable, ok := pb.variables[varRef]; ok {
			if variable.Required && !strings.Contains(tpl.Template, varRef) {
				return fmt.Errorf("required variable %s not found in template", varRef)
			}
		}
	}

	return nil
}

// ValidatePromptContent 验证 Prompt 内容是否符合要求
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

func ValidatePromptContent(prompt string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// 检查是否包含关键规则
	requiredRules := []struct {
		pattern string
		name    string
	}{
		{`内联.*style|inline.*style`, "内联样式说明"},
		{`IMG:\d+|IMG:index|图片.*占位`, "图片占位符说明"},
		{`HTML.*标签|HTML.*tag`, "HTML 标签说明"},
	}

	for _, rule := range requiredRules {
		matched, _ := regexp.MatchString(rule.pattern, prompt)
		if !matched {
			result.Warnings = append(result.Warnings, fmt.Sprintf("可能缺少 %s 的说明", rule.name))
		}
	}

	// 检查是否包含禁用的标签
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onload=`,
		`onerror=`,
	}

	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(pattern, prompt)
		if matched {
			result.Errors = append(result.Errors, fmt.Sprintf("包含不安全的内容: %s", pattern))
			result.Valid = false
		}
	}

	return result
}

// ExportPrompt 导出 Prompt
type ExportOptions struct {
	Format        string // json, text, markdown
	IncludeHeader bool   // 是否包含头部信息
	IncludeFooter bool   // 是否包含尾部说明
}

// ExportPrompt 导出 Prompt
func (pb *PromptBuilder) ExportPrompt(templateName string, vars map[string]string, opts *ExportOptions) (string, error) {
	if opts == nil {
		opts = &ExportOptions{
			Format:        "text",
			IncludeHeader: true,
			IncludeFooter: false,
		}
	}

	tpl, ok := pb.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	prompt, err := pb.BuildPrompt(templateName, vars)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if opts.IncludeHeader {
		switch opts.Format {
		case "markdown", "md":
			buf.WriteString(fmt.Sprintf("# %s\n\n", tpl.Name))
			if tpl.Description != "" {
				buf.WriteString(fmt.Sprintf("%s\n\n", tpl.Description))
			}
		case "json":
			buf.WriteString("{\n")
			buf.WriteString(fmt.Sprintf("  \"name\": \"%s\",\n", tpl.Name))
			buf.WriteString(fmt.Sprintf("  \"description\": \"%s\",\n", tpl.Description))
			buf.WriteString("  \"prompt\": \"")
		default:
			buf.WriteString(fmt.Sprintf("# Prompt: %s\n", tpl.Name))
			if tpl.Description != "" {
				buf.WriteString(fmt.Sprintf("# %s\n\n", tpl.Description))
			}
		}
	}

	switch opts.Format {
	case "json":
		// JSON 格式需要转义
		jsonPrompt := strings.ReplaceAll(prompt, "\n", "\\n")
		jsonPrompt = strings.ReplaceAll(jsonPrompt, "\"", "\\\"")
		buf.WriteString(jsonPrompt + "\"\n}")
	default:
		buf.WriteString(prompt)
	}

	if opts.IncludeFooter && opts.Format != "json" {
		buf.WriteString("\n\n---\n")
		buf.WriteString(fmt.Sprintf("Template: %s\n", tpl.Name))
		if len(tpl.Variables) > 0 {
			buf.WriteString("Variables: " + strings.Join(tpl.Variables, ", ") + "\n")
		}
	}

	return buf.String(), nil
}

// GetTemplate 获取模板
func (pb *PromptBuilder) GetTemplate(name string) (*PromptTemplate, error) {
	tpl, ok := pb.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return tpl, nil
}

// ListTemplates 列出所有模板
func (pb *PromptBuilder) ListTemplates() []string {
	var names []string
	for name := range pb.templates {
		names = append(names, name)
	}
	return names
}

// GetVariable 获取变量定义
func (pb *PromptBuilder) GetVariable(name string) (*PromptVariable, error) {
	v, ok := pb.variables[name]
	if !ok {
		return nil, fmt.Errorf("variable not found: %s", name)
	}
	return v, nil
}

// ListVariables 列出所有变量
func (pb *PromptBuilder) ListVariables() []string {
	var names []string
	for name := range pb.variables {
		names = append(names, name)
	}
	return names
}

// BuildPromptWithTemplate 使用 Go template 语法构建 Prompt
func (pb *PromptBuilder) BuildPromptWithTemplate(templateContent string, vars map[string]string) (string, error) {
	// 确保所有变量名都是有效的 Go template 标识符
	cleanVars := make(map[string]any)
	for k, v := range vars {
		// 将 {{VAR_NAME}} 转换为 .VarName 格式
		cleanKey := strings.ToLower(strings.ReplaceAll(k, "_", ""))
		cleanVars[cleanKey] = v
	}

	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cleanVars); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

// ParseMarkdownTitle 解析 Markdown 标题
func ParseMarkdownTitle(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// 移除 # 号并清理
			title := strings.TrimLeft(line, "#")
			title = strings.TrimSpace(title)
			if title != "" {
				return title
			}
		}
	}
	return "未命名文章"
}

// EstimateTokenCount 估算 token 数量（粗略估计：中文约 1.5 字符/token，英文约 4 字符/token）
func EstimateTokenCount(text string) int {
	chineseChars := 0
	otherChars := 0

	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			chineseChars++
		} else {
			otherChars++
		}
	}

	return (chineseChars / 1) + (otherChars / 4)
}
