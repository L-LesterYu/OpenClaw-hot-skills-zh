package image

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

// ModelScopeProvider ModelScope 图片生成服务提供者
// ModelScope 使用异步 API 模式，需要轮询任务状态
type ModelScopeProvider struct {
	apiKey       string
	baseURL      string
	model        string
	size         string
	client       *http.Client
	pollInterval time.Duration // 轮询间隔，默认 5s
	maxPollTime  time.Duration // 最大轮询时间，默认 120s
}

// NewModelScopeProvider 创建 ModelScope Provider
func NewModelScopeProvider(cfg *config.Config) (*ModelScopeProvider, error) {
	model := cfg.ImageModel
	if model == "" {
		model = "Tongyi-MAI/Z-Image-Turbo" // 默认模型
	}

	size := cfg.ImageSize
	if size == "" {
		size = "1024x1024" // 默认尺寸
	}

	baseURL := cfg.ImageAPIBase
	if baseURL == "" {
		baseURL = "https://api-inference.modelscope.cn/" // 默认 API 地址
	}

	return &ModelScopeProvider{
		apiKey:       cfg.ImageAPIKey,
		baseURL:      baseURL,
		model:        model,
		size:         size,
		pollInterval: 5 * time.Second,   // 默认轮询间隔 5 秒
		maxPollTime:  120 * time.Second, // 默认最大轮询时间 120 秒
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Name 返回提供者名称
func (p *ModelScopeProvider) Name() string {
	return "ModelScope"
}

// Generate 生成图片（异步模式）
func (p *ModelScopeProvider) Generate(ctx context.Context, prompt string) (*GenerateResult, error) {
	// 1. 发起异步请求，获取 task_id
	taskID, err := p.createTask(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// 2. 轮询任务状态直到完成
	imageURL, err := p.pollTaskStatus(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return &GenerateResult{
		URL:           imageURL,
		RevisedPrompt: "", // ModelScope 不返回优化后的提示词
		Model:         p.model,
		Size:          p.size,
	}, nil
}

// parseSize 解析尺寸字符串 (如 "1024x1024") 为宽度和高度
func parseSize(size string) (width, height int, err error) {
	parts := strings.Split(size, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid size format: %s, expected WIDTHxHEIGHT", size)
	}
	width, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %s", parts[0])
	}
	height, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %s", parts[1])
	}
	return width, height, nil
}

// createTask 创建图片生成任务，返回 task_id
func (p *ModelScopeProvider) createTask(ctx context.Context, prompt string) (string, error) {
	width, height, err := parseSize(p.size)
	if err != nil {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "invalid_size",
			Message:  fmt.Sprintf("图片尺寸格式错误: %v", err),
			Hint:     "请使用 WIDTHxHEIGHT 格式，如 1024x1024",
			Original: err,
		}
	}

	reqBody := map[string]any{
		"model":  p.model,
		"prompt": prompt,
		"n":      1,
		"width":  width,
		"height": height,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "marshal_error",
			Message:  "请求构造失败",
			Original: err,
		}
	}

	url := p.baseURL + "/v1/images/generations"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "request_error",
			Message:  "创建请求失败",
			Original: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	// ModelScope 异步模式请求头
	req.Header.Set("X-ModelScope-Async-Mode", "true")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "network_error",
			Message:  "网络请求失败，请检查网络连接",
			Hint:     "确认网络连接正常，API 地址正确",
			Original: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 处理错误响应
	if resp.StatusCode != http.StatusOK {
		return "", p.handleErrorResponse(resp)
	}

	// 解析响应获取 task_id
	var result struct {
		TaskID string `json:"task_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "decode_error",
			Message:  "响应解析失败",
			Original: err,
		}
	}

	if result.TaskID == "" {
		return "", &GenerateError{
			Provider: p.Name(),
			Code:     "no_task_id",
			Message:  "未获取到任务 ID",
			Hint:     "API 返回格式可能已变更，请检查 ModelScope 文档",
		}
	}

	return result.TaskID, nil
}

// pollTaskStatus 轮询任务状态直到完成或超时
func (p *ModelScopeProvider) pollTaskStatus(ctx context.Context, taskID string) (string, error) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()
	timeout := time.After(p.maxPollTime)

	for {
		select {
		case <-ctx.Done():
			return "", &GenerateError{
				Provider: p.Name(),
				Code:     "canceled",
				Message:  "操作已取消",
				Original: ctx.Err(),
			}
		case <-timeout:
			return "", &GenerateError{
				Provider: p.Name(),
				Code:     "timeout",
				Message:  fmt.Sprintf("图片生成超时（超过 %v）", p.maxPollTime),
				Hint:     "图片生成时间较长，请稍后在任务列表中查看结果，或尝试简化提示词",
			}
		case <-ticker.C:
			status, url, err := p.getTaskStatus(ctx, taskID)
			if err != nil {
				return "", err
			}
			if status == "SUCCEED" {
				return url, nil
			}
			if status == "FAILED" {
				return "", &GenerateError{
					Provider: p.Name(),
					Code:     "task_failed",
					Message:  "图片生成任务失败",
					Hint:     "提示词可能不符合内容政策，请尝试修改提示词",
				}
			}
			// PENDING/RUNNING: 继续轮询
		}
	}
}

// getTaskStatus 获取任务状态
func (p *ModelScopeProvider) getTaskStatus(ctx context.Context, taskID string) (string, string, error) {
	url := p.baseURL + "/v1/tasks/" + taskID
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", &GenerateError{
			Provider: p.Name(),
			Code:     "request_error",
			Message:  "创建请求失败",
			Original: err,
		}
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("X-ModelScope-Task-Type", "image_generation")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", "", &GenerateError{
			Provider: p.Name(),
			Code:     "network_error",
			Message:  "查询任务状态失败，请检查网络连接",
			Hint:     "确认网络连接正常",
			Original: err,
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", "", p.handleErrorResponse(resp)
	}

	// 解析响应
	var result struct {
		TaskStatus   string   `json:"task_status"`
		OutputImages []string `json:"output_images,omitempty"`
		ErrorMessage string   `json:"error_message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", &GenerateError{
			Provider: p.Name(),
			Code:     "decode_error",
			Message:  "任务状态响应解析失败",
			Original: err,
		}
	}

	// 如果有错误信息，返回失败状态
	if result.ErrorMessage != "" {
		return "FAILED", "", nil
	}

	// 返回状态和图片 URL（如果有）
	imageURL := ""
	if result.TaskStatus == "SUCCEED" && len(result.OutputImages) > 0 {
		imageURL = result.OutputImages[0]
	}

	return result.TaskStatus, imageURL, nil
}

