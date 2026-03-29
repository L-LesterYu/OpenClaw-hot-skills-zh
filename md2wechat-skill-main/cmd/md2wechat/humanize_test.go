package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHumanizeOutputsRequestAndWritesPromptFile(t *testing.T) {
	oldIntensity, oldShowChanges := intensityFlag, showChangesFlag
	oldOutput := outputFlag
	oldJSON := jsonOutput
	t.Cleanup(func() {
		intensityFlag, showChangesFlag = oldIntensity, oldShowChanges
		outputFlag = oldOutput
		jsonOutput = oldJSON
	})

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "article.md")
	outputPath := filepath.Join(tmpDir, "prompt.txt")
	content := "这是一段需要去痕处理的测试文本。"
	if err := os.WriteFile(inputPath, []byte(content), 0644); err != nil {
		t.Fatalf("write input file: %v", err)
	}

	intensityFlag = "gentle"
	showChangesFlag = true
	outputFlag = outputPath

	stdout := captureStdout(t, func() {
		if err := runHumanize(inputPath); err != nil {
			t.Fatalf("runHumanize() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true || response["code"] != "HUMANIZE_REQUEST_READY" {
		t.Fatalf("unexpected response: %#v", response)
	}
	if response["schema_version"] != "v1" || response["status"] != "action_required" || response["retryable"] != false {
		t.Fatalf("unexpected envelope: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["action"] != "humanize_request" {
		t.Fatalf("unexpected data payload: %#v", data)
	}
	if data["output_file"] != outputPath {
		t.Fatalf("output_file = %#v", data["output_file"])
	}

	request, ok := data["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected request block: %#v", data)
	}
	if request["content"] != content || request["intensity"] != "gentle" {
		t.Fatalf("unexpected request block: %#v", request)
	}

	promptData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read prompt file: %v", err)
	}
	if !strings.Contains(string(promptData), "Humanizer-zh") {
		t.Fatalf("unexpected prompt file: %s", promptData)
	}
}

func TestRunHumanizeJSONFlagWrapsResponseInEnvelope(t *testing.T) {
	oldIntensity, oldShowChanges := intensityFlag, showChangesFlag
	oldOutput, oldJSON := outputFlag, jsonOutput
	t.Cleanup(func() {
		intensityFlag, showChangesFlag = oldIntensity, oldShowChanges
		outputFlag, jsonOutput = oldOutput, oldJSON
	})

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "article.md")
	if err := os.WriteFile(inputPath, []byte("这是一段需要去痕处理的测试文本。"), 0644); err != nil {
		t.Fatalf("write input file: %v", err)
	}

	intensityFlag = "medium"
	showChangesFlag = false
	outputFlag = ""
	jsonOutput = true

	stdout := captureStdout(t, func() {
		if err := runHumanize(inputPath); err != nil {
			t.Fatalf("runHumanize() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true || response["code"] != "HUMANIZE_REQUEST_READY" {
		t.Fatalf("unexpected response: %#v", response)
	}
	if response["schema_version"] != "v1" || response["status"] != "action_required" || response["retryable"] != false {
		t.Fatalf("unexpected envelope: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["action"] != "humanize_request" {
		t.Fatalf("unexpected data payload: %#v", data)
	}
}

func TestRunHumanizeReturnsReadErrorForMissingFile(t *testing.T) {
	if err := runHumanize("/nonexistent/article.md"); err == nil || !strings.Contains(err.Error(), "读取文件失败") {
		t.Fatalf("runHumanize() error = %v", err)
	}
}

func TestParseHumanizeResponseReturnsStructuredFields(t *testing.T) {
	oldShowChanges := showChangesFlag
	t.Cleanup(func() {
		showChangesFlag = oldShowChanges
	})
	showChangesFlag = true

	aiResponse := strings.Join([]string{
		"# 人性化后的文本",
		"",
		"改写后的正文。",
		"",
		"# 修改说明",
		"",
		"整体删掉了明显的 AI 腔。",
		"",
		"# 质量评分",
		"",
		"| 维度 | 得分 | 说明 |",
		"|------|------|------|",
		"| 直接性 | 8/10 | 更直接 |",
		"| 节奏 | 7/10 | 更自然 |",
		"| 信任度 | 8/10 | 更可信 |",
		"| 真实性 | 9/10 | 更像人写的 |",
		"| 精炼度 | 8/10 | 更紧凑 |",
		"| 总分 | 40/50 | 良好 |",
	}, "\n")

	output := parseHumanizeResponse(aiResponse, "原文", "medium")
	if output["success"] != true || output["content"] != "改写后的正文。" {
		t.Fatalf("unexpected output: %#v", output)
	}
	if output["report"] != "整体删掉了明显的 AI 腔。" {
		t.Fatalf("report = %#v", output["report"])
	}
	score, ok := output["score"].(map[string]any)
	if !ok {
		t.Fatalf("expected score block: %#v", output)
	}
	if score["total"] != 40 {
		t.Fatalf("score = %#v", score)
	}
}
