package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/image"
)

func TestResolveGenerateImagePromptWithPresetAndArticle(t *testing.T) {
	oldPreset, oldArticle := generateImageCmdPreset, generateImageCmdArticle
	oldTitle, oldSummary := generateImageCmdTitle, generateImageCmdSummary
	oldKeywords, oldStyle := generateImageCmdKeywords, generateImageCmdStyle
	oldAspect := generateImageCmdAspect
	t.Cleanup(func() {
		generateImageCmdPreset = oldPreset
		generateImageCmdArticle = oldArticle
		generateImageCmdTitle = oldTitle
		generateImageCmdSummary = oldSummary
		generateImageCmdKeywords = oldKeywords
		generateImageCmdStyle = oldStyle
		generateImageCmdAspect = oldAspect
	})

	article := strings.Join([]string{
		"---",
		"title: AI 时代的写作系统",
		"digest: 一篇关于写作工作流、提示词和信息组织的总结。",
		"---",
		"",
		"# 忽略这个标题",
		"",
		"正文第一段。",
	}, "\n")
	articlePath := filepath.Join(t.TempDir(), "article.md")
	if err := os.WriteFile(articlePath, []byte(article), 0600); err != nil {
		t.Fatalf("write article: %v", err)
	}

	generateImageCmdPreset = "cover-hero"
	generateImageCmdArticle = articlePath

	prompt, err := resolveGenerateImagePrompt(generateImageInput{
		Preset:  generateImageCmdPreset,
		Article: generateImageCmdArticle,
	})
	if err != nil {
		t.Fatalf("resolveGenerateImagePrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "AI 时代的写作系统") {
		t.Fatalf("prompt missing title: %q", prompt)
	}
	if !strings.Contains(prompt, "一篇关于写作工作流") {
		t.Fatalf("prompt missing summary: %q", prompt)
	}
	if !strings.Contains(prompt, "16:9") {
		t.Fatalf("prompt missing default aspect ratio: %q", prompt)
	}
}

func TestRunGenerateImageUsesPresetPrompt(t *testing.T) {
	oldCfg := cfg
	oldPreset, oldArticle := generateImageCmdPreset, generateImageCmdArticle
	oldTitle, oldSummary := generateImageCmdTitle, generateImageCmdSummary
	oldKeywords, oldStyle := generateImageCmdKeywords, generateImageCmdStyle
	oldAspect, oldSize, oldModel := generateImageCmdAspect, generateImageCmdSize, generateImageCmdModel
	oldNewImageProcessor, oldNewImageProcessorWithConfig := newImageProcessor, newImageProcessorWithConfig
	t.Cleanup(func() {
		cfg = oldCfg
		generateImageCmdPreset = oldPreset
		generateImageCmdArticle = oldArticle
		generateImageCmdTitle = oldTitle
		generateImageCmdSummary = oldSummary
		generateImageCmdKeywords = oldKeywords
		generateImageCmdStyle = oldStyle
		generateImageCmdAspect = oldAspect
		generateImageCmdSize = oldSize
		generateImageCmdModel = oldModel
		newImageProcessor = oldNewImageProcessor
		newImageProcessorWithConfig = oldNewImageProcessorWithConfig
	})

	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
	}
	generateImageCmdPreset = "infographic-comparison"
	generateImageCmdTitle = "提示词系统设计"
	generateImageCmdSummary = "比较不同图片提示词组织方式的优缺点"
	generateImageCmdStyle = "technical schematic"

	expectedPrompt, err := resolveGenerateImagePrompt(generateImageInput{
		Preset:  generateImageCmdPreset,
		Title:   generateImageCmdTitle,
		Summary: generateImageCmdSummary,
		Style:   generateImageCmdStyle,
	})
	if err != nil {
		t.Fatalf("resolveGenerateImagePrompt() error = %v", err)
	}

	processor := &fakeImageProcessor{
		generateResults: map[string]*image.GenerateAndUploadResult{
			expectedPrompt: {
				Prompt:      expectedPrompt,
				OriginalURL: "https://provider.example/image.png",
				MediaID:     "media-123",
				WechatURL:   "https://wechat.local/media-123",
			},
		},
	}
	newImageProcessor = func() imageProcessor { return processor }
	newImageProcessorWithConfig = func(runtimeCfg *config.Config) imageProcessor { return processor }

	stdout := captureStdout(t, func() {
		if err := runGenerateImage(nil); err != nil {
			t.Fatalf("runGenerateImage() error = %v", err)
		}
	})

	if len(processor.generateCalls) != 1 {
		t.Fatalf("generateCalls = %#v", processor.generateCalls)
	}
	if processor.generateCalls[0] != expectedPrompt {
		t.Fatalf("generate prompt = %q, want %q", processor.generateCalls[0], expectedPrompt)
	}

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true {
		t.Fatalf("unexpected response: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["media_id"] != "media-123" {
		t.Fatalf("unexpected response data: %#v", data)
	}
}

func TestRunGenerateImageUsesModelOverride(t *testing.T) {
	oldCfg := cfg
	oldNewImageProcessor := newImageProcessor
	oldNewImageProcessorWithConfig := newImageProcessorWithConfig
	t.Cleanup(func() {
		cfg = oldCfg
		newImageProcessor = oldNewImageProcessor
		newImageProcessorWithConfig = oldNewImageProcessorWithConfig
	})

	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
		ImageModel:   "default-model",
	}

	processor := &fakeImageProcessor{
		generateResults: map[string]*image.GenerateAndUploadResult{
			"test prompt": {
				Prompt:      "test prompt",
				OriginalURL: "https://provider.example/image.png",
				MediaID:     "media-override",
				WechatURL:   "https://wechat.local/media-override",
			},
		},
	}

	newImageProcessor = func() imageProcessor {
		t.Fatal("newImageProcessor should not be used when --model is set")
		return nil
	}

	newImageProcessorWithConfig = func(runtimeCfg *config.Config) imageProcessor {
		if runtimeCfg == cfg {
			t.Fatal("expected model override to use a config copy")
		}
		if runtimeCfg.ImageModel != "override-model" {
			t.Fatalf("ImageModel = %q, want override-model", runtimeCfg.ImageModel)
		}
		if cfg.ImageModel != "default-model" {
			t.Fatalf("original cfg.ImageModel mutated = %q", cfg.ImageModel)
		}
		return processor
	}

	stdout := captureStdout(t, func() {
		if err := runGenerateImageWithInput(generateImageInput{
			RawPrompt: "test prompt",
			Model:     "override-model",
		}); err != nil {
			t.Fatalf("runGenerateImageWithInput() error = %v", err)
		}
	})

	if len(processor.generateCalls) != 1 || processor.generateCalls[0] != "test prompt" {
		t.Fatalf("generateCalls = %#v", processor.generateCalls)
	}

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestResolveGenerateImagePromptRejectsMixedRawPromptAndPreset(t *testing.T) {
	oldPreset := generateImageCmdPreset
	t.Cleanup(func() { generateImageCmdPreset = oldPreset })

	generateImageCmdPreset = "cover-default"
	if _, err := resolveGenerateImagePrompt(generateImageInput{
		RawPrompt: "raw prompt",
		Preset:    generateImageCmdPreset,
	}); err == nil {
		t.Fatal("expected error for mixed raw prompt and preset")
	}
}

func TestRunGeneratePresetImageRejectsWrongArchetype(t *testing.T) {
	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
	}

	err := runGeneratePresetImage("cover", "cover-default", generateImageInput{
		Preset:  "infographic-default",
		Title:   "标题",
		Summary: "摘要",
	})
	if err == nil || !strings.Contains(err.Error(), "expected cover") {
		t.Fatalf("runGeneratePresetImage() error = %v", err)
	}
}

