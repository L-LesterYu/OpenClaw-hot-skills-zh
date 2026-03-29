package publish

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/image"
)

// AssetPipeline isolates publish-time asset resolution, upload, and HTML replacement.
type AssetPipeline struct {
	processor AssetProcessor
}

// NewAssetPipeline creates a publish asset pipeline.
func NewAssetPipeline(processor AssetProcessor) *AssetPipeline {
	return &AssetPipeline{processor: processor}
}

// ProcessInput is the normalized asset-processing request.
type ProcessInput struct {
	HTML        string
	Assets      []AssetRef
	MarkdownDir string
}

// ProcessOutput is the normalized asset-processing result.
type ProcessOutput struct {
	HTML   string
	Assets []AssetRef
}

// Process uploads or generates assets and rewrites the HTML to the published URLs.
func (p *AssetPipeline) Process(input *ProcessInput) (*ProcessOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("asset process input is required")
	}
	if len(input.Assets) == 0 {
		return &ProcessOutput{
			HTML:   input.HTML,
			Assets: input.Assets,
		}, nil
	}
	if p.processor == nil {
		return nil, fmt.Errorf("asset processor is required")
	}

	output := &ProcessOutput{
		HTML:   InsertAssetPlaceholders(input.HTML, input.Assets),
		Assets: append([]AssetRef(nil), input.Assets...),
	}

	var failed []string
	for i, asset := range output.Assets {
		var uploadResult *image.UploadResult
		var err error

		switch asset.Kind {
		case AssetKindLocal:
			localPath := asset.ResolvedSource
			if localPath == "" {
				localPath = asset.Source
			}
			if !filepath.IsAbs(localPath) && input.MarkdownDir != "" {
				localPath = filepath.Join(input.MarkdownDir, localPath)
			}
			uploadResult, err = p.processor.UploadLocalImage(localPath)
			output.Assets[i].ResolvedSource = localPath
		case AssetKindRemote:
			uploadResult, err = p.processor.DownloadAndUpload(asset.Source)
		case AssetKindAI:
			var genResult *image.GenerateAndUploadResult
			genResult, err = p.processor.GenerateAndUpload(asset.Prompt)
			if err == nil {
				uploadResult = &image.UploadResult{
					MediaID:   genResult.MediaID,
					WechatURL: genResult.WechatURL,
				}
			}
		default:
			err = fmt.Errorf("unsupported asset kind: %s", asset.Kind)
		}

		if err != nil {
			failed = append(failed, fmt.Sprintf("%d:%s", i, err.Error()))
			continue
		}

		output.Assets[i].MediaID = uploadResult.MediaID
		output.Assets[i].PublicURL = uploadResult.WechatURL
	}

	output.HTML = ReplaceAssetPlaceholders(output.HTML, output.Assets)
	if len(failed) > 0 {
		return nil, fmt.Errorf("image processing failed for %d image(s): %s", len(failed), strings.Join(failed, "; "))
	}

	return output, nil
}
