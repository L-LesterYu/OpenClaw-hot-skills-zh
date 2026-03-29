// Package writer provides assisted writing functionality with customizable creator styles
package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/assets"
	"gopkg.in/yaml.v3"
)

const writersDirEnvVar = "MD2WECHAT_WRITERS_DIR"

// StyleManager 风格管理器
type StyleManager struct {
	styles      map[string]*WriterStyle
	writersDir  string
	initialized bool
}

// NewStyleManager 创建风格管理器
func NewStyleManager() *StyleManager {
	return &StyleManager{
		styles: make(map[string]*WriterStyle),
	}
}

// LoadStyles 加载所有风格配置
func (sm *StyleManager) LoadStyles() error {
	sm.styles = make(map[string]*WriterStyle)

	if err := sm.loadBuiltinStyles(); err != nil {
		return fmt.Errorf("加载内置 writers: %w", err)
	}

	for _, dir := range sm.getWriterDirs() {
		if err := sm.loadStylesFromDir(dir); err != nil {
			return err
		}
	}

	sm.initialized = true
	return nil
}

// loadStyle 加载单个风格文件
func (sm *StyleManager) loadStyle(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件: %w", err)
	}

	return sm.loadStyleData(data)
}

func (sm *StyleManager) loadStyleData(data []byte) error {
	var style WriterStyle
	if err := yaml.Unmarshal(data, &style); err != nil {
		return fmt.Errorf("解析 YAML: %w", err)
	}

	// 验证必需字段
	if style.EnglishName == "" {
		return fmt.Errorf("缺少必需字段: english_name")
	}

	if style.Name == "" {
		style.Name = style.EnglishName
	}

	// 设置默认值
	if style.Category == "" {
		style.Category = "自定义"
	}
	if style.Version == "" {
		style.Version = "1.0"
	}

	sm.styles[style.EnglishName] = &style
	return nil
}

func (sm *StyleManager) loadBuiltinStyles() error {
	names, err := assets.ListBuiltinWriters()
	if err != nil {
		return err
	}

	for _, name := range names {
		data, err := assets.ReadBuiltinWriter(name)
		if err != nil {
			return fmt.Errorf("读取内置 writer %s: %w", name, err)
		}
		if err := sm.loadStyleData(data); err != nil {
			return fmt.Errorf("加载内置 writer %s: %w", name, err)
		}
	}

	return nil
}

func (sm *StyleManager) loadStylesFromDir(dir string) error {
	if dir == "" {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取 writers 目录 %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		stylePath := filepath.Join(dir, name)
		if err := sm.loadStyle(stylePath); err != nil {
			// 记录错误但继续加载其他风格
			continue
		}
	}

	return nil
}

// getWriterDirs 获取 writers 目录路径，按优先级从低到高返回，后加载者覆盖前加载者。
func (sm *StyleManager) getWriterDirs() []string {
	homeDir, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(homeDir, ".md2wechat-writers"),
		filepath.Join(homeDir, ".config", "md2wechat", "writers"),
		"writers",
	}

	if explicitDir := strings.TrimSpace(os.Getenv(writersDirEnvVar)); explicitDir != "" {
		dirs = append(dirs, explicitDir)
	}

	return dirs
}

// GetStyle 获取风格
func (sm *StyleManager) GetStyle(name string) (*WriterStyle, error) {
	// 确保已加载
	if !sm.initialized {
		if err := sm.LoadStyles(); err != nil {
			return nil, err
		}
	}

	// 支持名称映射（中文转英文）
	englishName := sm.mapToEnglishName(name)

	style, ok := sm.styles[englishName]
	if !ok {
		return nil, NewStyleNotFoundError(name)
	}

	return style, nil
}

// mapToEnglishName 将各种名称映射到英文标识
func (sm *StyleManager) mapToEnglishName(name string) string {
	// 名称映射表
	nameMap := map[string]string{
		// Dan Koe 别名
		"dan-koe": "dan-koe",
		"dankoe":  "dan-koe",
		"dan":     "dan-koe",
		"koe":     "dan-koe",
		// 未来可以添加更多映射
	}

	// 先检查精确匹配
	if english, ok := nameMap[strings.ToLower(name)]; ok {
		return english
	}

	// 检查是否已经是英文名
	if _, ok := sm.styles[name]; ok {
		return name
	}

	// 检查中文名称映射
	for _, style := range sm.styles {
		if style.Name == name {
			return style.EnglishName
		}
	}

	return name
}