// handleErrorResponse 处理错误响应
func (p *ModelScopeProvider) handleErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp struct {
		Message string `json:"message"`
		Code    string `json:"code"`
		Data    any    `json:"data"`
	}

	_ = json.Unmarshal(body, &errResp)

	message := errResp.Message
	if message == "" {
		var rawErr struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &rawErr)
		message = rawErr.Error
	}
	if message == "" {
		message = string(body)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unauthorized",
			Message:  "ModelScope API Key 无效或已过期",
			Hint:     "请检查配置文件中的 api.image_key 是否正确，或前往 ModelScope 控制台获取新的 API Key",
			Original: fmt.Errorf("status 401: %s", string(body)),
		}
	case http.StatusTooManyRequests:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "rate_limit",
			Message:  "请求过于频繁，请稍后重试",
			Hint:     "ModelScope API 有速率限制，请等待一段时间后再试",
			Original: fmt.Errorf("status 429: %s", string(body)),
		}
	case http.StatusBadRequest:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "bad_request",
			Message:  fmt.Sprintf("请求参数错误: %s", message),
			Hint:     "请检查图片尺寸、模型名称等参数是否正确。支持的模型: Tongyi-MAI/Z-Image-Turbo",
			Original: fmt.Errorf("status 400: %s", string(body)),
		}
	case http.StatusPaymentRequired, http.StatusForbidden:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "payment_required",
			Message:  "ModelScope 账户余额不足或访问受限",
			Hint:     "请前往 ModelScope 控制台检查账户余额和 API 使用权限",
			Original: fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		}
	default:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unknown",
			Message:  fmt.Sprintf("ModelScope API 返回错误 (HTTP %d)", resp.StatusCode),
			Hint:     "请稍后重试，或访问 ModelScope 控制台查看服务状态",
			Original: fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		}
	}
}

// GetSupportedModels 返回 ModelScope 支持的模型列表
func GetModelScopeSupportedModels() []string {
	return []string{
		"Tongyi-MAI/Z-Image-Turbo",
	}
}
