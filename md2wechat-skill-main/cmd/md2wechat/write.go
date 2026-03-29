// Package main provides the md2wechat CLI tool
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/humanizer"
	"github.com/geekjourneyx/md2wechat-skill/internal/writer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// writeCmd 写作命令
var writeCmd = &cobra.Command{
	Use:   "write [input]",
	Short: "Writer Style Assistant - Write with creator styles",
	Long: `Assisted writing with customizable creator styles.

Default style: Dan Koe (profound, sharp, grounded)

Examples:
  # Interactive mode
  md2wechat write

  # Write from idea
  md2wechat write --style dan-koe

  # Refine existing content
  md2wechat write --style dan-koe --input-type fragment article.md

  # Generate with cover
  md2wechat write --style dan-koe --cover

  # Write with AI trace removal
  md2wechat write --style dan-koe --humanize
  md2wechat write --style dan-koe --humanize --humanize-intensity aggressive`,
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWrite(cmd, args)
	},
}

// write 命令参数
var (
	writeStyle             string
	writeInputType         string
	writeArticleType       string
	writeLength            string
	writeTitle             string
	writeOutput            string
	writeCover             bool
	writeCoverOnly         bool
	writeListStyles        bool
	writeStyleDetail       bool
	writeHumanize          bool
	writeHumanizeIntensity string
)

func init() {
	// 添加 flags
	writeCmd.Flags().StringVar(&writeStyle, "style", "dan-koe", "Writer style")
	writeCmd.Flags().StringVar(&writeInputType, "input-type", "idea", "Input type: idea/fragment/outline/title")
	writeCmd.Flags().StringVar(&writeArticleType, "article-type", "essay", "Article type: essay/commentary/story/tutorial/review")
	writeCmd.Flags().StringVar(&writeLength, "length", "medium", "Article length: short/medium/long")
	writeCmd.Flags().StringVar(&writeTitle, "title", "", "Article title")
	writeCmd.Flags().StringVarP(&writeOutput, "output", "o", "", "Output file")
	writeCmd.Flags().BoolVar(&writeCover, "cover", false, "Generate matching cover")
	writeCmd.Flags().BoolVar(&writeCoverOnly, "cover-only", false, "Generate cover only")
	writeCmd.Flags().BoolVar(&writeListStyles, "list", false, "List all available styles")
	writeCmd.Flags().BoolVar(&writeStyleDetail, "detail", false, "Show detailed style info")

	// Humanizer flags
	writeCmd.Flags().BoolVar(&writeHumanize, "humanize", false, "Enable AI trace removal")
	writeCmd.Flags().StringVar(&writeHumanizeIntensity, "humanize-intensity", "medium", "Humanize intensity: gentle/medium/aggressive")
}

// runWrite 执行写作命令
func runWrite(cmd *cobra.Command, args []string) error {
	// 处理列出风格
	if writeListStyles {
		return runListStyles()
	}

	// 获取输入内容
	input := ""
	if len(args) > 0 {
		// 从文件读取
		content, err := os.ReadFile(args[0])
		if err != nil {
			return wrapCLIError(codeWriteReadFailed, err, fmt.Sprintf("读取文件: %v", err))
		}
		input = string(content)

		// 如果没有明确指定输入类型，默认为 fragment
		if writeInputType == "idea" {
			writeInputType = "fragment"
		}
	} else {
		// 检查 stdin 是否有输入
		stdinContent, err := readStdin()
		if err == nil && stdinContent != "" {
			input = stdinContent
		}
	}

	// 如果没有输入，进入交互模式
	if input == "" {
		return runInteractiveWrite()
	}

	// 执行写作
	return executeWrite(input)
}

// runListStyles 列出所有风格
func runListStyles() error {
	asst := writer.NewAssistant()
	result := asst.ListStyles()

	if !result.Success {
		return newCLIError(codeWriteFailed, result.Error)
	}

	if writeStyleDetail {
		// 详细模式
		for _, style := range result.Styles {
			fmt.Println(writer.FormatStyleSummary(style))
			fmt.Println("---")
		}
	} else {
		// 简洁模式
		fmt.Println(writer.FormatStyleList(result.Styles))
	}

	return nil
}

