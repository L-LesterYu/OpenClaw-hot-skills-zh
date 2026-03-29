package humanizer

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
	"github.com/geekjourneyx/md2wechat-skill/internal/promptcatalog"
)

func TestHumanizeMethodsShareAIRequestContract(t *testing.T) {
	h := NewHumanizer()
	req := &HumanizeRequest{Content: "这是一段需要去痕处理的测试文本。"}

	results := []*HumanizeResult{
		h.Humanize(req),
		h.HumanizeWithResult(&HumanizeRequest{Content: "这是一段需要去痕处理的测试文本。"}),
	}

	for i, result := range results {
		if !result.Success {
			t.Fatalf("result %d Success = false: %+v", i, result)
		}
		if result.Status != action.StatusActionRequired {
			t.Fatalf("result %d Status = %q", i, result.Status)
		}
		if result.Action != HumanizeActionAIRequest {
			t.Fatalf("result %d Action = %q", i, result.Action)
		}
		if !result.RequiresAI() {
			t.Fatalf("result %d expected RequiresAI()", i)
		}
		if result.Prompt == "" {
			t.Fatalf("result %d missing prompt", i)
		}
		if result.Content != "" {
			t.Fatalf("result %d content = %q, want empty", i, result.Content)
		}
	}
}

func TestParseAIResponseMarksCompletedState(t *testing.T) {
	h := NewHumanizer()
	req := &HumanizeRequest{Content: "原文"}
	response := strings.Join([]string{
		"# 人性化后的文本",
		"",
		"改写后的正文。",
		"",
		"# 修改说明",
		"",
		"删掉了明显的 AI 腔。",
	}, "\n")

	result := h.ParseAIResponse(response, req)
	if !result.Success {
		t.Fatalf("ParseAIResponse() failed: %+v", result)
	}
	if result.Status != action.StatusCompleted {
		t.Fatalf("Status = %q", result.Status)
	}
	if result.Action != HumanizeActionCompleted {
		t.Fatalf("Action = %q", result.Action)
	}
	if result.RequiresAI() {
		t.Fatalf("completed result should not require AI: %+v", result)
	}
	if result.Content != "改写后的正文。" {
		t.Fatalf("Content = %q", result.Content)
	}
	if result.Report != "删掉了明显的 AI 腔。" {
		t.Fatalf("Report = %q", result.Report)
	}
}

func TestParseAIResponseFallsBackToOriginalContentOnParseFailure(t *testing.T) {
	h := NewHumanizer()
	req := &HumanizeRequest{Content: "原始文本"}

	result := h.ParseAIResponse("", req)
	if result.Success {
		t.Fatalf("expected parse failure result: %+v", result)
	}
	if result.Status != action.StatusFailed {
		t.Fatalf("Status = %q", result.Status)
	}
	if result.Content != "原始文本" {
		t.Fatalf("Content = %q", result.Content)
	}
	if result.Error == "" {
		t.Fatalf("expected error message")
	}
}

func TestBuildPromptUsesBundledPromptAssets(t *testing.T) {
	promptcatalog.ResetDefaultCatalogForTests()

	prompt := BuildPrompt(&HumanizeRequest{
		Content:      "需要处理的文本",
		Intensity:    IntensityAggressive,
		ShowChanges:  true,
		IncludeScore: true,
	})

	if !strings.Contains(prompt, "激进模式") {
		t.Fatalf("prompt missing aggressive section: %q", prompt)
	}
	if !strings.Contains(prompt, "# Humanizer-zh: 去除 AI 写作痕迹") {
		t.Fatalf("prompt missing bundled base template: %q", prompt)
	}
	if !strings.Contains(prompt, "需要处理的文本") {
		t.Fatalf("prompt missing content: %q", prompt)
	}
}
