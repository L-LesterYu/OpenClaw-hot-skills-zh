package preview

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/inspect"
)

type PageData struct {
	Title           string
	SourceFile      string
	Context         inspect.Context
	Metadata        inspect.MetadataState
	Structure       inspect.Structure
	Checks          []inspect.Check
	Readiness       inspect.Readiness
	ArticleHTML     template.HTML
	BodyMarkdown    string
	ExactHTML       bool
	RenderError     string
	DegradedMessage string
}

func Render(data PageData) (string, error) {
	tmpl, err := template.New("preview").Funcs(template.FuncMap{
		"trim": trim,
	}).Parse(pageTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f5f4ef;
      --panel: #fffdf7;
      --ink: #1f1f1f;
      --muted: #5f5a52;
      --border: #ddd6c7;
      --warn: #8c5c00;
      --error: #9b1c1c;
      --info: #1f5d8a;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Georgia, "Noto Serif SC", serif;
      color: var(--ink);
      background: linear-gradient(180deg, #efede5 0%, #f9f7f1 100%);
    }
    .shell {
      display: grid;
      grid-template-columns: minmax(280px, 360px) minmax(0, 1fr);
      gap: 24px;
      padding: 24px;
      min-height: 100vh;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 18px;
      padding: 18px;
      box-shadow: 0 12px 32px rgba(59, 43, 17, 0.08);
    }
    h1, h2, h3 { margin: 0 0 12px; }
    h1 { font-size: 24px; }
    h2 { font-size: 15px; text-transform: uppercase; letter-spacing: 0.08em; color: var(--muted); margin-top: 24px; }
    .meta-row { margin-bottom: 14px; }
    .meta-label { font-size: 12px; text-transform: uppercase; color: var(--muted); letter-spacing: 0.08em; }
    .meta-value { font-size: 16px; margin-top: 4px; word-break: break-word; }
    .meta-source { font-size: 12px; color: var(--muted); margin-top: 4px; }
    .pill {
      display: inline-block;
      border-radius: 999px;
      padding: 4px 10px;
      font-size: 12px;
      margin-right: 8px;
      margin-bottom: 8px;
      border: 1px solid var(--border);
      background: #fbf8f0;
    }
    .check {
      border-left: 4px solid var(--border);
      padding: 10px 12px;
      background: #fcfaf4;
      margin-bottom: 10px;
      border-radius: 8px;
    }
    .check.warn { border-color: var(--warn); }
    .check.error { border-color: var(--error); }
    .check.info { border-color: var(--info); }
    .check-code { font-size: 12px; color: var(--muted); margin-bottom: 4px; }
    .article-shell {
      max-width: 760px;
      margin: 0 auto;
      background: #fff;
      border: 1px solid var(--border);
      border-radius: 20px;
      padding: 28px;
      box-shadow: 0 16px 40px rgba(40, 33, 18, 0.10);
    }
    .banner {
      border: 1px dashed var(--border);
      border-radius: 14px;
      padding: 14px;
      margin-bottom: 18px;
      background: #faf6ec;
      color: var(--muted);
    }
    pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #f6f3eb;
      border-radius: 12px;
      padding: 16px;
      border: 1px solid var(--border);
      overflow: auto;
    }
    @media (max-width: 960px) {
      .shell { grid-template-columns: 1fr; }
      .article-shell { max-width: 100%; }
    }
  </style>
</head>
<body>
  <div class="shell">
    <aside class="panel">
      <h1>Publish Preview</h1>
      <div class="meta-source">{{.SourceFile}}</div>

      <h2>Resolved Metadata</h2>
      <div class="meta-row">
        <div class="meta-label">Title</div>
        <div class="meta-value">{{.Metadata.Title.Value}}</div>
        <div class="meta-source">{{.Metadata.Title.Source}} · {{.Metadata.Title.Length}} / {{.Metadata.Title.Limit}}</div>
      </div>
      <div class="meta-row">
        <div class="meta-label">Author</div>
        <div class="meta-value">{{if .Metadata.Author.Value}}{{.Metadata.Author.Value}}{{else}}—{{end}}</div>
        <div class="meta-source">{{.Metadata.Author.Source}} · {{.Metadata.Author.Length}} / {{.Metadata.Author.Limit}}</div>
      </div>
      <div class="meta-row">
        <div class="meta-label">Digest</div>
        <div class="meta-value">{{if .Metadata.Digest.Value}}{{.Metadata.Digest.Value}}{{else}}—{{end}}</div>
        <div class="meta-source">{{.Metadata.Digest.Source}} · {{.Metadata.Digest.Length}} / {{.Metadata.Digest.Limit}}</div>
      </div>

      <h2>Context</h2>
      <span class="pill">mode: {{.Context.Mode}}</span>
      <span class="pill">theme: {{.Context.Theme}}</span>
      <span class="pill">font: {{.Context.FontSize}}</span>
      <span class="pill">background: {{.Context.BackgroundType}}</span>
      <span class="pill">preview: {{.Readiness.PreviewFidelity}}</span>

      <h2>Structure</h2>
      <div class="meta-row">
        <div class="meta-label">Body H1</div>
        <div class="meta-value">{{if .Structure.BodyH1.Present}}{{.Structure.BodyH1.Text}}{{else}}—{{end}}</div>
      </div>
      <div class="meta-row">
        <div class="meta-label">Images</div>
        <div class="meta-value">{{.Structure.Images.Total}} total · {{.Structure.Images.Local}} local · {{.Structure.Images.Remote}} remote · {{.Structure.Images.AI}} ai</div>
      </div>

      <h2>Checks</h2>
      {{if .Checks}}
        {{range .Checks}}
          <div class="check {{.Level}}">
            <div class="check-code">{{.Level}} · {{.Code}}</div>
            <div>{{.Message}}</div>
            {{if .SuggestedFix}}<div class="meta-source">{{.SuggestedFix}}</div>{{end}}
          </div>
        {{end}}
      {{else}}
        <div class="meta-source">No checks.</div>
      {{end}}
    </aside>

    <main class="panel">
      <div class="article-shell">
        {{if .DegradedMessage}}
          <div class="banner">{{.DegradedMessage}}</div>
        {{end}}
        {{if .RenderError}}
          <div class="banner">Render error: {{.RenderError}}</div>
        {{end}}
        {{if .ExactHTML}}
          {{.ArticleHTML}}
        {{else}}
          <h1 style="margin-bottom:12px;">{{.Metadata.Title.Value}}</h1>
          {{if .Metadata.Author.Value}}<p style="color:#666; margin-top:0;">{{.Metadata.Author.Value}}</p>{{end}}
          {{if .Metadata.Digest.Value}}<p style="color:#666;">{{.Metadata.Digest.Value}}</p>{{end}}
          <h2 style="margin-top:24px;">Markdown Body</h2>
          <pre>{{trim .BodyMarkdown}}</pre>
        {{end}}
      </div>
    </main>
  </div>
</body>
</html>`

func trim(value string) string {
	return strings.TrimSpace(value)
}
