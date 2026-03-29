package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromptBuilderAddBuildAndExport(t *testing.T) {
	builder := NewPromptBuilder()
	err := builder.AddTemplate(&PromptTemplate{
		Name:        "article",
		Description: "Article template",
		Template:    "Title: {{TITLE}}\nBody: {{MARKDOWN}}\nFont: {{FONT_SIZE}}",
	})
	if err != nil {
		t.Fatalf("AddTemplate() error = %v", err)
	}

	prompt, err := builder.BuildPrompt("article", map[string]string{
		"TITLE":    "Hello",
		"MARKDOWN": "World",
	})
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "Hello") || !strings.Contains(prompt, "World") || !strings.Contains(prompt, "16px") {
		t.Fatalf("prompt = %q", prompt)
	}

	exported, err := builder.ExportPrompt("article", map[string]string{
		"TITLE":    "Hello",
		"MARKDOWN": "World",
	}, &ExportOptions{Format: "markdown", IncludeHeader: true})
	if err != nil {
		t.Fatalf("ExportPrompt() error = %v", err)
	}
	if !strings.Contains(exported, "# article") || !strings.Contains(exported, "Article template") {
		t.Fatalf("exported = %q", exported)
	}

	if _, err := builder.GetTemplate("article"); err != nil {
		t.Fatalf("GetTemplate() error = %v", err)
	}
	if _, err := builder.GetVariable("{{TITLE}}"); err != nil {
		t.Fatalf("GetVariable() error = %v", err)
	}
	if len(builder.ListTemplates()) != 1 || len(builder.ListVariables()) == 0 {
		t.Fatalf("unexpected lists: templates=%v variables=%v", builder.ListTemplates(), builder.ListVariables())
	}
}

func TestPromptBuilderErrorPaths(t *testing.T) {
	builder := NewPromptBuilder()

	if err := builder.AddTemplate(&PromptTemplate{}); err == nil {
		t.Fatal("expected empty template name error")
	}
	if _, err := builder.BuildPrompt("missing", nil); err == nil {
		t.Fatal("expected missing template error")
	}
	if _, err := builder.GetTemplate("missing"); err == nil {
		t.Fatal("expected missing template error")
	}
	if _, err := builder.GetVariable("missing"); err == nil {
		t.Fatal("expected missing variable error")
	}
	if err := builder.ValidateTemplate("missing"); err == nil {
		t.Fatal("expected validate missing template error")
	}
}

func TestThemeManagerInMemoryLookups(t *testing.T) {
	tm := NewThemeManager()
	tm.themes["api-theme"] = Theme{
		Name:        "api-theme",
		Type:        "api",
		Description: "API theme",
		APITheme:    "default",
		Colors:      map[string]string{"primary": "#000"},
	}
	tm.themes["ai-theme"] = Theme{
		Name:        "ai-theme",
		Type:        "ai",
		Description: "AI theme",
		Prompt:      "Prompt body",
	}

	if len(tm.ListThemes()) != 2 || len(tm.ListAIThemes()) != 1 || len(tm.ListAPIThemes()) != 1 {
		t.Fatalf("unexpected theme lists")
	}
	if apiTheme, err := tm.GetAPITheme("api-theme"); err != nil || apiTheme != "default" {
		t.Fatalf("GetAPITheme() = %q, %v", apiTheme, err)
	}
	if aiPrompt, err := tm.GetAIPrompt("ai-theme"); err != nil || aiPrompt != "Prompt body" {
		t.Fatalf("GetAIPrompt() = %q, %v", aiPrompt, err)
	}
	if tm.GetThemeDescription("missing") != "未知主题" {
		t.Fatalf("unexpected missing theme description")
	}
	if !tm.IsAPITheme("api-theme") || !tm.IsAITheme("ai-theme") {
		t.Fatalf("theme type checks failed")
	}
	colors, err := tm.GetThemeColors("api-theme")
	if err != nil || colors["primary"] != "#000" {
		t.Fatalf("GetThemeColors() = %#v, %v", colors, err)
	}
	if _, err := tm.GetAPITheme("ai-theme"); err == nil {
		t.Fatal("expected api type mismatch error")
	}
	if _, err := tm.GetAIPrompt("api-theme"); err == nil {
		t.Fatal("expected ai type mismatch error")
	}
	if _, err := tm.GetThemeInfo("ai-theme"); err != nil {
		t.Fatalf("GetThemeInfo() error = %v", err)
	}
	if err := tm.EnsureLoaded(); err != nil {
		t.Fatalf("EnsureLoaded() error = %v", err)
	}
}

func TestThemeManagerLoadThemeAppliesDefaults(t *testing.T) {
	tm := NewThemeManager()
	path := filepath.Join(t.TempDir(), "custom.yaml")
	content := []byte("name: custom-theme\nprompt: hello\n")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write theme file: %v", err)
	}

	if err := tm.LoadTheme(path); err != nil {
		t.Fatalf("LoadTheme() error = %v", err)
	}

	theme, err := tm.GetTheme("custom-theme")
	if err != nil {
		t.Fatalf("GetTheme() error = %v", err)
	}
	if theme.Type != "ai" || theme.Description != "custom-theme" {
		t.Fatalf("theme defaults not applied: %#v", theme)
	}
}

func TestBuildCustomAIPromptAndEstimateTokenCount(t *testing.T) {
	prompt := BuildCustomAIPrompt("请帮我排版")
	if !strings.Contains(prompt, "重要规则") || !strings.Contains(prompt, "请转换以下 Markdown内容：") {
		t.Fatalf("prompt = %q", prompt)
	}

	if count := EstimateTokenCount("中文abcde"); count <= 0 {
		t.Fatalf("EstimateTokenCount() = %d", count)
	}
}
