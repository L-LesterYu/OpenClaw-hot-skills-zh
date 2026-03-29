package image

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

func TestNewModelScopeProvider(t *testing.T) {
	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://api-inference.modelscope.cn/",
		ImageModel:   "Tongyi-MAI/Z-Image-Turbo",
	}

	p, err := NewModelScopeProvider(cfg)
	if err != nil {
		t.Fatalf("NewModelScopeProvider() error = %v", err)
	}

	if p.Name() != "ModelScope" {
		t.Errorf("Name() = %v, want ModelScope", p.Name())
	}

	if p.model != "Tongyi-MAI/Z-Image-Turbo" {
		t.Errorf("model = %v, want Tongyi-MAI/Z-Image-Turbo", p.model)
	}

	if p.apiKey != "test-key" {
		t.Errorf("apiKey = %v, want test-key", p.apiKey)
	}
}

func TestNewModelScopeProviderDefaults(t *testing.T) {
	cfg := &config.Config{
		ImageAPIKey: "test-key",
		// ImageAPIBase 和 ImageModel 为空，应使用默认值
	}

	p, err := NewModelScopeProvider(cfg)
	if err != nil {
		t.Fatalf("NewModelScopeProvider() error = %v", err)
	}

	if p.baseURL != "https://api-inference.modelscope.cn/" {
		t.Errorf("baseURL = %v, want https://api-inference.modelscope.cn/", p.baseURL)
	}

	if p.model != "Tongyi-MAI/Z-Image-Turbo" {
		t.Errorf("model = %v, want Tongyi-MAI/Z-Image-Turbo", p.model)
	}

	if p.pollInterval != 5*time.Second {
		t.Errorf("pollInterval = %v, want 5s", p.pollInterval)
	}

	if p.maxPollTime != 120*time.Second {
		t.Errorf("maxPollTime = %v, want 120s", p.maxPollTime)
	}
}

func TestParseModelScopeSize(t *testing.T) {
	tests := []struct {
		name       string
		size       string
		wantWidth  int
		wantHeight int
		wantErr    string
	}{
		{
			name:       "square",
			size:       "1024x1024",
			wantWidth:  1024,
			wantHeight: 1024,
		},
		{
			name:       "portrait",
			size:       "1536x2048",
			wantWidth:  1536,
			wantHeight: 2048,
		},
		{
			name:       "landscape",
			size:       "1920x1080",
			wantWidth:  1920,
			wantHeight: 1080,
		},
		{
			name:    "aspect ratio not supported",
			size:    "16:9",
			wantErr: "expected WIDTHxHEIGHT",
		},
		{
			name:    "ultrawide aspect ratio not supported",
			size:    "21:9",
			wantErr: "expected WIDTHxHEIGHT",
		},
		{
			name:    "invalid text",
			size:    "foo",
			wantErr: "expected WIDTHxHEIGHT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight, err := parseSize(tt.size)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("parseSize(%q) should return error", tt.size)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("parseSize(%q) error = %v, want substring %q", tt.size, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseSize(%q) error = %v", tt.size, err)
			}
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Fatalf("parseSize(%q) = %dx%d, want %dx%d", tt.size, gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestModelScopeProvider_CreateTask_InvalidSize(t *testing.T) {
	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
		ImageSize:    "16:9",
	}

	p, err := NewModelScopeProvider(cfg)
	if err != nil {
		t.Fatalf("NewModelScopeProvider() error = %v", err)
	}

	_, err = p.createTask(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("createTask() should return error for aspect ratio size")
	}

	genErr, ok := err.(*GenerateError)
	if !ok {
		t.Fatalf("Error type = %T, want *GenerateError", err)
	}

	if genErr.Code != "invalid_size" {
		t.Fatalf("Error code = %v, want invalid_size", genErr.Code)
	}

	if !strings.Contains(genErr.Message, "WIDTHxHEIGHT") {
		t.Fatalf("Error message = %q, want WIDTHxHEIGHT hint", genErr.Message)
	}
}

func TestModelScopeProvider_CreateTask(t *testing.T) {
	taskID := "test-task-123"
	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
		ImageModel:   "Tongyi-MAI/Z-Image-Turbo",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		// 验证请求
		if r.Method != "POST" {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		if r.URL.Path != "/v1/images/generations" {
			t.Errorf("Path = %v, want /v1/images/generations", r.URL.Path)
		}
		if r.Header.Get("X-ModelScope-Async-Mode") != "true" {
			t.Errorf("X-ModelScope-Async-Mode header missing")
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Authorization header incorrect")
		}

		return jsonResponse(http.StatusOK, map[string]string{
			"task_id": taskID,
		}), nil
	})
	gotTaskID, err := p.createTask(context.Background(), "a golden cat")
	if err != nil {
		t.Fatalf("createTask() error = %v", err)
	}

	if gotTaskID != taskID {
		t.Errorf("taskID = %v, want %v", gotTaskID, taskID)
	}
}

