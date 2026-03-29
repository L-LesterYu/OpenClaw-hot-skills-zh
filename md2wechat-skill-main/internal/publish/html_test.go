package publish

import (
	"strings"
	"testing"
)

func TestInsertAndReplaceAssetPlaceholders(t *testing.T) {
	html := `<p>a</p><img src="./a.png"><p>b</p><img src="https://example.com/b.png">`
	assets := []AssetRef{
		{Index: 0, Source: "./a.png", Placeholder: "<!-- IMG:0 -->", PublicURL: "https://wechat.local/a"},
		{Index: 1, Source: "https://example.com/b.png", Placeholder: "<!-- IMG:1 -->", PublicURL: "https://wechat.local/b"},
	}

	withPlaceholders := InsertAssetPlaceholders(html, assets)
	if strings.Count(withPlaceholders, "<!-- IMG:") != 2 {
		t.Fatalf("placeholder HTML = %s", withPlaceholders)
	}

	replaced := ReplaceAssetPlaceholders(withPlaceholders, assets)
	if !strings.Contains(replaced, "https://wechat.local/a") || !strings.Contains(replaced, "https://wechat.local/b") {
		t.Fatalf("replaced HTML = %s", replaced)
	}
}

func TestInsertAssetPlaceholdersFallsBackToDocumentOrder(t *testing.T) {
	html := `<div><img src="https://cdn.example.com/1"><img src="https://cdn.example.com/2"></div>`
	assets := []AssetRef{
		{Index: 0, Source: "./a.png", Placeholder: "<!-- IMG:0 -->", PublicURL: "https://wechat.local/a"},
		{Index: 1, Source: "./b.png", Placeholder: "<!-- IMG:1 -->", PublicURL: "https://wechat.local/b"},
	}

	withPlaceholders := InsertAssetPlaceholders(html, assets)
	if !strings.Contains(withPlaceholders, "<!-- IMG:0 -->") || !strings.Contains(withPlaceholders, "<!-- IMG:1 -->") {
		t.Fatalf("fallback placeholder HTML = %s", withPlaceholders)
	}

	replaced := ReplaceAssetPlaceholders(withPlaceholders, assets)
	if !strings.Contains(replaced, "https://wechat.local/a") || !strings.Contains(replaced, "https://wechat.local/b") {
		t.Fatalf("fallback replaced HTML = %s", replaced)
	}
}
