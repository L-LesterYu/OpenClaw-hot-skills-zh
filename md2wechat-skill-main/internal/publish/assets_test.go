package publish

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/image"
)

func TestAssetPipelineProcessRewritesHTMLAndAssets(t *testing.T) {
	dir := t.TempDir()
	localPath := filepath.Join(dir, "images", "local.png")

	pipeline := NewAssetPipeline(&fakeAssetProcessor{
		localResults: map[string]*image.UploadResult{
			localPath: {MediaID: "m-local", WechatURL: "https://wechat.local/local"},
		},
		onlineResults: map[string]*image.UploadResult{
			"https://example.com/remote.png": {MediaID: "m-remote", WechatURL: "https://wechat.local/remote"},
		},
		generateResults: map[string]*image.GenerateAndUploadResult{
			"draw fox": {MediaID: "m-ai", WechatURL: "https://wechat.local/ai"},
		},
	})

	output, err := pipeline.Process(&ProcessInput{
		HTML: `<img src="images/local.png"><img src="https://example.com/remote.png"><img src="https://old.example/ai.png">`,
		Assets: []AssetRef{
			{Index: 0, Kind: AssetKindLocal, Source: filepath.Join("images", "local.png"), Placeholder: "<!-- IMG:0 -->"},
			{Index: 1, Kind: AssetKindRemote, Source: "https://example.com/remote.png", Placeholder: "<!-- IMG:1 -->"},
			{Index: 2, Kind: AssetKindAI, Source: "draw fox", Prompt: "draw fox", Placeholder: "<!-- IMG:2 -->"},
		},
		MarkdownDir: dir,
	})
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	for _, expected := range []string{
		"https://wechat.local/local",
		"https://wechat.local/remote",
		"https://wechat.local/ai",
	} {
		if !strings.Contains(output.HTML, expected) {
			t.Fatalf("output HTML missing %q: %s", expected, output.HTML)
		}
	}
	if output.Assets[0].ResolvedSource != localPath {
		t.Fatalf("resolved source = %q, want %q", output.Assets[0].ResolvedSource, localPath)
	}
	if output.Assets[2].MediaID != "m-ai" {
		t.Fatalf("AI asset media id = %q", output.Assets[2].MediaID)
	}
}

func TestAssetPipelineProcessReturnsTypedStageErrorInput(t *testing.T) {
	pipeline := NewAssetPipeline(&fakeAssetProcessor{
		onlineErrs: map[string]error{
			"https://example.com/remote.png": fmt.Errorf("download failed"),
		},
	})

	_, err := pipeline.Process(&ProcessInput{
		HTML: `<img src="https://example.com/remote.png">`,
		Assets: []AssetRef{
			{Index: 0, Kind: AssetKindRemote, Source: "https://example.com/remote.png", Placeholder: "<!-- IMG:0 -->"},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "download failed") {
		t.Fatalf("Process() error = %v", err)
	}
}