// runInteractiveWrite 交互式写作模式
func runInteractiveWrite() error {
	fmt.Println("📝 Writer Style Assistant")
	fmt.Println()

	// 显示可用风格
	asst := writer.NewAssistant()
	styles := asst.GetAvailableStyles()

	fmt.Printf("可用风格 (%d 个):\n", len(styles))
	for _, styleName := range styles {
		style, _ := asst.GetStyleInfo(styleName)
		fmt.Printf("  - %s (%s)\n", style.Name, style.EnglishName)
	}
	fmt.Println()

	// 获取输入
	fmt.Print("请选择风格 [默认: dan-koe]: ")
	styleInput := readLine()
	if styleInput == "" {
		styleInput = "dan-koe"
	}

	fmt.Print("请输入你的观点或内容 (Ctrl+D 结束):\n")
	input := readMultiline()
	if strings.TrimSpace(input) == "" {
		return newCLIError(codeWriteInputInvalid, "输入不能为空")
	}

	// 构建请求
	req := &writer.WriteRequest{
		Input:     input,
		InputType: writer.GetInputTypeFromString(writeInputType),
		StyleName: styleInput,
		Length:    writer.GetLengthFromString(writeLength),
	}

	// 执行写作
	result := asst.Write(req)

	if result.IsAIRequest {
		// AI 模式：返回提示词
		output := map[string]interface{}{
			"mode":   "ai",
			"action": "ai_write_request",
			"style":  result.Style.Name,
			"prompt": result.Prompt,
		}

		// 如果启用了 humanizer，添加 humanizer 提示词
		if writeHumanize {
			h := humanizer.NewHumanizer()
			hReq := &humanizer.HumanizeRequest{
				Intensity:     humanizer.ParseIntensity(writeHumanizeIntensity),
				PreserveStyle: true,
				OriginalStyle: result.Style.EnglishName,
				ShowChanges:   true,
				IncludeScore:  true,
			}
			output["humanizer"] = map[string]interface{}{
				"enabled":         true,
				"intensity":       writeHumanizeIntensity,
				"prompt_template": h.BuildAIRequestForAI(hReq),
				"instruction":     "先生成文章，然后使用 humanizer prompt 去除 AI 痕迹",
			}
		}

		if writeCover {
			coverGen := writer.NewCoverGenerator(asst.GetStyleManager())
			coverResult, _ := coverGen.GeneratePrompt(&writer.GenerateCoverRequest{
				StyleName:      styleInput,
				ArticleContent: input,
			})
			if coverResult.Success {
				output["cover_prompt"] = coverResult.Prompt
			}
		}

		responseActionRequiredWith(codeWriteAIRequestReady, "AI write request prepared", output)
		return nil
	}

	if !result.Success {
		return newCLIError(codeWriteFailed, result.Error)
	}

	// 输出结果
	if writeOutput != "" {
		if err := os.WriteFile(writeOutput, []byte(result.Article), 0644); err != nil {
			return wrapCLIError(codeWriteFailed, err, fmt.Sprintf("保存文件: %v", err))
		}
		log.Info("article saved", zap.String("file", writeOutput))
	} else {
		fmt.Println("\n=== 生成文章 ===")
		fmt.Println(result.Article)
		fmt.Println("\n=== 金句 ===")
		for i, quote := range result.Quotes {
			fmt.Printf("%d. %s\n", i+1, quote)
		}
	}

	return nil
}