// ListStyles 列出所有风格
func (sm *StyleManager) ListStyles() []StyleSummary {
	// 确保已加载
	if !sm.initialized {
		_ = sm.LoadStyles()
	}

	result := make([]StyleSummary, 0, len(sm.styles))
	for _, style := range sm.styles {
		result = append(result, StyleSummary{
			Name:        style.Name,
			EnglishName: style.EnglishName,
			Category:    style.Category,
			Description: style.Description,
			CoverStyle:  style.CoverStyle,
		})
	}
	return result
}

// ListStyleNames 列出所有风格的英文名称
func (sm *StyleManager) ListStyleNames() []string {
	// 确保已加载
	if !sm.initialized {
		_ = sm.LoadStyles()
	}

	names := make([]string, 0, len(sm.styles))
	for name := range sm.styles {
		names = append(names, name)
	}
	return names
}

// GetWritersDir 获取 writers 目录路径（公开方法）
func (sm *StyleManager) GetWritersDir() string {
	if sm.writersDir == "" {
		if explicitDir := strings.TrimSpace(os.Getenv(writersDirEnvVar)); explicitDir != "" {
			sm.writersDir = explicitDir
			return sm.writersDir
		}

		// 优先返回最接近用户可编辑的目录
		paths := []string{
			"writers",
			filepath.Join(os.Getenv("HOME"), ".config", "md2wechat", "writers"),
			filepath.Join(os.Getenv("HOME"), ".md2wechat-writers"),
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				sm.writersDir = path
				return sm.writersDir
			}
		}

		sm.writersDir = "writers"
	}
	return sm.writersDir
}

// ReloadStyles 重新加载所有风格
func (sm *StyleManager) ReloadStyles() error {
	sm.styles = make(map[string]*WriterStyle)
	sm.initialized = false
	return sm.LoadStyles()
}

// GetStyleCount 获取风格数量
func (sm *StyleManager) GetStyleCount() int {
	if !sm.initialized {
		_ = sm.LoadStyles()
	}
	return len(sm.styles)
}

// HasStyle 检查是否存在指定风格
func (sm *StyleManager) HasStyle(name string) bool {
	_, err := sm.GetStyle(name)
	return err == nil
}

// GetDefaultStyle 获取默认风格
func (sm *StyleManager) GetDefaultStyle() (*WriterStyle, error) {
	return sm.GetStyle(DefaultStyleName)
}

// ValidateStyle 验证风格配置
func (sm *StyleManager) ValidateStyle(style *WriterStyle) error {
	if style.EnglishName == "" {
		return fmt.Errorf("english_name 是必需的")
	}
	if style.WritingPrompt == "" {
		return fmt.Errorf("writing_prompt 是必需的")
	}
	return nil
}

// CreateStyleDirectory 创建 writers 目录
func (sm *StyleManager) CreateStyleDirectory() error {
	dir := sm.GetWritersDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录: %w", err)
	}
	return nil
}

// ExportStyle 导出风格配置到文件
func (sm *StyleManager) ExportStyle(style *WriterStyle, destPath string) error {
	if err := sm.ValidateStyle(style); err != nil {
		return err
	}

	data, err := yaml.Marshal(style)
	if err != nil {
		return fmt.Errorf("序列化 YAML: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("写入文件: %w", err)
	}

	return nil
}

// GetStyleByCategory 按分类获取风格
func (sm *StyleManager) GetStyleByCategory(category string) []StyleSummary {
	// 确保已加载
	if !sm.initialized {
		_ = sm.LoadStyles()
	}

	result := []StyleSummary{}
	for _, style := range sm.styles {
		if style.Category == category {
			result = append(result, StyleSummary{
				Name:        style.Name,
				EnglishName: style.EnglishName,
				Category:    style.Category,
				Description: style.Description,
				CoverStyle:  style.CoverStyle,
			})
		}
	}
	return result
}

// ListCategories 列出所有分类
func (sm *StyleManager) ListCategories() []string {
	// 确保已加载
	if !sm.initialized {
		_ = sm.LoadStyles()
	}

	categoryMap := make(map[string]bool)
	for _, style := range sm.styles {
		if style.Category != "" {
			categoryMap[style.Category] = true
		}
	}

	categories := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	return categories
}

// GetStyleWithPrompt 获取风格并填充提示词模板
func (sm *StyleManager) GetStyleWithPrompt(name string, templateVars map[string]string) (*WriterStyle, error) {
	style, err := sm.GetStyle(name)
	if err != nil {
		return nil, err
	}

	// 替换提示词中的变量
	prompt := style.WritingPrompt
	for key, value := range templateVars {
		placeholder := "{" + key + "}"
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}

	// 创建副本，不修改原风格
	styleCopy := *style
	styleCopy.WritingPrompt = prompt

	return &styleCopy, nil
}
