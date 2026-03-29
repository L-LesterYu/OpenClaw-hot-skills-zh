package main

import (
	"fmt"
	"os"

	"github.com/geekjourneyx/md2wechat-skill/internal/publish"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var testHTMLCmd = &cobra.Command{
	Use:   "test-draft <html_file> <cover_image>",
	Short: "Test creating WeChat draft from HTML file",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := runTestDraft(args[0], args[1])
		if err != nil {
			return err
		}
		responseSuccessWith(codeTestDraftCreated, "Draft created successfully", response)
		return nil
	},
}

func runTestDraft(htmlFile, coverImage string) (map[string]any, error) {
	if err := cfg.ValidateForWeChat(); err != nil {
		return nil, wrapCLIError(codeConfigInvalid, err, err.Error())
	}

	html, err := os.ReadFile(htmlFile)
	if err != nil {
		return nil, wrapCLIError(codeTestDraftReadFailed, err, fmt.Sprintf("read HTML file: %v", err))
	}

	log.Info("testing draft creation",
		zap.Int("html_length", len(html)),
		zap.String("cover", coverImage))

	log.Info("uploading cover image", zap.String("path", coverImage))
	coverMediaID, err := uploadCoverImageFn(coverImage)
	if err != nil {
		return nil, wrapCLIError(codeTestDraftCoverFailed, err, fmt.Sprintf("upload cover: %v", err))
	}
	log.Info("cover uploaded", zap.String("media_id", maskMediaID(coverMediaID)))

	svc := newDraftCreator()
	result, err := svc.CreateDraft(publish.Artifact{
		HTML: string(html),
		Metadata: publish.Metadata{
			Title:  "AI生成测试文章",
			Digest: "这是AI生成的微信公众号文章测试",
		},
		CoverMediaID: coverMediaID,
	})
	if err != nil {
		return nil, wrapCLIError(codeTestDraftCreateFailed, err, fmt.Sprintf("create draft: %v", err))
	}

	response := map[string]any{
		"media_id": result.MediaID,
		"message":  "Draft created successfully! You can check it in WeChat backend.",
	}
	if result.DraftURL != "" {
		response["draft_url"] = result.DraftURL
	}

	return response, nil
}