// executeWrite 执行写作
func executeWrite(input string) error {
	asst := writer.NewAssistant()

	req := &writer.WriteRequest{
		Input:     input,
		InputType: writer.GetInputTypeFromString(writeInputType),
		StyleName: writer.ParseStyleInput(writeStyle),
		Title:     writeTitle,
		Length:    writer.GetLengthFromString(writeLength),
	}

	result := asst.Write(req)

	if result.IsAIRequest {
		// AI 模式：返回提示词
		output := map[string]interface{}{
			"mode":   "ai",
			"action": "ai_write_request",
			"style":  result.Style.Name,
			"prompt": result.Prompt,
		}

		// 如果启用了 humanizer，添加 humanizer 提示词
		if writeHumanize {
			h := humanizer.NewHumanizer()
			hReq := &humanizer.HumanizeRequest{
				// Content 将在 AI 生成后填充
				Intensity:     humanizer.ParseIntensity(writeHumanizeIntensity),
				PreserveStyle: true, // 风格优先
				OriginalStyle: result.Style.EnglishName,
				ShowChanges:   true,
				IncludeScore:  true,
			}
			output["humanizer"] = map[string]interface{}{
				"enabled":         true,
				"intensity":       writeHumanizeIntensity,
				"prompt_template": h.BuildAIRequestForAI(hReq),
				"instruction":     "先生成文章，然后使用 humanizer prompt 去除 AI 痕迹",
			}
		}

		if writeCover || writeCoverOnly {
			coverGen := writer.NewCoverGenerator(asst.GetStyleManager())
			coverResult, err := coverGen.GeneratePrompt(&writer.GenerateCoverRequest{
				StyleName:      req.StyleName,
				ArticleTitle:   req.Title,
				ArticleContent: input,
			})
			if err == nil && coverResult.Success {
				output["cover_prompt"] = coverResult.Prompt
				output["cover_explanation"] = coverResult.Explanation
			}
		}

		responseActionRequiredWith(codeWriteAIRequestReady, "AI write request prepared", output)
		return nil
	}

	if !result.Success {
		return newCLIError(codeWriteFailed, result.Error)
	}

	// 只生成封面
	if writeCoverOnly {
		return generateCover(asst, req)
	}

	// 输出文章
	if writeOutput != "" {
		if err := os.WriteFile(writeOutput, []byte(result.Article), 0644); err != nil {
			return wrapCLIError(codeWriteFailed, err, fmt.Sprintf("保存文件: %v", err))
		}
		log.Info("article saved", zap.String("file", writeOutput))
	} else {
		fmt.Println("\n=== 生成文章 ===")
		fmt.Println(result.Article)
		fmt.Println("\n=== 金句 ===")
		for i, quote := range result.Quotes {
			fmt.Printf("%d. %s\n", i+1, quote)
		}
	}

	// 如果需要封面
	if writeCover {
		return generateCover(asst, req)
	}

	return nil
}

// generateCover 生成封面
func generateCover(asst *writer.Assistant, req *writer.WriteRequest) error {
	coverGen := writer.NewCoverGenerator(asst.GetStyleManager())

	coverReq := &writer.GenerateCoverRequest{
		StyleName:      req.StyleName,
		ArticleTitle:   req.Title,
		ArticleContent: req.Input,
	}

	result, err := coverGen.GeneratePrompt(coverReq)
	if err != nil {
		return wrapCLIError(codeWriteFailed, err, fmt.Sprintf("生成封面提示词: %v", err))
	}

	fmt.Println("\n=== 封面提示词 ===")
	fmt.Println(result.Prompt)

	if result.Explanation != "" {
		fmt.Println("\n---")
		fmt.Println("📖 隐喻说明:", result.Explanation)
	}

	return nil
}

// readLine 读取一行输入
func readLine() string {
	var line string
	if _, err := fmt.Scanln(&line); err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

// readMultiline 读取多行输入
func readMultiline() string {
	var lines []string
	for {
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// readStdin 读取标准输入，如果 stdin 为空或来自终端则返回空字符串
func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// 检查 stdin 是否来自管道或重定向（而非终端）
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil // 来自终端，无管道输入
	}

	// 读取所有 stdin 内容
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}
