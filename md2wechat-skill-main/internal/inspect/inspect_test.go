package inspect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

func TestRunResolvesMetadataAndDetectsDuplicateH1(t *testing.T) {
	dir := t.TempDir()
	markdownPath := filepath.Join(dir, "article.md")
	markdown := strings.Join([]string{
		"---",
		"title: Frontmatter 标题",
		"author: 张三",
		"summary: Frontmatter 摘要",
		"---",
		"",
		"# Frontmatter 标题",
		"",
		"正文",
	}, "\n")

	result, err := Run(&Input{
		MarkdownFile: markdownPath,
		Markdown:     markdown,
		Mode:         "api",
		Theme:        "default",
		FontSize:     "medium",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result.Metadata.Title.Value != "Frontmatter 标题" || result.Metadata.Title.Source != "frontmatter.title" {
		t.Fatalf("title = %#v", result.Metadata.Title)
	}
	if result.Metadata.Author.Value != "张三" || result.Metadata.Author.Source != "frontmatter.author" {
		t.Fatalf("author = %#v", result.Metadata.Author)
	}
	if result.Metadata.Digest.Value != "Frontmatter 摘要" || result.Metadata.Digest.Source != "frontmatter.summary" {
		t.Fatalf("digest = %#v", result.Metadata.Digest)
	}
	if !result.Structure.BodyH1.Present || result.Structure.BodyH1.Text != "Frontmatter 标题" {
		t.Fatalf("body_h1 = %#v", result.Structure.BodyH1)
	}
	if !result.Structure.DuplicateTitleRisk {
		t.Fatal("expected duplicate title risk")
	}
}

func TestRunDoesNotUseFirstBodyLineAsTitleAndFlagsMissingRequirements(t *testing.T) {
	dir := t.TempDir()
	markdownPath := filepath.Join(dir, "article.md")
	localImage := filepath.Join("images", "missing.png")
	markdown := strings.Join([]string{
		"正文第一行不是标题",
		"",
		"![missing](images/missing.png)",
	}, "\n")

	result, err := Run(&Input{
		MarkdownFile:    markdownPath,
		Markdown:        markdown,
		Mode:            "api",
		Theme:           "default",
		FontSize:        "medium",
		DraftRequested:  true,
		UploadRequested: true,
		Config:          &config.Config{},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result.Metadata.Title.Value != "未命名文章" || result.Metadata.Title.Source != "fallback.untitled" {
		t.Fatalf("title = %#v", result.Metadata.Title)
	}
	if result.Readiness.ConvertReady {
		t.Fatal("expected convert_ready false without API key")
	}
	if result.Readiness.UploadReady {
		t.Fatal("expected upload_ready false without WeChat config")
	}
	if result.Readiness.DraftReady {
		t.Fatal("expected draft_ready false without cover and config")
	}

	foundLocalImageMissing := false
	foundMissingAPIKey := false
	foundMissingCover := false
	foundMissingWeChatConfig := false
	for _, check := range result.Checks {
		switch check.Code {
		case "LOCAL_IMAGE_MISSING":
			foundLocalImageMissing = strings.Contains(check.Message, localImage)
		case "MISSING_API_KEY":
			foundMissingAPIKey = true
		case "MISSING_COVER":
			foundMissingCover = true
		case "MISSING_WECHAT_CONFIG":
			foundMissingWeChatConfig = true
		}
	}
	if !foundLocalImageMissing || !foundMissingAPIKey || !foundMissingCover || !foundMissingWeChatConfig {
		t.Fatalf("checks = %#v", result.Checks)
	}
}

func TestRunMarksLocalImageExists(t *testing.T) {
	dir := t.TempDir()
	imageDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	imagePath := filepath.Join(imageDir, "ok.png")
	if err := os.WriteFile(imagePath, []byte("png"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := Run(&Input{
		MarkdownFile: filepath.Join(dir, "article.md"),
		Markdown:     "![ok](images/ok.png)\n",
		Mode:         "ai",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(result.Assets) != 1 || !result.Assets[0].Exists {
		t.Fatalf("assets = %#v", result.Assets)
	}
	if result.Readiness.PreviewFidelity != PreviewFidelityDegraded {
		t.Fatalf("preview_fidelity = %q", result.Readiness.PreviewFidelity)
	}
}

func TestRunMetadataSourcePrecedence(t *testing.T) {
	t.Run("cli overrides frontmatter", func(t *testing.T) {
		dir := t.TempDir()
		result, err := Run(&Input{
			MarkdownFile: filepath.Join(dir, "article.md"),
			Markdown: strings.Join([]string{
				"---",
				"title: Frontmatter 标题",
				"author: Frontmatter 作者",
				"digest: Frontmatter 摘要",
				"---",
				"",
				"# 正文标题",
			}, "\n"),
			TitleOverride:  "CLI 标题",
			AuthorOverride: "CLI 作者",
			DigestOverride: "CLI 摘要",
		})
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

		if result.Metadata.Title.Source != "cli.title" || result.Metadata.Title.Value != "CLI 标题" {
			t.Fatalf("title = %#v", result.Metadata.Title)
		}
		if result.Metadata.Author.Source != "cli.author" || result.Metadata.Author.Value != "CLI 作者" {
			t.Fatalf("author = %#v", result.Metadata.Author)
		}
		if result.Metadata.Digest.Source != "cli.digest" || result.Metadata.Digest.Value != "CLI 摘要" {
			t.Fatalf("digest = %#v", result.Metadata.Digest)
		}
	})

	t.Run("empty frontmatter title falls back to markdown heading", func(t *testing.T) {
		dir := t.TempDir()
		result, err := Run(&Input{
			MarkdownFile: filepath.Join(dir, "article.md"),
			Markdown: strings.Join([]string{
				"---",
				"title: \"\"",
				"author: 张三",
				"---",
				"",
				"# Heading 标题",
				"",
				"正文",
			}, "\n"),
		})
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

		if result.Metadata.Title.Source != "markdown.heading" || result.Metadata.Title.Value != "Heading 标题" {
			t.Fatalf("title = %#v", result.Metadata.Title)
		}
		if result.Metadata.Author.Source != "frontmatter.author" {
			t.Fatalf("author = %#v", result.Metadata.Author)
		}
	})
}

func TestRunDigestSourcePrecedence(t *testing.T) {
	cases := []struct {
		name        string
		frontmatter []string
		wantValue   string
		wantSource  string
	}{
		{
			name: "digest wins",
			frontmatter: []string{
				"digest: Digest 文本",
				"summary: Summary 文本",
				"description: Description 文本",
			},
			wantValue:  "Digest 文本",
			wantSource: "frontmatter.digest",
		},
		{
			name: "summary wins when digest missing",
			frontmatter: []string{
				"summary: Summary 文本",
				"description: Description 文本",
			},
			wantValue:  "Summary 文本",
			wantSource: "frontmatter.summary",
		},
		{
			name: "description used last",
			frontmatter: []string{
				"description: Description 文本",
			},
			wantValue:  "Description 文本",
			wantSource: "frontmatter.description",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			markdown := strings.Join(append([]string{"---"}, append(tc.frontmatter, "---", "", "# 标题")...), "\n")
			result, err := Run(&Input{
				MarkdownFile: filepath.Join(dir, "article.md"),
				Markdown:     markdown,
			})
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if result.Metadata.Digest.Value != tc.wantValue || result.Metadata.Digest.Source != tc.wantSource {
				t.Fatalf("digest = %#v", result.Metadata.Digest)
			}
		})
	}
}

func TestRunReadinessMatrix(t *testing.T) {
	fullCfg := &config.Config{
		MD2WechatAPIKey: "api-key",
		WechatAppID:     "appid",
		WechatSecret:    "secret",
	}
	apiOnlyCfg := &config.Config{MD2WechatAPIKey: "api-key"}
	validCoverPath := filepath.Join(t.TempDir(), "cover.jpg")
	if err := os.WriteFile(validCoverPath, []byte("cover"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cases := []struct {
		name         string
		input        Input
		wantConvert  bool
		wantUpload   bool
		wantDraft    bool
		wantPreview  string
		wantChecks   []string
		absentChecks []string
	}{
		{
			name: "api without config or cover",
			input: Input{
				MarkdownFile:    filepath.Join(t.TempDir(), "article.md"),
				Markdown:        "# 标题\n",
				Mode:            "api",
				DraftRequested:  true,
				UploadRequested: true,
				Config:          &config.Config{},
			},
			wantConvert: false,
			wantUpload:  false,
			wantDraft:   false,
			wantPreview: PreviewFidelityDegraded,
			wantChecks:  []string{"MISSING_API_KEY", "MISSING_WECHAT_CONFIG", "MISSING_COVER"},
		},
		{
			name: "ai without config or cover",
			input: Input{
				MarkdownFile:    filepath.Join(t.TempDir(), "article.md"),
				Markdown:        "# 标题\n",
				Mode:            "ai",
				DraftRequested:  true,
				UploadRequested: true,
				Config:          &config.Config{},
			},
			wantConvert:  true,
			wantUpload:   false,
			wantDraft:    false,
			wantPreview:  PreviewFidelityDegraded,
			wantChecks:   []string{"AI_MODE_ACTION_REQUIRED", "MISSING_WECHAT_CONFIG", "MISSING_COVER"},
			absentChecks: []string{"MISSING_API_KEY"},
		},
		{
			name: "api with api key only",
			input: Input{
				MarkdownFile:   filepath.Join(t.TempDir(), "article.md"),
				Markdown:       "# 标题\n",
				Mode:           "api",
				DraftRequested: true,
				Config:         apiOnlyCfg,
			},
			wantConvert:  true,
			wantUpload:   false,
			wantDraft:    false,
			wantPreview:  PreviewFidelityExact,
			wantChecks:   []string{"MISSING_WECHAT_CONFIG", "MISSING_COVER"},
			absentChecks: []string{"MISSING_API_KEY"},
		},
		{
			name: "api fully ready",
			input: Input{
				MarkdownFile:    filepath.Join(t.TempDir(), "article.md"),
				Markdown:        "# 标题\n",
				Mode:            "api",
				DraftRequested:  true,
				UploadRequested: true,
				CoverImagePath:  validCoverPath,
				Config:          fullCfg,
			},
			wantConvert:  true,
			wantUpload:   true,
			wantDraft:    true,
			wantPreview:  PreviewFidelityExact,
			absentChecks: []string{"MISSING_API_KEY", "MISSING_WECHAT_CONFIG", "MISSING_COVER"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Run(&tc.input)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if result.Readiness.ConvertReady != tc.wantConvert {
				t.Fatalf("convert_ready = %t", result.Readiness.ConvertReady)
			}
			if result.Readiness.UploadReady != tc.wantUpload {
				t.Fatalf("upload_ready = %t", result.Readiness.UploadReady)
			}
			if result.Readiness.DraftReady != tc.wantDraft {
				t.Fatalf("draft_ready = %t", result.Readiness.DraftReady)
			}
			if result.Readiness.PreviewFidelity != tc.wantPreview {
				t.Fatalf("preview_fidelity = %q", result.Readiness.PreviewFidelity)
			}
			for _, code := range tc.wantChecks {
				if !hasCheckCode(result.Checks, code) {
					t.Fatalf("missing check %s in %#v", code, result.Checks)
				}
			}
			for _, code := range tc.absentChecks {
				if hasCheckCode(result.Checks, code) {
					t.Fatalf("unexpected check %s in %#v", code, result.Checks)
				}
			}
		})
	}
}

func TestRunMakesReadinessFalseForBlockingChecks(t *testing.T) {
	fullCfg := &config.Config{
		MD2WechatAPIKey: "api-key",
		WechatAppID:     "appid",
		WechatSecret:    "secret",
	}

	t.Run("metadata limit errors block convert and downstream readiness", func(t *testing.T) {
		result, err := Run(&Input{
			MarkdownFile:   filepath.Join(t.TempDir(), "article.md"),
			Markdown:       "# 标题\n",
			Mode:           "api",
			TitleOverride:  strings.Repeat("标", 33),
			DraftRequested: true,
			CoverImagePath: filepath.Join(t.TempDir(), "cover.png"),
			Config:         fullCfg,
		})
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		if !hasCheckCode(result.Checks, "TITLE_TOO_LONG") {
			t.Fatalf("checks = %#v", result.Checks)
		}
		if result.Readiness.ConvertReady {
			t.Fatal("expected convert_ready false when title is too long")
		}
		if result.Readiness.UploadReady {
			t.Fatal("expected upload_ready false when convert is blocked")
		}
		if result.Readiness.DraftReady {
			t.Fatal("expected draft_ready false when convert is blocked")
		}
	})

	t.Run("cover path errors block draft readiness", func(t *testing.T) {
		result, err := Run(&Input{
			MarkdownFile:   filepath.Join(t.TempDir(), "article.md"),
			Markdown:       "# 标题\n",
			Mode:           "api",
			DraftRequested: true,
			CoverImagePath: filepath.Join(t.TempDir(), "missing-cover.png"),
			Config:         fullCfg,
		})
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		if !hasCheckCode(result.Checks, "COVER_IMAGE_MISSING") {
			t.Fatalf("checks = %#v", result.Checks)
		}
		if !result.Readiness.ConvertReady {
			t.Fatal("expected convert_ready true when only the cover path is invalid")
		}
		if !result.Readiness.UploadReady {
			t.Fatal("expected upload_ready true when config is valid and no upload blocker exists")
		}
		if result.Readiness.DraftReady {
			t.Fatal("expected draft_ready false when the cover image path is invalid")
		}
	})
}

func TestRunRejectsInvalidMode(t *testing.T) {
	_, err := Run(&Input{
		MarkdownFile: filepath.Join(t.TempDir(), "article.md"),
		Markdown:     "# 标题\n",
		Mode:         "foo",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid convert mode") {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRunInspectsAbsoluteRemoteAndAIAssets(t *testing.T) {
	dir := t.TempDir()
	absPath := filepath.Join(dir, "cover.png")
	if err := os.WriteFile(absPath, []byte("png"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	markdown := strings.Join([]string{
		"# 标题",
		"",
		"![local](" + absPath + ")",
		"![remote](https://example.com/a.png)",
		"![ai](__generate:draw a cover__)",
	}, "\n")

	result, err := Run(&Input{
		MarkdownFile: filepath.Join(dir, "article.md"),
		Markdown:     markdown,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Structure.Images.Total != 3 || result.Structure.Images.Local != 1 || result.Structure.Images.Remote != 1 || result.Structure.Images.AI != 1 {
		t.Fatalf("image counters = %#v", result.Structure.Images)
	}
	if len(result.Assets) != 3 {
		t.Fatalf("assets = %#v", result.Assets)
	}
	if result.Assets[0].ResolvedSource != absPath || !result.Assets[0].Exists {
		t.Fatalf("local asset = %#v", result.Assets[0])
	}
	if result.Assets[1].Kind != "online" || !result.Assets[1].Exists {
		t.Fatalf("remote asset = %#v", result.Assets[1])
	}
	if result.Assets[2].Kind != "ai" || !result.Assets[2].Exists {
		t.Fatalf("ai asset = %#v", result.Assets[2])
	}
}

func TestRunAddsMetadataAndImageIntentChecks(t *testing.T) {
	dir := t.TempDir()
	imagePath := filepath.Join(dir, "cover.png")
	if err := os.WriteFile(imagePath, []byte("png"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	markdown := strings.Join([]string{
		"---",
		"title: 草稿标题",
		"summary: 这是摘要",
		"---",
		"",
		"# 正文一级标题",
		"",
		"![local](cover.png)",
	}, "\n")

	result, err := Run(&Input{
		MarkdownFile: filepath.Join(dir, "article.md"),
		Markdown:     markdown,
		Mode:         "ai",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, code := range []string{
		"TITLE_BODY_MISMATCH",
		"DIGEST_METADATA_ONLY",
		"IMAGE_REPLACEMENT_REQUIRES_UPLOAD_OR_DRAFT",
		"AI_MODE_ACTION_REQUIRED",
	} {
		if !hasCheckCode(result.Checks, code) {
			t.Fatalf("missing check %s in %#v", code, result.Checks)
		}
	}
	if hasCheckCode(result.Checks, "DUPLICATE_H1") {
		t.Fatalf("unexpected duplicate h1 check in %#v", result.Checks)
	}
}

func hasCheckCode(checks []Check, code string) bool {
	for _, check := range checks {
		if check.Code == code {
			return true
		}
	}
	return false
}
