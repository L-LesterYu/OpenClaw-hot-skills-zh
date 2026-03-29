package preview

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/inspect"
)

func TestRenderExactHTMLIncludesResolvedMetadataAndChecks(t *testing.T) {
	html, err := Render(PageData{
		Title:      "预览标题",
		SourceFile: "/tmp/article.md",
		Context: inspect.Context{
			Mode:           "api",
			Theme:          "default",
			FontSize:       "medium",
			BackgroundType: "none",
		},
		Metadata: inspect.MetadataState{
			Title:  inspect.MetadataField{Value: "预览标题", Source: "frontmatter.title", Length: 4, Limit: 32},
			Author: inspect.MetadataField{Value: "作者", Source: "cli.author", Length: 2, Limit: 16},
			Digest: inspect.MetadataField{Value: "摘要", Source: "frontmatter.summary", Length: 2, Limit: 128},
		},
		Structure: inspect.Structure{
			BodyH1: inspect.HeadingInfo{Present: true, Text: "预览标题"},
			Images: inspect.ImageCounters{Total: 1, Local: 1},
		},
		Readiness:   inspect.Readiness{PreviewFidelity: inspect.PreviewFidelityExact},
		ArticleHTML: "<p>exact html</p>",
		ExactHTML:   true,
		Checks: []inspect.Check{{
			Level:        inspect.LevelWarn,
			Code:         "DUPLICATE_H1",
			Message:      "Body H1 matches the final article title",
			SuggestedFix: "remove the body H1 or change the final title",
		}},
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, want := range []string{
		"Publish Preview",
		"frontmatter.title",
		"cli.author",
		"frontmatter.summary",
		"warn · DUPLICATE_H1",
		"remove the body H1 or change the final title",
		"<p>exact html</p>",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("rendered page missing %q\n%s", want, html)
		}
	}
}

func TestRenderDegradedPageShowsFallbackMarkdownAndBanners(t *testing.T) {
	html, err := Render(PageData{
		Title:      "降级预览",
		SourceFile: "/tmp/article.md",
		Context: inspect.Context{
			Mode:           "ai",
			Theme:          "autumn-warm",
			FontSize:       "medium",
			BackgroundType: "none",
		},
		Metadata: inspect.MetadataState{
			Title:  inspect.MetadataField{Value: "降级预览", Source: "markdown.heading", Length: 4, Limit: 32},
			Author: inspect.MetadataField{Source: "fallback.empty", Limit: 16},
			Digest: inspect.MetadataField{Source: "fallback.empty", Limit: 128},
		},
		Readiness: inspect.Readiness{PreviewFidelity: inspect.PreviewFidelityDegraded},
		BodyMarkdown: strings.Join([]string{
			"# 降级预览",
			"",
			"正文段落",
		}, "\n"),
		ExactHTML:       false,
		RenderError:     "converter offline",
		DegradedMessage: "Preview degraded: exact HTML could not be rendered.",
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, want := range []string{
		"Preview degraded: exact HTML could not be rendered.",
		"Render error: converter offline",
		"Markdown Body",
		"# 降级预览",
		"正文段落",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("rendered page missing %q\n%s", want, html)
		}
	}
}

func TestRenderWithoutChecksShowsNoChecks(t *testing.T) {
	html, err := Render(PageData{
		Title:      "无检查",
		SourceFile: "/tmp/article.md",
		Metadata: inspect.MetadataState{
			Title:  inspect.MetadataField{Value: "无检查", Source: "frontmatter.title", Length: 3, Limit: 32},
			Author: inspect.MetadataField{Source: "fallback.empty", Limit: 16},
			Digest: inspect.MetadataField{Source: "fallback.empty", Limit: 128},
		},
		Readiness:   inspect.Readiness{PreviewFidelity: inspect.PreviewFidelityExact},
		ArticleHTML: "<p>ok</p>",
		ExactHTML:   true,
	})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(html, "No checks.") {
		t.Fatalf("rendered page missing no-checks state\n%s", html)
	}
}
