package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func chdirTemp(t *testing.T) string {
	t.Helper()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
	})
	return tmp
}

func TestLoadStyleAppliesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	stylePath := filepath.Join(tmpDir, "style.yaml")
	content := strings.Join([]string{
		`english_name: custom-style`,
		`writing_prompt: Write with intent.`,
	}, "\n")
	if err := os.WriteFile(stylePath, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	sm := NewStyleManager()
	if err := sm.loadStyle(stylePath); err != nil {
		t.Fatalf("loadStyle() error = %v", err)
	}

	style, ok := sm.styles["custom-style"]
	if !ok {
		t.Fatal("expected custom-style to be loaded")
	}
	if style.Name != "custom-style" {
		t.Fatalf("Name = %q", style.Name)
	}
	if style.Category != "自定义" {
		t.Fatalf("Category = %q", style.Category)
	}
	if style.Version != "1.0" {
		t.Fatalf("Version = %q", style.Version)
	}
}

func TestGetStyleWithPromptInterpolatesTemplateVariables(t *testing.T) {
	sm := &StyleManager{
		styles: map[string]*WriterStyle{
			DefaultStyleName: {
				Name:          "Dan Koe",
				EnglishName:   DefaultStyleName,
				WritingPrompt: "Title: {title}\nBody: {body}",
			},
		},
		initialized: true,
	}

	style, err := sm.GetStyleWithPrompt(DefaultStyleName, map[string]string{
		"title": "Test Title",
		"body":  "Test Body",
	})
	if err != nil {
		t.Fatalf("GetStyleWithPrompt() error = %v", err)
	}
	if !strings.Contains(style.WritingPrompt, "Test Title") || !strings.Contains(style.WritingPrompt, "Test Body") {
		t.Fatalf("WritingPrompt = %q", style.WritingPrompt)
	}
	if sm.styles[DefaultStyleName].WritingPrompt != "Title: {title}\nBody: {body}" {
		t.Fatalf("original style was mutated: %q", sm.styles[DefaultStyleName].WritingPrompt)
	}
}

func TestValidateStyleRequiresCoreFields(t *testing.T) {
	sm := NewStyleManager()

	if err := sm.ValidateStyle(&WriterStyle{WritingPrompt: "prompt"}); err == nil {
		t.Fatal("expected english_name validation error")
	}
	if err := sm.ValidateStyle(&WriterStyle{EnglishName: "custom"}); err == nil {
		t.Fatal("expected writing_prompt validation error")
	}
}

func TestGetWritersDirPrefersExplicitEnvironmentVariable(t *testing.T) {
	customDir := filepath.Join(t.TempDir(), "custom-writers")
	t.Setenv(writersDirEnvVar, customDir)

	sm := NewStyleManager()
	if got := sm.GetWritersDir(); got != customDir {
		t.Fatalf("GetWritersDir() = %q, want %q", got, customDir)
	}
}

func TestLoadStylesFallsBackToBuiltinDanKoe(t *testing.T) {
	chdirTemp(t)
	t.Setenv(writersDirEnvVar, "")
	t.Setenv("HOME", filepath.Join(t.TempDir(), "home"))

	sm := NewStyleManager()
	style, err := sm.GetDefaultStyle()
	if err != nil {
		t.Fatalf("GetDefaultStyle() error = %v", err)
	}
	if style.EnglishName != DefaultStyleName {
		t.Fatalf("EnglishName = %q, want %q", style.EnglishName, DefaultStyleName)
	}
	if style.Name != "Dan Koe" {
		t.Fatalf("Name = %q, want Dan Koe", style.Name)
	}
}

func TestLoadStylesRespectsPriorityOrder(t *testing.T) {
	tmp := chdirTemp(t)
	homeDir := filepath.Join(tmp, "home")
	configDir := filepath.Join(homeDir, ".config", "md2wechat", "writers")
	legacyDir := filepath.Join(homeDir, ".md2wechat-writers")
	cwdDir := filepath.Join(tmp, "writers")
	envDir := filepath.Join(tmp, "env-writers")

	for _, dir := range []string{configDir, legacyDir, cwdDir, envDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("MkdirAll(%s) error = %v", dir, err)
		}
	}

	t.Setenv("HOME", homeDir)
	t.Setenv(writersDirEnvVar, envDir)

	writeStyle := func(dir, description string) {
		t.Helper()
		content := strings.Join([]string{
			`english_name: dan-koe`,
			`name: ` + description,
			`description: ` + description,
			`writing_prompt: Prompt for ` + description,
		}, "\n")
		if err := os.WriteFile(filepath.Join(dir, "dan-koe.yaml"), []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", dir, err)
		}
	}

	writeStyle(legacyDir, "legacy-home")
	writeStyle(configDir, "config-home")
	writeStyle(cwdDir, "cwd")
	writeStyle(envDir, "env")

	sm := NewStyleManager()
	style, err := sm.GetDefaultStyle()
	if err != nil {
		t.Fatalf("GetDefaultStyle() error = %v", err)
	}
	if style.Description != "env" {
		t.Fatalf("Description = %q, want env", style.Description)
	}
	if style.WritingPrompt != "Prompt for env" {
		t.Fatalf("WritingPrompt = %q, want Prompt for env", style.WritingPrompt)
	}
}

func TestLoadStylesFallsBackThroughUserDirsBeforeBuiltin(t *testing.T) {
	tmp := chdirTemp(t)
	homeDir := filepath.Join(tmp, "home")
	configDir := filepath.Join(homeDir, ".config", "md2wechat", "writers")
	legacyDir := filepath.Join(homeDir, ".md2wechat-writers")

	for _, dir := range []string{configDir, legacyDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("MkdirAll(%s) error = %v", dir, err)
		}
	}

	t.Setenv("HOME", homeDir)
	t.Setenv(writersDirEnvVar, "")

	writeStyle := func(dir, description string) {
		t.Helper()
		content := strings.Join([]string{
			`english_name: dan-koe`,
			`name: ` + description,
			`description: ` + description,
			`writing_prompt: Prompt for ` + description,
		}, "\n")
		if err := os.WriteFile(filepath.Join(dir, "dan-koe.yaml"), []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", dir, err)
		}
	}

	writeStyle(legacyDir, "legacy-home")
	writeStyle(configDir, "config-home")

	sm := NewStyleManager()
	style, err := sm.GetDefaultStyle()
	if err != nil {
		t.Fatalf("GetDefaultStyle() error = %v", err)
	}
	if style.Description != "config-home" {
		t.Fatalf("Description = %q, want config-home", style.Description)
	}
}
