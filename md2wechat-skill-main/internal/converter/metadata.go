package converter

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// ArticleMetadata 表示 Markdown 中可用于发布的元信息。
type ArticleMetadata struct {
	Title  string
	Author string
	Digest string
}

// ArticleDocument 表示解析后的文章元信息与正文。
type ArticleDocument struct {
	Metadata ArticleMetadata
	Body     string
}

type frontMatter struct {
	Title       string `yaml:"title"`
	Author      string `yaml:"author"`
	Digest      string `yaml:"digest"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
}

// ParseArticleMetadata 提取 frontmatter 和正文中的元信息。
func ParseArticleMetadata(markdown string) ArticleMetadata {
	return ParseArticleDocument(markdown).Metadata
}

// ParseArticleDocument 提取 frontmatter 元信息，并返回去除 frontmatter 后的正文。
func ParseArticleDocument(markdown string) ArticleDocument {
	meta := ArticleMetadata{}
	body := markdown

	if fm, parsedBody, ok := parseFrontMatter(markdown); ok {
		meta.Title = strings.TrimSpace(fm.Title)
		meta.Author = strings.TrimSpace(fm.Author)
		meta.Digest = firstNonEmpty(fm.Digest, fm.Summary, fm.Description)
		body = parsedBody
	}

	if meta.Title == "" {
		meta.Title = ParseMarkdownTitle(body)
	}

	return ArticleDocument{
		Metadata: meta,
		Body:     body,
	}
}

func parseFrontMatter(markdown string) (frontMatter, string, bool) {
	var fm frontMatter
	normalized := normalizeMarkdownNewlines(markdown)
	lines := strings.Split(normalized, "\n")

	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return fm, markdown, false
	}

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "---" {
			continue
		}

		frontMatterBody := strings.Join(lines[1:i], "\n")
		if err := yaml.Unmarshal([]byte(frontMatterBody), &fm); err != nil {
			return frontMatter{}, markdown, false
		}

		body := strings.Join(lines[i+1:], "\n")
		return fm, body, true
	}

	return fm, markdown, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeMarkdownNewlines(markdown string) string {
	markdown = strings.TrimPrefix(markdown, "\uFEFF")
	return strings.ReplaceAll(markdown, "\r\n", "\n")
}
