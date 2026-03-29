package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/assets"
	"gopkg.in/yaml.v3"
)

const themesDirEnvVar = "MD2WECHAT_THEMES_DIR"

// Theme 主题定义
type Theme struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"` // "api" | "ai"
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	StyleInfo   ThemeStyleInfo    `yaml:"style_info,omitempty"`
	Colors      map[string]string `yaml:"colors,omitempty"`
	APITheme    string            `yaml:"api_theme,omitempty"`
	Prompt      string            `yaml:"prompt,omitempty"`
}

// ThemeStyleInfo 主题风格信息
type ThemeStyleInfo struct {
	Mood    string `yaml:"mood"`
	Colors  string `yaml:"colors"`
	BestFor string `yaml:"best_for"`
}

// ThemeManager 主题管理器
type ThemeManager struct {
	themes map[string]Theme
}

// NewThemeManager 创建主题管理器
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		themes: make(map[string]Theme),
	}
}

// LoadThemes 从 YAML 文件加载主题
func (tm *ThemeManager) LoadThemes() error {
	if err := tm.loadBuiltinThemes(); err != nil {
		return fmt.Errorf("load builtin themes: %w", err)
	}

	for _, themeDir := range tm.getThemeDirs() {
		if err := tm.loadThemesFromDir(themeDir); err != nil {
			return err
		}
	}

	return nil
}

// loadThemeFromFile 从文件加载单个主题
func (tm *ThemeManager) loadThemeFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return tm.loadThemeData(data)
}

func (tm *ThemeManager) loadThemeData(data []byte) error {
	var theme Theme
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	// 验证主题
	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}
	if theme.Type == "" {
		theme.Type = "ai" // 默认为 AI 模式
	}

	// 如果 description 为空，设置默认值
	if theme.Description == "" {
		theme.Description = theme.Name
	}

	tm.themes[theme.Name] = theme
	return nil
}

func (tm *ThemeManager) loadThemesFromDir(themeDir string) error {
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read theme directory %s: %w", themeDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}
		themePath := filepath.Join(themeDir, entry.Name())
		if err := tm.loadThemeFromFile(themePath); err != nil {
			return fmt.Errorf("load theme from %s: %w", themePath, err)
		}
	}

	return nil
}

func (tm *ThemeManager) loadBuiltinThemes() error {
	names, err := assets.ListBuiltinThemes()
	if err != nil {
		return err
	}

	for _, name := range names {
		data, err := assets.ReadBuiltinTheme(name)
		if err != nil {
			return fmt.Errorf("read builtin theme %s: %w", name, err)
		}
		if err := tm.loadThemeData(data); err != nil {
			return fmt.Errorf("load builtin theme %s: %w", name, err)
		}
	}

	return nil
}

// getThemeDirs 获取主题目录，优先级高的目录排在前面，由后加载覆盖内置资产
func (tm *ThemeManager) getThemeDirs() []string {
	dirs := make([]string, 0, 3)
	add := func(dir string) {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return
		}
		for _, existing := range dirs {
			if existing == dir {
				return
			}
		}
		dirs = append(dirs, dir)
	}

	add(os.Getenv(themesDirEnvVar))
	add("themes")

	homeDir, _ := os.UserHomeDir()
	add(filepath.Join(homeDir, ".config", "md2wechat", "themes"))

	return dirs
}

// LoadTheme 加载单个主题（支持自定义路径）
func (tm *ThemeManager) LoadTheme(path string) error {
	return tm.loadThemeFromFile(path)
}

// GetTheme 获取主题
func (tm *ThemeManager) GetTheme(name string) (*Theme, error) {
	// 如果主题未加载，尝试从文件加载
	if _, ok := tm.themes[name]; !ok {
		if err := tm.LoadThemes(); err != nil {
			return nil, fmt.Errorf("theme not found: %s (load error: %w)", name, err)
		}
	}

	theme, ok := tm.themes[name]
	if !ok {
		return nil, fmt.Errorf("theme not found: %s", name)
	}
	return &theme, nil
}

// ListThemes 列出所有主题
func (tm *ThemeManager) ListThemes() []string {
	var names []string
	for name := range tm.themes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (tm *ThemeManager) ListThemeDefinitions() []Theme {
	themes := make([]Theme, 0, len(tm.themes))
	for _, theme := range tm.themes {
		themes = append(themes, theme)
	}
	sort.Slice(themes, func(i, j int) bool {
		if themes[i].Type == themes[j].Type {
			return themes[i].Name < themes[j].Name
		}
		return themes[i].Type < themes[j].Type
	})
	return themes
}

// ListAIThemes 列出所有 AI 主题
func (tm *ThemeManager) ListAIThemes() []string {
	var names []string
	for name, theme := range tm.themes {
		if theme.Type == "ai" {
			names = append(names, name)
		}
	}
	return names
}

// ListAPIThemes 列出所有 API 主题
func (tm *ThemeManager) ListAPIThemes() []string {
	var names []string
	for name, theme := range tm.themes {
		if theme.Type == "api" {
			names = append(names, name)
		}
	}
	return names
}

// GetAPITheme 获取 API 模式的主题名
func (tm *ThemeManager) GetAPITheme(name string) (string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}
	if theme.Type != "api" {
		return "", fmt.Errorf("theme '%s' is not an API theme", name)
	}
	return theme.APITheme, nil
}

// GetAIPrompt 获取 AI 模式的提示词
func (tm *ThemeManager) GetAIPrompt(name string) (string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}
	if theme.Type != "ai" {
		return "", fmt.Errorf("theme '%s' is not an AI theme", name)
	}
	if theme.Prompt == "" {
		return "", fmt.Errorf("theme '%s' has no prompt defined", name)
	}
	return theme.Prompt, nil
}

// BuildCustomAIPrompt 构建自定义 AI 提示词
func BuildCustomAIPrompt(customPrompt string) string {
	if customPrompt == "" {
		return customPrompt
	}

	// 确保包含基本规则
	baseRules := `

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->
5. 返回完整的 HTML，不需要其他说明文字

`

	if !strings.Contains(customPrompt, "重要规则") && !strings.Contains(customPrompt, "规则") {
		customPrompt += baseRules
	}

	if !strings.Contains(customPrompt, "请转换") {
		customPrompt += "\n\n请转换以下 Markdown内容："
	}

	return customPrompt
}

// IsAPITheme 检查是否是 API 主题
func (tm *ThemeManager) IsAPITheme(name string) bool {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return false
	}
	return theme.Type == "api"
}

// IsAITheme 检查是否是 AI 主题
func (tm *ThemeManager) IsAITheme(name string) bool {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return false
	}
	return theme.Type == "ai"
}

// GetThemeDescription 获取主题描述
func (tm *ThemeManager) GetThemeDescription(name string) string {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "未知主题"
	}
	return theme.Description
}

// GetThemeColors 获取主题颜色配置
func (tm *ThemeManager) GetThemeColors(name string) (map[string]string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return nil, err
	}
	return theme.Colors, nil
}

// ReloadThemes 重新加载所有主题
func (tm *ThemeManager) ReloadThemes() error {
	tm.themes = make(map[string]Theme)
	return tm.LoadThemes()
}

// GetThemeInfo 获取主题完整信息（用于调试）
func (tm *ThemeManager) GetThemeInfo(name string) (*Theme, error) {
	return tm.GetTheme(name)
}

// EnsureLoaded 确保主题已加载
func (tm *ThemeManager) EnsureLoaded() error {
	if len(tm.themes) == 0 {
		return tm.LoadThemes()
	}
	return nil
}
