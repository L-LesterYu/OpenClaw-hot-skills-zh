package image

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"go.uber.org/zap"
)

type fakeProvider struct {
	result *GenerateResult
	err    error
}

func (f *fakeProvider) Name() string { return "fake" }

func (f *fakeProvider) Generate(ctx context.Context, prompt string) (*GenerateResult, error) {
	return f.result, f.err
}

func TestNewProcessorInitializesDownloadHelper(t *testing.T) {
	downloadCalled := false
	p := NewProcessor(
		&config.Config{},
		zap.NewNop(),
		WithDownloadFunc(func(url string) (string, error) {
			downloadCalled = true
			return "", nil
		}),
	)
	if p.downloadFile == nil {
		t.Fatal("download helper should be initialized")
	}
	if _, err := p.download("https://example.com/a.png"); err != nil {
		t.Fatalf("download() error = %v", err)
	}
	if !downloadCalled {
		t.Fatal("expected injected download helper to be used")
	}
}

func TestDownloadAndUploadUsesInjectedHelpers(t *testing.T) {
	dir := t.TempDir()
	downloadedPath := filepath.Join(dir, "downloaded.png")
	if err := os.WriteFile(downloadedPath, []byte("png"), 0644); err != nil {
		t.Fatalf("write temp image: %v", err)
	}

	var gotDownloadURL string
	var gotUploadPath string
	p := &Processor{
		cfg: &config.Config{
			CompressImages: false,
		},
		log:        zap.NewNop(),
		compressor: NewCompressor(zap.NewNop(), 0, 0),
		downloadFile: func(url string) (string, error) {
			gotDownloadURL = url
			return downloadedPath, nil
		},
		uploadMaterial: func(filePath string) (*UploadResult, error) {
			gotUploadPath = filePath
			return &UploadResult{
				MediaID:   "media-1",
				WechatURL: "https://wechat.local/1",
			}, nil
		},
	}

	result, err := p.DownloadAndUpload("https://example.com/image.png")
	if err != nil {
		t.Fatalf("DownloadAndUpload() error = %v", err)
	}
	if gotDownloadURL != "https://example.com/image.png" {
		t.Fatalf("download url = %q", gotDownloadURL)
	}
	if gotUploadPath != downloadedPath {
		t.Fatalf("upload path = %q, want %q", gotUploadPath, downloadedPath)
	}
	if result.MediaID != "media-1" || result.WechatURL != "https://wechat.local/1" {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(downloadedPath); !os.IsNotExist(err) {
		t.Fatalf("downloaded file should be cleaned up, stat err = %v", err)
	}
}

func TestGenerateAndUploadUsesInjectedHelpers(t *testing.T) {
	dir := t.TempDir()
	downloadedPath := filepath.Join(dir, "generated.png")
	if err := os.WriteFile(downloadedPath, []byte("png"), 0644); err != nil {
		t.Fatalf("write temp image: %v", err)
	}

	var gotDownloadURL string
	var gotUploadPath string
	p := &Processor{
		cfg: &config.Config{
			WechatAppID:    "appid",
			WechatSecret:   "secret",
			ImageAPIKey:    "image-key",
			CompressImages: false,
		},
		log:        zap.NewNop(),
		compressor: NewCompressor(zap.NewNop(), 0, 0),
		downloadFile: func(url string) (string, error) {
			gotDownloadURL = url
			return downloadedPath, nil
		},
		uploadMaterial: func(filePath string) (*UploadResult, error) {
			gotUploadPath = filePath
			return &UploadResult{
				MediaID:   "media-2",
				WechatURL: "https://wechat.local/2",
			}, nil
		},
		provider: &fakeProvider{
			result: &GenerateResult{
				URL:   "https://provider.example/generated.png",
				Model: "fake-model",
				Size:  "2K",
			},
		},
	}

	result, err := p.GenerateAndUpload("draw a fox")
	if err != nil {
		t.Fatalf("GenerateAndUpload() error = %v", err)
	}
	if gotDownloadURL != "https://provider.example/generated.png" {
		t.Fatalf("download url = %q", gotDownloadURL)
	}
	if gotUploadPath != downloadedPath {
		t.Fatalf("upload path = %q, want %q", gotUploadPath, downloadedPath)
	}
	if result.OriginalURL != "https://provider.example/generated.png" {
		t.Fatalf("original url = %q", result.OriginalURL)
	}
	if result.MediaID != "media-2" || result.WechatURL != "https://wechat.local/2" {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(downloadedPath); !os.IsNotExist(err) {
		t.Fatalf("downloaded file should be cleaned up, stat err = %v", err)
	}
}
