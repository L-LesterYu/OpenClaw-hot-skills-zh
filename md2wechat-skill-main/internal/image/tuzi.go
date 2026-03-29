package image

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

// TuZiProvider TuZi 图片生成服务提供者
type TuZiProvider struct {
	apiKey  string
	baseURL string
	model   string
	size    string
	client  *http.Client
}

// NewTuZiProvider 创建 TuZi Provider
func NewTuZiProvider(cfg *config.Config) (*TuZiProvider, error) {
	model := cfg.ImageModel
	if model == "" {
		model = "doubao-seedream-4-5-251128" // 默认使用豆包 Seedream
	}

	size := cfg.ImageSize
	if size == "" {
		size = "2048x2048" // 默认正方形（TuZi 要求最小 3686400 像素）
	}

	return &TuZiProvider{
		apiKey:  cfg.ImageAPIKey,
		baseURL: cfg.ImageAPIBase,
		model:   model,
		size:    size,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// Name 返回提供者名称
func (p *TuZiProvider) Name() string {
	return "TuZi"
}

// Generate 生成图片
func (p *TuZiProvider) Generate(ctx context.Context, prompt string) (*GenerateResult, error) {
	// 构造请求
	reqBody := map[string]any{
		"model":           p.model,
		"prompt":          prompt,
		"n":               1,
		"size":            p.size,
		"response_format": "url",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "marshal_error",
			Message:  "请求构造失败",
			Original: err,
		}
	}

	// 创建请求
	url := p.baseURL + "/images/generations"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "request_error",
			Message:  "创建请求失败",
			Original: err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	// TuZi 特殊请求头
	req.Header.Set("HTTP-Referer", "https://md2wechat.cn")
	req.Header.Set("X-Title", "WeChat Markdown Editor")

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, &GenerateError{
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
		return nil, p.handleErrorResponse(resp)
	}

	// 解析响应 (OpenAI 兼容格式)
	var result struct {
		Data []struct {
			URL           string `json:"url"`
			RevisedPrompt string `json:"revised_prompt,omitempty"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "decode_error",
			Message:  "响应解析失败",
			Original: err,
		}
	}

	if len(result.Data) == 0 {
		return nil, &GenerateError{
			Provider: p.Name(),
			Code:     "no_image",
			Message:  "未生成图片",
			Hint:     "提示词可能不符合内容政策，请尝试修改提示词",
		}
	}

	return &GenerateResult{
		URL:           result.Data[0].URL,
		RevisedPrompt: result.Data[0].RevisedPrompt,
		Model:         p.model,
		Size:          p.size,
	}, nil
}

// handleErrorResponse 处理错误响应
func (p *TuZiProvider) handleErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	// 尝试解析 OpenAI 兼容错误格式
	_ = json.Unmarshal(body, &errResp)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unauthorized",
			Message:  "TuZi API Key 无效或已过期",
			Hint:     "请检查配置文件中的 api.image_key 是否正确，或前往 TuZi 控制台获取新的 API Key",
			Original: fmt.Errorf("status 401: %s", string(body)),
		}
	case http.StatusTooManyRequests:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "rate_limit",
			Message:  "请求过于频繁，请稍后重试",
			Hint:     "TuZi API 有速率限制，请等待一段时间后再试，或考虑升级套餐",
			Original: fmt.Errorf("status 429: %s", string(body)),
		}
	case http.StatusBadRequest:
		msg := errResp.Error.Message
		if msg == "" {
			msg = string(body)
		}
		return &GenerateError{
			Provider: p.Name(),
			Code:     "bad_request",
			Message:  fmt.Sprintf("请求参数错误: %s", msg),
			Hint:     "请检查图片尺寸、模型名称等参数是否正确。支持的模型: gemini-3-pro-image-preview, doubao-seedream-4-5-251128",
			Original: fmt.Errorf("status 400: %s", string(body)),
		}
	case http.StatusPaymentRequired, http.StatusForbidden:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "payment_required",
			Message:  "TuZi 账户余额不足或访问受限",
			Hint:     "请前往 TuZi 控制台充值或检查 API 使用权限",
			Original: fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		}
	default:
		return &GenerateError{
			Provider: p.Name(),
			Code:     "unknown",
			Message:  fmt.Sprintf("TuZi API 返回错误 (HTTP %d)", resp.StatusCode),
			Hint:     "请稍后重试，或访问 TuZi 控制台查看服务状态",
			Original: fmt.Errorf("status %d: %s", resp.StatusCode, string(body)),
		}
	}
}

// GetSupportedModels 返回 TuZi 支持的模型列表
func GetSupportedModels() []string {
	return []string{
		"gemini-3-pro-image-preview",
		"doubao-seedream-4-5-251128",
	}
}

// GetSupportedSizes 返回 TuZi 支持的尺寸列表
// 注意：TuZi 要求最小 3686400 像素
func GetSupportedSizes() []string {
	return []string{
		"2048x2048", // 1:1 正方形（默认，4.2M 像素）
		"1920x1920", // 1:1 正方形（最小要求，3.7M 像素）
		"2560x1440", // 16:9 横版（3.7M 像素）
		"1440x2560", // 9:16 竖版（3.7M 像素）
		"3072x2048", // 3:2 横版（6.3M 像素）
		"2048x3072", // 2:3 竖版（6.3M 像素）
		"3840x2160", // 16:9 超宽横版（8.3M 像素）
		"2160x3840", // 9:16 超高竖版（8.3M 像素）
	}
}