func TestModelScopeProvider_PollTaskStatus(t *testing.T) {
	imageURL := "https://example.com/generated-image.png"
	pollCount := 0
	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != "GET" {
			t.Errorf("Method = %v, want GET", r.Method)
		}
		if r.Header.Get("X-ModelScope-Task-Type") != "image_generation" {
			t.Errorf("X-ModelScope-Task-Type header missing")
		}

		pollCount++
		var response any
		if pollCount < 3 {
			// 前两次返回 PENDING/RUNNING
			status := "PENDING"
			if pollCount == 2 {
				status = "RUNNING"
			}
			response = map[string]any{
				"task_status": status,
			}
		} else {
			// 第三次返回成功
			response = map[string]any{
				"task_status":   "SUCCEED",
				"output_images": []string{imageURL},
			}
		}

		return jsonResponse(http.StatusOK, response), nil
	})
	p.pollInterval = 10 * time.Millisecond // 加快测试速度

	gotURL, err := p.pollTaskStatus(context.Background(), "test-task-id")
	if err != nil {
		t.Fatalf("pollTaskStatus() error = %v", err)
	}

	if gotURL != imageURL {
		t.Errorf("URL = %v, want %v", gotURL, imageURL)
	}

	if pollCount != 3 {
		t.Errorf("poll count = %v, want 3", pollCount)
	}
}

func TestModelScopeProvider_Generate(t *testing.T) {
	imageURL := "https://example.com/generated-image.png"
	taskID := "test-task-456"
	var createCalled bool
	var statusCalled int

	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method == "POST" && r.URL.Path == "/v1/images/generations" {
			createCalled = true
			return jsonResponse(http.StatusOK, map[string]string{
				"task_id": taskID,
			}), nil
		} else if r.Method == "GET" && r.URL.Path == "/v1/tasks/"+taskID {
			statusCalled++
			// 第一次直接返回成功，简化测试
			return jsonResponse(http.StatusOK, map[string]any{
				"task_status":   "SUCCEED",
				"output_images": []string{imageURL},
			}), nil
		}
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "not found"}), nil
	})
	p.pollInterval = 10 * time.Millisecond

	result, err := p.Generate(context.Background(), "a golden cat")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !createCalled {
		t.Error("createTask was not called")
	}

	if statusCalled == 0 {
		t.Error("getTaskStatus was not called")
	}

	if result.URL != imageURL {
		t.Errorf("URL = %v, want %v", result.URL, imageURL)
	}

	if result.Model != "Tongyi-MAI/Z-Image-Turbo" {
		t.Errorf("Model = %v, want Tongyi-MAI/Z-Image-Turbo", result.Model)
	}
}

func TestModelScopeProvider_Generate_TaskFailed(t *testing.T) {
	taskID := "test-task-failed"

	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method == "POST" && r.URL.Path == "/v1/images/generations" {
			return jsonResponse(http.StatusOK, map[string]string{
				"task_id": taskID,
			}), nil
		} else if r.Method == "GET" && r.URL.Path == "/v1/tasks/"+taskID {
			return jsonResponse(http.StatusOK, map[string]any{
				"task_status": "FAILED",
			}), nil
		}
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "not found"}), nil
	})
	p.pollInterval = 10 * time.Millisecond

	_, err := p.Generate(context.Background(), "test prompt")
	if err == nil {
		t.Fatal("Generate() should return error for failed task")
	}

	genErr, ok := err.(*GenerateError)
	if !ok {
		t.Fatalf("Error type = %T, want *GenerateError", err)
	}

	if genErr.Code != "task_failed" {
		t.Errorf("Error code = %v, want task_failed", genErr.Code)
	}
}

func TestModelScopeProvider_HandleErrorResponse_Unauthorized(t *testing.T) {
	cfg := &config.Config{
		ImageAPIKey:  "invalid-key",
		ImageAPIBase: "https://mock.local",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"message": "invalid token"}`), nil
	})
	_, err := p.Generate(context.Background(), "test")

	if err == nil {
		t.Fatal("Expected error for unauthorized request")
	}

	genErr, ok := err.(*GenerateError)
	if !ok {
		t.Fatalf("Error type = %T, want *GenerateError", err)
	}

	if genErr.Code != "unauthorized" {
		t.Errorf("Error code = %v, want unauthorized", genErr.Code)
	}

	if genErr.Provider != "ModelScope" {
		t.Errorf("Provider = %v, want ModelScope", genErr.Provider)
	}
}

func TestModelScopeProvider_HandleErrorResponse_RateLimit(t *testing.T) {
	cfg := &config.Config{
		ImageAPIKey:  "test-key",
		ImageAPIBase: "https://mock.local",
	}

	p, _ := NewModelScopeProvider(cfg)
	p.client = newMockHTTPClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusTooManyRequests, `{"message": "rate limit exceeded"}`), nil
	})
	_, err := p.Generate(context.Background(), "test")

	if err == nil {
		t.Fatal("Expected error for rate limit")
	}

	genErr, ok := err.(*GenerateError)
	if !ok {
		t.Fatalf("Error type = %T, want *GenerateError", err)
	}

	if genErr.Code != "rate_limit" {
		t.Errorf("Error code = %v, want rate_limit", genErr.Code)
	}
}

func TestGetModelScopeSupportedModels(t *testing.T) {
	models := GetModelScopeSupportedModels()
	if len(models) == 0 {
		t.Error("GetModelScopeSupportedModels() returned empty list")
	}

	found := false
	for _, m := range models {
		if m == "Tongyi-MAI/Z-Image-Turbo" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tongyi-MAI/Z-Image-Turbo not found in supported models")
	}
}
