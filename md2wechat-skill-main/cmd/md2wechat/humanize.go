package main

import (
	"fmt"
	"os"

	"github.com/geekjourneyx/md2wechat-skill/internal/humanizer"
	"github.com/spf13/cobra"
)

var (
	intensityFlag   string
	showChangesFlag bool
	outputFlag      string
)

// humanizeCmd - AI 写作去痕命令
var humanizeCmd = &cobra.Command{
	Use:   "humanize <file>",
	Short: "AI 写作去痕 - 去除文本中的 AI 生成痕迹",
	Long: `去除文本中的 AI 生成痕迹，使文章听起来更自然、更像人类书写。

基于 humanizer-zh 方法，检测并处理 24 种 AI 写作痕迹模式：
  • 内容模式：过度强调、夸大意义、宣传语言、模糊归因
  • 语言语法：AI 词汇、否定排比、三段式、同义词循环
  • 风格模式：破折号过度、粗体滥用、表情符号
  • 填充词回避：填充短语、过度限定、通用结论
  • 协作痕迹：对话式填充、知识截止免责声明

处理强度:
  gentle      - 温和处理，只修改明显的问题
  medium      - 中等强度 (默认)
  aggressive  - 激进处理，深度去除 AI 痕迹

示例:
  # 基本用法
  md2wechat humanize article.md

  # 指定强度
  md2wechat humanize article.md --intensity gentle

  # 显示修改对比和质量评分
  md2wechat humanize article.md --show-changes

  # 输出到文件
  md2wechat humanize article.md -o output.md

  # 与写作风格组合使用
  md2wechat write --style dan-koe --humanize
  md2wechat write --style dan-koe --humanize=aggressive`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHumanize(args[0])
	},
}

func runHumanize(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return wrapCLIError(codeHumanizeReadFailed, err, fmt.Sprintf("读取文件失败: %v", err))
	}

	req := &humanizer.HumanizeRequest{
		Content:       string(content),
		Intensity:     humanizer.ParseIntensity(intensityFlag),
		ShowChanges:   showChangesFlag,
		IncludeScore:  true,
		PreserveStyle: false,
	}

	h := humanizer.NewHumanizer()
	prompt := h.BuildAIRequestForAI(req)

	response := map[string]interface{}{
		"action": "humanize_request",
		"request": map[string]interface{}{
			"content":   req.Content,
			"intensity": req.Intensity.String(),
			"prompt":    prompt,
		},
	}

	if outputFlag != "" {
		if err := os.WriteFile(outputFlag, []byte(prompt), 0644); err != nil {
			return wrapCLIError(codeHumanizeWriteFailed, err, fmt.Sprintf("保存文件失败: %v", err))
		}
		response["output_file"] = outputFlag
	}

	responseActionRequiredWith(codeHumanizeRequestReady, "Humanize AI request prepared", response)
	return nil
}

// 从 AI 响应解析结果
func parseHumanizeResponse(aiResponse string, originalContent string, intensity humanizer.HumanizeIntensity) map[string]interface{} {
	h := humanizer.NewHumanizer()
	req := &humanizer.HumanizeRequest{
		Content:      originalContent,
		Intensity:    intensity,
		ShowChanges:  showChangesFlag,
		IncludeScore: true,
	}
	result := h.ParseAIResponse(aiResponse, req)

	// 构建输出
	output := map[string]interface{}{
		"success": result.Success,
		"content": result.Content,
	}

	if result.Error != "" {
		output["error"] = result.Error
	}
	if result.Report != "" {
		output["report"] = result.Report
	}
	if result.HasChanges() {
		output["changes_count"] = result.ChangeCount()
		output["changes"] = result.Changes
	}
	if result.Score != nil {
		output["score"] = map[string]interface{}{
			"total":        result.Score.Total,
			"directness":   result.Score.Directness,
			"rhythm":       result.Score.Rhythm,
			"trust":        result.Score.Trust,
			"authenticity": result.Score.Authenticity,
			"conciseness":  result.Score.Conciseness,
			"rating":       result.Score.Rating(),
		}
	}

	return output
}

func init() {
	humanizeCmd.Flags().StringVarP(&intensityFlag, "intensity", "i", "medium", "处理强度: gentle/medium/aggressive")
	humanizeCmd.Flags().BoolVarP(&showChangesFlag, "show-changes", "c", false, "显示修改对比和质量评分")
	humanizeCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "输出文件路径")

	// 添加强度别名
	humanizeCmd.Flags().Lookup("intensity").NoOptDefVal = "medium"
}
