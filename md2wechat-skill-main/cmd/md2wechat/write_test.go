package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func chdirToRepoRoot(t *testing.T) {
	t.Helper()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join(oldWD, "..", ".."))
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir(%q) error = %v", repoRoot, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	})
}

func captureStdout(t *testing.T, fn func()) []byte {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}

	return buf.Bytes()
}

func TestExecuteWriteOutputsAIRequestWithHumanizerAndCover(t *testing.T) {
	chdirToRepoRoot(t)

	oldStyle, oldInputType := writeStyle, writeInputType
	oldLength, oldTitle := writeLength, writeTitle
	oldOutput, oldCover := writeOutput, writeCover
	oldCoverOnly, oldHumanize := writeCoverOnly, writeHumanize
	oldIntensity := writeHumanizeIntensity
	t.Cleanup(func() {
		writeStyle, writeInputType = oldStyle, oldInputType
		writeLength, writeTitle = oldLength, oldTitle
		writeOutput, writeCover = oldOutput, oldCover
		writeCoverOnly, writeHumanize = oldCoverOnly, oldHumanize
		writeHumanizeIntensity = oldIntensity
	})

	writeStyle = "dan-koe"
	writeInputType = "idea"
	writeLength = "long"
	writeTitle = "测试标题"
	writeOutput = ""
	writeCover = true
	writeCoverOnly = false
	writeHumanize = true
	writeHumanizeIntensity = "aggressive"

	stdout := captureStdout(t, func() {
		if err := executeWrite("这是一段足够长的测试输入，用来验证 write 命令的 AI 请求输出、humanizer 模板和封面提示词是否都被正确组装。"); err != nil {
			t.Fatalf("executeWrite() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}

	if response["success"] != true || response["code"] != "WRITE_AI_REQUEST_READY" {
		t.Fatalf("unexpected response: %#v", response)
	}
	if response["schema_version"] != "v1" || response["status"] != "action_required" || response["retryable"] != false {
		t.Fatalf("unexpected envelope: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["mode"] != "ai" || data["action"] != "ai_write_request" {
		t.Fatalf("unexpected data payload: %#v", data)
	}
	if _, ok := data["cover_prompt"].(string); !ok {
		t.Fatalf("expected cover_prompt in response: %#v", data)
	}

	humanizerBlock, ok := data["humanizer"].(map[string]any)
	if !ok {
		t.Fatalf("expected humanizer block: %#v", data)
	}
	if humanizerBlock["enabled"] != true || humanizerBlock["intensity"] != "aggressive" {
		t.Fatalf("unexpected humanizer block: %#v", humanizerBlock)
	}
	if promptTemplate, _ := humanizerBlock["prompt_template"].(string); !strings.Contains(promptTemplate, "Humanizer-zh") {
		t.Fatalf("unexpected humanizer prompt: %#v", humanizerBlock)
	}
}

func TestExecuteWriteJSONFlagWrapsAIResponseInEnvelope(t *testing.T) {
	chdirToRepoRoot(t)

	oldStyle, oldInputType := writeStyle, writeInputType
	oldLength, oldOutput := writeLength, writeOutput
	oldCover, oldHumanize := writeCover, writeHumanize
	oldJSON := jsonOutput
	t.Cleanup(func() {
		writeStyle, writeInputType = oldStyle, oldInputType
		writeLength, writeOutput = oldLength, oldOutput
		writeCover, writeHumanize = oldCover, oldHumanize
		jsonOutput = oldJSON
	})

	writeStyle = "dan-koe"
	writeInputType = "idea"
	writeLength = "medium"
	writeOutput = ""
	writeCover = false
	writeHumanize = false
	jsonOutput = true

	stdout := captureStdout(t, func() {
		if err := executeWrite("这是一段足够长的测试输入，用来验证 --json 模式下 write 命令会返回稳定 envelope。"); err != nil {
			t.Fatalf("executeWrite() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true || response["code"] != "WRITE_AI_REQUEST_READY" {
		t.Fatalf("unexpected response: %#v", response)
	}
	if response["schema_version"] != "v1" || response["status"] != "action_required" || response["retryable"] != false {
		t.Fatalf("unexpected envelope: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["action"] != "ai_write_request" {
		t.Fatalf("unexpected data payload: %#v", data)
	}
}

func TestRunWriteReadsFileAndSwitchesInputTypeToFragment(t *testing.T) {
	chdirToRepoRoot(t)

	oldStyle, oldInputType := writeStyle, writeInputType
	oldLength, oldTitle := writeLength, writeTitle
	oldOutput, oldCover := writeOutput, writeCover
	oldCoverOnly, oldHumanize := writeCoverOnly, writeHumanize
	oldList, oldDetail := writeListStyles, writeStyleDetail
	oldStdin := os.Stdin
	t.Cleanup(func() {
		writeStyle, writeInputType = oldStyle, oldInputType
		writeLength, writeTitle = oldLength, oldTitle
		writeOutput, writeCover = oldOutput, oldCover
		writeCoverOnly, writeHumanize = oldCoverOnly, oldHumanize
		writeListStyles, writeStyleDetail = oldList, oldDetail
		os.Stdin = oldStdin
	})

	tmpDir := t.TempDir()
	markdownPath := filepath.Join(tmpDir, "article.md")
	if err := os.WriteFile(markdownPath, []byte("这是一段来自文件的测试内容，用来验证 runWrite 会自动把输入类型从 idea 切到 fragment。"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	devNull, err := os.Open("/dev/null")
	if err != nil {
		t.Fatalf("open /dev/null: %v", err)
	}
	defer func() {
		_ = devNull.Close()
	}()
	os.Stdin = devNull

	writeStyle = "dan-koe"
	writeInputType = "idea"
	writeLength = "medium"
	writeTitle = ""
	writeOutput = ""
	writeCover = false
	writeCoverOnly = false
	writeHumanize = false
	writeListStyles = false
	writeStyleDetail = false

	stdout := captureStdout(t, func() {
		if err := runWrite(nil, []string{markdownPath}); err != nil {
			t.Fatalf("runWrite() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	data, _ := response["data"].(map[string]any)
	prompt, _ := data["prompt"].(string)
	if !strings.Contains(prompt, "输入类型: fragment") {
		t.Fatalf("prompt did not use fragment input type: %q", prompt)
	}
}

func TestRunListStylesOutputsAvailableStyles(t *testing.T) {
	chdirToRepoRoot(t)

	oldList, oldDetail := writeListStyles, writeStyleDetail
	t.Cleanup(func() {
		writeListStyles, writeStyleDetail = oldList, oldDetail
	})

	writeListStyles = true
	writeStyleDetail = false

	stdout := captureStdout(t, func() {
		if err := runListStyles(); err != nil {
			t.Fatalf("runListStyles() error = %v", err)
		}
	})

	if !strings.Contains(string(stdout), "dan-koe") {
		t.Fatalf("expected dan-koe in styles output: %s", stdout)
	}
}
