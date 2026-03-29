package draft

import (
	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/publish"
	"github.com/geekjourneyx/md2wechat-skill/internal/wechat"
	"go.uber.org/zap"
)

// ImagePostCreator adapts the WeChat newspic draft API to the publish-layer contract.
type ImagePostCreator struct {
	creator newspicDraftCreator
}

type newspicDraftCreator interface {
	CreateNewspicDraft(articles []wechat.NewspicArticle) (*wechat.CreateDraftResult, error)
}

// NewImagePostCreator creates a publish-layer image-post adapter.
func NewImagePostCreator(cfg *config.Config, log *zap.Logger) *ImagePostCreator {
	return &ImagePostCreator{
		creator: wechat.NewService(cfg, log),
	}
}

// CreateImagePost publishes a canonical image-post artifact as a WeChat newspic draft.
func (c *ImagePostCreator) CreateImagePost(artifact publish.ImagePostArtifact) (*publish.ImagePostResult, error) {
	imageList := make([]wechat.NewspicImageItem, 0, len(artifact.Assets))
	uploadedIDs := make([]string, 0, len(artifact.Assets))
	for _, asset := range artifact.Assets {
		imageList = append(imageList, wechat.NewspicImageItem{
			ImageMediaID: asset.MediaID,
		})
		uploadedIDs = append(uploadedIDs, asset.MediaID)
	}

	article := wechat.NewspicArticle{
		Title:       artifact.Title,
		Content:     artifact.Content,
		ArticleType: string(ArticleTypeNewspic),
		ImageInfo: wechat.NewspicImageInfo{
			ImageList: imageList,
		},
	}
	if artifact.OpenComment {
		article.NeedOpenComment = 1
		if artifact.FansOnly {
			article.OnlyFansCanComment = 1
		}
	}

	result, err := c.creator.CreateNewspicDraft([]wechat.NewspicArticle{article})
	if err != nil {
		return nil, err
	}

	return &publish.ImagePostResult{
		MediaID:     result.MediaID,
		DraftURL:    result.DraftURL,
		ImageCount:  len(artifact.Assets),
		UploadedIDs: uploadedIDs,
	}, nil
}
