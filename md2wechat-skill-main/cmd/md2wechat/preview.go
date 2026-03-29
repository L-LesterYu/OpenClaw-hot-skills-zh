package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
	"github.com/geekjourneyx/md2wechat-skill/internal/converter"
	inspectpkg "github.com/geekjourneyx/md2wechat-skill/internal/inspect"
	previewpkg "github.com/geekjourneyx/md2wechat-skill/internal/preview"
	"github.com/spf13/cobra"
)

var (
	previewMode           string
	previewTheme          string
	previewFontSize       string
	previewBackgroundType string
	previewTitle          string
	previewAuthor         string
	previewDigest         string
	previewCover          string
	previewUpload         bool
	previewDraft          bool
	previewOutput         string
)

var previewCmd = &cobra.Command{
	Use:   "preview <markdown_file>",
	Short: "Generate a standalone article preview HTML file",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		result, render, err := runPreview(args[0])
		if err != nil {
			return err
		}
		if jsonOutput {
			responseSuccessWith(codePreviewReady, "Preview ready", map[string]any{
				"output_file": previewOutput,
				"inspect":     result,
				"render":      render,
			})
			return nil
		}
		fmt.Printf("Preview written to %s\n", previewOutput)
		if render["fidelity"] != inspectpkg.PreviewFidelityExact {
			fmt.Printf("Preview fidelity: %s\n", render["fidelity"])
		}
		return nil
	},
}

func init() {
	previewCmd.Flags().StringVar(&previewMode, "mode", "api", "Preview context mode: api or ai")
	previewCmd.Flags().StringVar(&previewTheme, "theme", "default", "Theme name")
	previewCmd.Flags().StringVar(&previewFontSize, "font-size", "medium", "Font size: small/medium/large")
	previewCmd.Flags().StringVar(&previewBackgroundType, "background-type", "none", "Background type: default/grid/none")
	previewCmd.Flags().StringVar(&previewTitle, "title", "", "Override article title")
	previewCmd.Flags().StringVar(&previewAuthor, "author", "", "Override article author")
	previewCmd.Flags().StringVar(&previewDigest, "digest", "", "Override article digest")
	previewCmd.Flags().StringVar(&previewCover, "cover", "", "Cover image path to validate draft readiness")
	previewCmd.Flags().BoolVar(&previewUpload, "upload", false, "Evaluate upload readiness in preview")
	previewCmd.Flags().BoolVar(&previewDraft, "draft", false, "Evaluate draft readiness in preview")
	previewCmd.Flags().StringVarP(&previewOutput, "output", "o", "", "Write preview HTML to this file (default: temp file)")
}

func runPreview(markdownFile string) (*inspectpkg.Result, map[string]any, error) {
	markdown, err := os.ReadFile(markdownFile)
	if err != nil {
		return nil, nil, wrapCLIError(codeConvertReadFailed, err, fmt.Sprintf("read markdown file: %v", err))
	}

	result, err := runInspectWithInput(markdownFile, string(markdown), inspectpkg.Input{
		Mode:            previewMode,
		Theme:           previewTheme,
		FontSize:        previewFontSize,
		BackgroundType:  previewBackgroundType,
		TitleOverride:   previewTitle,
		AuthorOverride:  previewAuthor,
		DigestOverride:  previewDigest,
		CoverImagePath:  previewCover,
		UploadRequested: previewUpload,
		DraftRequested:  previewDraft,
	})
	if err != nil {
		return nil, nil, err
	}

	renderHTML, exactHTML, renderError, degradedMessage := buildPreviewRender(result)
	if !exactHTML {
		result.Readiness.PreviewFidelity = inspectpkg.PreviewFidelityDegraded
	}
	page, err := previewpkg.Render(previewpkg.PageData{
		Title:           result.Metadata.Title.Value,
		SourceFile:      result.SourceFile,
		Context:         result.Context,
		Metadata:        result.Metadata,
		Structure:       result.Structure,
		Checks:          result.Checks,
		Readiness:       result.Readiness,
		ArticleHTML:     template.HTML(renderHTML),
		BodyMarkdown:    result.Body,
		ExactHTML:       exactHTML,
		RenderError:     renderError,
		DegradedMessage: degradedMessage,
	})
	if err != nil {
		return nil, nil, wrapCLIError(codePreviewFailed, err, err.Error())
	}

	outputFile := previewOutput
	if outputFile == "" {
		tmp, err := os.CreateTemp("", "md2wechat-preview-*.html")
		if err != nil {
			return nil, nil, wrapCLIError(codePreviewFailed, err, err.Error())
		}
		outputFile = tmp.Name()
		if err := tmp.Close(); err != nil {
			return nil, nil, wrapCLIError(codePreviewFailed, err, err.Error())
		}
	}
	if err := os.WriteFile(outputFile, []byte(page), 0644); err != nil {
		return nil, nil, wrapCLIError(codePreviewFailed, err, err.Error())
	}
	previewOutput = outputFile

	return result, map[string]any{
		"mode":       result.Context.Mode,
		"theme":      result.Context.Theme,
		"fidelity":   result.Readiness.PreviewFidelity,
		"exact_html": exactHTML,
		"error":      renderError,
	}, nil
}

func buildPreviewRender(result *inspectpkg.Result) (html string, exact bool, renderError string, degradedMessage string) {
	if result == nil {
		return "", false, "inspect result is required", "Preview degraded: inspect result missing."
	}

	req := &converter.ConvertRequest{
		Markdown:       result.Body,
		Mode:           converter.ConvertMode(result.Context.Mode),
		Theme:          result.Context.Theme,
		FontSize:       result.Context.FontSize,
		BackgroundType: result.Context.BackgroundType,
		Metadata: converter.ArticleMetadata{
			Title:  result.Metadata.Title.Value,
			Author: result.Metadata.Author.Value,
			Digest: result.Metadata.Digest.Value,
		},
	}
	conv := newMarkdownConverter()
	if conv == nil {
		return "", false, "markdown converter is not available", "Preview degraded: converter not available."
	}
	converted := conv.Convert(req)
	if converted == nil {
		return "", false, "converter returned nil result", "Preview degraded: converter returned no result."
	}
	if converter.IsAIRequest(converted) || converted.Status == action.StatusActionRequired {
		return "", false, "", "Preview degraded: AI mode currently yields a prompt/request instead of final HTML."
	}
	if !converted.Success {
		return "", false, converted.Error, "Preview degraded: exact HTML could not be rendered in the current environment."
	}
	return converted.HTML, true, "", ""
}