func TestRunGeneratePresetImageAllowsCompatibleCoverUseCase(t *testing.T) {
	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
	}

	_, err := resolveGenerateImagePrompt(generateImageInput{
		Preset:            "infographic-victorian-engraving-banner",
		Title:             "标题",
		Summary:           "摘要",
		RequiredArchetype: "cover",
	})
	if err != nil {
		t.Fatalf("resolveGenerateImagePrompt() error = %v", err)
	}
}

func TestResolveGenerateImagePromptUsesSpecDefaultAspectRatio(t *testing.T) {
	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
	}

	prompt, err := resolveGenerateImagePrompt(generateImageInput{
		Preset:  "infographic-victorian-engraving-banner",
		Title:   "标题",
		Summary: "摘要",
	})
	if err != nil {
		t.Fatalf("resolveGenerateImagePrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "21:9") {
		t.Fatalf("expected 21:9 default aspect ratio in prompt: %q", prompt)
	}
}

func TestRunGenerateCoverUsesDefaultPreset(t *testing.T) {
	oldCfg := cfg
	oldNewImageProcessor := newImageProcessor
	t.Cleanup(func() {
		cfg = oldCfg
		newImageProcessor = oldNewImageProcessor
	})

	cfg = &config.Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
		ImageAPIKey:  "image-key",
	}

	expectedPrompt, err := resolveGenerateImagePrompt(generateImageInput{
		Preset:            "cover-default",
		Title:             "标题",
		Summary:           "摘要",
		RequiredArchetype: "cover",
	})
	if err != nil {
		t.Fatalf("resolveGenerateImagePrompt() error = %v", err)
	}

	processor := &fakeImageProcessor{
		generateResults: map[string]*image.GenerateAndUploadResult{
			expectedPrompt: {
				Prompt:      expectedPrompt,
				OriginalURL: "https://provider.example/cover.png",
				MediaID:     "cover-1",
				WechatURL:   "https://wechat.local/cover-1",
			},
		},
	}
	newImageProcessor = func() imageProcessor { return processor }

	stdout := captureStdout(t, func() {
		if err := runGeneratePresetImage("cover", "cover-default", generateImageInput{
			Title:   "标题",
			Summary: "摘要",
		}); err != nil {
			t.Fatalf("runGeneratePresetImage() error = %v", err)
		}
	})

	if len(processor.generateCalls) != 1 || processor.generateCalls[0] != expectedPrompt {
		t.Fatalf("generateCalls = %#v", processor.generateCalls)
	}

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true {
		t.Fatalf("unexpected response: %#v", response)
	}
}
