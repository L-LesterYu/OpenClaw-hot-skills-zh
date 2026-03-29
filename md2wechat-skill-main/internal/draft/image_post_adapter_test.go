package draft

import (
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/publish"
	"github.com/geekjourneyx/md2wechat-skill/internal/wechat"
)

type fakeNewspicCreator struct {
	articles [][]wechat.NewspicArticle
	result   *wechat.CreateDraftResult
	err      error
}

func (f *fakeNewspicCreator) CreateNewspicDraft(articles []wechat.NewspicArticle) (*wechat.CreateDraftResult, error) {
	copied := append([]wechat.NewspicArticle(nil), articles...)
	f.articles = append(f.articles, copied)
	if f.err != nil {
		return nil, f.err
	}
	if f.result != nil {
		return f.result, nil
	}
	return &wechat.CreateDraftResult{MediaID: "draft-1"}, nil
}

func TestImagePostCreatorBuildsNewspicArticle(t *testing.T) {
	ws := &fakeNewspicCreator{
		result: &wechat.CreateDraftResult{MediaID: "draft-1"},
	}
	creator := &ImagePostCreator{creator: ws}

	result, err := creator.CreateImagePost(publish.ImagePostArtifact{
		Title:       "Title",
		Content:     "Body",
		OpenComment: true,
		FansOnly:    true,
		Assets: []publish.AssetRef{
			{MediaID: "img-1"},
			{MediaID: "img-2"},
		},
	})
	if err != nil {
		t.Fatalf("CreateImagePost() error = %v", err)
	}
	if result.MediaID != "draft-1" || result.ImageCount != 2 {
		t.Fatalf("result = %#v", result)
	}
	if len(result.UploadedIDs) != 2 || result.UploadedIDs[0] != "img-1" || result.UploadedIDs[1] != "img-2" {
		t.Fatalf("uploaded ids = %#v", result.UploadedIDs)
	}
	if len(ws.articles) != 1 || len(ws.articles[0]) != 1 {
		t.Fatalf("articles = %#v", ws.articles)
	}
	article := ws.articles[0][0]
	if article.Title != "Title" || article.Content != "Body" || article.ArticleType != string(ArticleTypeNewspic) {
		t.Fatalf("article = %#v", article)
	}
	if article.NeedOpenComment != 1 || article.OnlyFansCanComment != 1 {
		t.Fatalf("comment settings = %#v", article)
	}
	if len(article.ImageInfo.ImageList) != 2 || article.ImageInfo.ImageList[0].ImageMediaID != "img-1" {
		t.Fatalf("image info = %#v", article.ImageInfo)
	}
}
