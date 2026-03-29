package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThemeManagerLoadsBuiltinThemesWithoutFilesystemAssets(t *testing.T) {
	t.Chdir(t.TempDir())

	tm := NewThemeManager()
	theme, err := tm.GetTheme("default")
	if err != nil {
		t.Fatalf("GetTheme(default) error = %v", err)
	}
	if theme.Name != "default" || theme.Type != "api" {
		t.Fatalf("unexpected builtin theme: %#v", theme)
	}

	aiTheme, err := tm.GetTheme("autumn-warm")
	if err != nil {
		t.Fatalf("GetTheme(autumn-warm) error = %v", err)
	}
	if aiTheme.Type != "ai" {
		t.Fatalf("expected AI builtin theme, got %#v", aiTheme)
	}
}

func TestThemeManagerPrefersCurrentDirectoryThemeOverBuiltin(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	if err := os.MkdirAll("themes", 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	override := []byte("name: default\ntype: api\ndescription: local override\napi_theme: default\n")
	if err := os.WriteFile(filepath.Join("themes", "default.yaml"), override, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tm := NewThemeManager()
	theme, err := tm.GetTheme("default")
	if err != nil {
		t.Fatalf("GetTheme(default) error = %v", err)
	}
	if theme.Description != "local override" {
		t.Fatalf("Description = %q, want local override", theme.Description)
	}
}

func TestThemeManagerPrefersExplicitThemesDirOverBuiltin(t *testing.T) {
	customDir := filepath.Join(t.TempDir(), "custom-themes")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	override := []byte("name: default\ntype: api\ndescription: env override\napi_theme: default\n")
	if err := os.WriteFile(filepath.Join(customDir, "default.yaml"), override, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv(themesDirEnvVar, customDir)
	t.Chdir(t.TempDir())

	tm := NewThemeManager()
	theme, err := tm.GetTheme("default")
	if err != nil {
		t.Fatalf("GetTheme(default) error = %v", err)
	}
	if theme.Description != "env override" {
		t.Fatalf("Description = %q, want env override", theme.Description)
	}
}
