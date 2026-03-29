package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/converter"
	"github.com/geekjourneyx/md2wechat-skill/internal/promptcatalog"
)

var (
	generateImageCmdSize     string
	generateImageCmdModel    string
	generateImageCmdPreset   string
	generateImageCmdArticle  string
	generateImageCmdTitle    string
	generateImageCmdSummary  string
	generateImageCmdKeywords string
	generateImageCmdStyle    string
	generateImageCmdAspect   string
)

type generateImageInput struct {
	RawPrompt         string
	Preset            string
	Article           string
	Title             string
	Summary           string
	Keywords          string
	Style             string
	Aspect            string
	Size              string
	Model             string
	RequiredArchetype string
}

type generateImageContext struct {
	Title     string
	Summary   string
	Keywords  string
	KeyPoints string
}

func runGenerateImage(args []string) error {
	input := generateImageInput{
		Preset:   generateImageCmdPreset,
		Article:  generateImageCmdArticle,
		Title:    generateImageCmdTitle,
		Summary:  generateImageCmdSummary,
		Keywords: generateImageCmdKeywords,
		Style:    generateImageCmdStyle,
		Aspect:   generateImageCmdAspect,
		Size:     generateImageCmdSize,
		Model:    generateImageCmdModel,
	}
	if len(args) > 0 {
		input.RawPrompt = args[0]
	}
	return runGenerateImageWithInput(input)
}

func runGeneratePresetImage(archetype, defaultPreset string, input generateImageInput) error {
	input.RequiredArchetype = archetype
	if strings.TrimSpace(input.Preset) == "" {
		input.Preset = defaultPreset
	}
	return runGenerateImageWithInput(input)
}

func runGenerateImageWithInput(input generateImageInput) error {
	if err := cfg.ValidateForImageGeneration(); err != nil {
		return wrapCLIError(codeConfigInvalid, err, err.Error())
	}

	prompt, err := resolveGenerateImagePrompt(input)
	if err != nil {
		return newCLIError(codeConfigInvalid, err.Error())
	}

	processor := resolveImageProcessor(input.Model)
	if input.Size != "" {
		result, err := processor.GenerateAndUploadWithSize(prompt, input.Size)
		if err != nil {
			return wrapCLIError(codeImageGenerateFailed, err, err.Error())
		}
		responseSuccess(result)
		return nil
	}

	result, err := processor.GenerateAndUpload(prompt)
	if err != nil {
		return wrapCLIError(codeImageGenerateFailed, err, err.Error())
	}
	responseSuccess(result)
	return nil
}

func resolveImageProcessor(model string) imageProcessor {
	model = strings.TrimSpace(model)
	if model == "" {
		return newImageProcessor()
	}

	cfgCopy := *cfg
	cfgCopy.ImageModel = model
	return newImageProcessorWithConfig(&cfgCopy)
}

func resolveGenerateImagePrompt(input generateImageInput) (string, error) {
	if strings.TrimSpace(input.Preset) == "" {
		if strings.TrimSpace(input.RawPrompt) == "" {
			return "", fmt.Errorf("generate_image requires a prompt or --preset")
		}
		return input.RawPrompt, nil
	}

	if strings.TrimSpace(input.RawPrompt) != "" {
		return "", fmt.Errorf("do not pass a raw prompt when --preset is used")
	}

	cat, err := promptcatalog.DefaultCatalog()
	if err != nil {
		return "", err
	}
	spec, err := cat.Get("image", input.Preset)
	if err != nil {
		return "", err
	}
	if input.RequiredArchetype != "" && !promptcatalog.SupportsUseCase(spec, input.RequiredArchetype) {
		return "", fmt.Errorf("preset %s is %s/%s, expected %s", spec.Name, spec.Archetype, spec.PrimaryUseCase, input.RequiredArchetype)
	}

	ctx, err := buildGenerateImageContext(input)
	if err != nil {
		return "", err
	}

	rendered, _, err := cat.Render("image", input.Preset, map[string]string{
		"ARTICLE_TITLE":   ctx.Title,
		"ARTICLE_SUMMARY": ctx.Summary,
		"KEYWORDS":        ctx.Keywords,
		"KEY_POINTS":      ctx.KeyPoints,
		"VISUAL_STYLE":    defaultString(input.Style, defaultVisualStyle(spec.Archetype)),
		"ASPECT_RATIO":    defaultString(input.Aspect, spec.DefaultAspectRatio, defaultAspectRatio(spec.Archetype)),
	})
	if err != nil {
		return "", err
	}
	return rendered, nil
}

func buildGenerateImageContext(input generateImageInput) (*generateImageContext, error) {
	ctx := &generateImageContext{
		Title:    strings.TrimSpace(input.Title),
		Summary:  strings.TrimSpace(input.Summary),
		Keywords: strings.TrimSpace(input.Keywords),
	}

	if input.Article != "" {
		markdown, err := os.ReadFile(input.Article)
		if err != nil {
			return nil, fmt.Errorf("read article: %w", err)
		}
		meta := converter.ParseArticleMetadata(string(markdown))
		if ctx.Title == "" {
			ctx.Title = strings.TrimSpace(meta.Title)
		}
		if ctx.Summary == "" {
			ctx.Summary = firstNonEmptyString(strings.TrimSpace(meta.Digest), deriveMarkdownSummary(string(markdown)))
		}
	}

	if ctx.Title == "" && ctx.Summary == "" {
		return nil, fmt.Errorf("--preset requires --article, --title, or --summary")
	}

	if ctx.Keywords == "" {
		ctx.Keywords = deriveKeywords(ctx.Title, ctx.Summary)
	}
	ctx.KeyPoints = firstNonEmptyString(ctx.Summary, ctx.Keywords, ctx.Title)
	return ctx, nil
}

func deriveMarkdownSummary(markdown string) string {
	normalized := strings.ReplaceAll(markdown, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	var body []string
	inFrontMatter := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i == 0 && trimmed == "---" {
			inFrontMatter = true
			continue
		}
		if inFrontMatter {
			if trimmed == "---" {
				inFrontMatter = false
			}
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "![") {
			continue
		}
		body = append(body, cleanedMarkdownLine(trimmed))
		if len(strings.Join(body, " ")) >= 140 {
			break
		}
	}

	return strings.TrimSpace(strings.Join(body, " "))
}

func cleanedMarkdownLine(line string) string {
	replacer := strings.NewReplacer("**", "", "__", "", "*", "", "`", "", ">", "", "-", "", "|", " ")
	return strings.Join(strings.Fields(replacer.Replace(line)), " ")
}

func deriveKeywords(values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		parts = append(parts, value)
	}
	return strings.Join(parts, "；")
}

func defaultVisualStyle(archetype string) string {
	switch strings.ToLower(strings.TrimSpace(archetype)) {
	case "infographic":
		return "clear information design"
	case "cover":
		return "editorial clean"
	default:
		return "clean visual style"
	}
}

func defaultAspectRatio(archetype string) string {
	switch strings.ToLower(strings.TrimSpace(archetype)) {
	case "infographic":
		return "3:4"
	case "cover":
		return "16:9"
	default:
		return "1:1"
	}
}

func defaultString(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
