package publish

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/image"
)

type fakeImagePostCreator struct {
	artifacts []ImagePostArtifact
	result    *ImagePostResult
	err       error
}

func (f *fakeImagePostCreator) CreateImagePost(artifact ImagePostArtifact) (*ImagePostResult, error) {
	f.artifacts = append(f.artifacts, artifact)
	if f.err != nil {
		return nil, f.err
	}
	if f.result != nil {
		return f.result, nil
	}
	return &ImagePostResult{
		MediaID:     "draft-1",
		ImageCount:  len(artifact.Assets),
		UploadedIDs: []string{"img-1"},
	}, nil
}

func TestExtractAssetsFromMarkdownUsesSharedParser(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "article.md")
	content := []byte(`
![cover](./a.png "cover")
![nested](images/b.png)
![parent](../c.png)
![absolute](/tmp/d.png)
![remote](https://example.com/e.png)
![angle](<images/my cat.png>)
`)
	if err := os.WriteFile(mdPath, content, 0600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	got, err := ExtractAssetsFromMarkdown(mdPath)
	if err != nil {
		t.Fatalf("ExtractAssetsFromMarkdown() error = %v", err)
	}

	want := []AssetRef{
		{Index: 0, Kind: AssetKindLocal, Source: "./a.png", ResolvedSource: filepath.Join(dir, "./a.png")},
		{Index: 1, Kind: AssetKindLocal, Source: "images/b.png", ResolvedSource: filepath.Join(dir, "images/b.png")},
		{Index: 2, Kind: AssetKindLocal, Source: "../c.png", ResolvedSource: filepath.Join(dir, "../c.png")},
		{Index: 3, Kind: AssetKindLocal, Source: "/tmp/d.png", ResolvedSource: "/tmp/d.png"},
		{Index: 4, Kind: AssetKindRemote, Source: "https://example.com/e.png"},
		{Index: 5, Kind: AssetKindLocal, Source: "images/my cat.png", ResolvedSource: filepath.Join(dir, "images/my cat.png")},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractAssetsFromMarkdown() = %#v, want %#v", got, want)
	}
}

func TestImagePostServicePreviewValidatesAndIncludesFileDetails(t *testing.T) {
	svc := NewImagePostService(nil, nil)

	if _, err := svc.PreviewImagePost(&ImagePostInput{Title: "Title"}); err == nil {
		t.Fatal("expected no images error")
	}

	images := make([]string, 21)
	for i := range images {
		images[i] = "image.jpg"
	}
	if _, err := svc.PreviewImagePost(&ImagePostInput{Title: "Title", Images: images}); err == nil {
		t.Fatal("expected too many images error")
	}

	dir := t.TempDir()
	imagePath := filepath.Join(dir, "image.jpg")
	if err := os.WriteFile(imagePath, []byte("12345"), 0600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	preview, err := svc.PreviewImagePost(&ImagePostInput{
		Title:  "Preview",
		Images: []string{imagePath},
	})
	if err != nil {
		t.Fatalf("PreviewImagePost() error = %v", err)
	}
	if preview.ImageCount != 1 || len(preview.Images) != 1 {
		t.Fatalf("preview = %#v", preview)
	}
	if !preview.Images[0].Exists || preview.Images[0].Size != int64(5) {
		t.Fatalf("preview image = %#v", preview.Images[0])
	}
}

func TestImagePostServiceCreateUploadsAssetsAndDispatchesArtifact(t *testing.T) {
	dir := t.TempDir()
	imagePath := filepath.Join(dir, "image.jpg")
	if err := os.WriteFile(imagePath, []byte("12345"), 0600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	assets := &fakeAssetProcessor{
		localResults: map[string]*image.UploadResult{
			imagePath: {MediaID: "img-1", WechatURL: "https://wechat.local/image"},
		},
	}
	creator := &fakeImagePostCreator{
		result: &ImagePostResult{
			MediaID:     "draft-1",
			ImageCount:  1,
			UploadedIDs: []string{"img-1"},
		},
	}
	svc := NewImagePostService(assets, creator)

	result, err := svc.CreateImagePost(&ImagePostInput{
		Title:       "Title",
		Content:     "Body",
		Images:      []string{imagePath},
		OpenComment: true,
		FansOnly:    true,
	})
	if err != nil {
		t.Fatalf("CreateImagePost() error = %v", err)
	}
	if result.MediaID != "draft-1" || result.ImageCount != 1 {
		t.Fatalf("result = %#v", result)
	}
	if len(assets.localCalls) != 1 || assets.localCalls[0] != imagePath {
		t.Fatalf("local calls = %#v", assets.localCalls)
	}
	if len(creator.artifacts) != 1 {
		t.Fatalf("artifacts = %#v", creator.artifacts)
	}
	artifact := creator.artifacts[0]
	if artifact.Title != "Title" || artifact.Content != "Body" || !artifact.OpenComment || !artifact.FansOnly {
		t.Fatalf("artifact = %#v", artifact)
	}
	if len(artifact.Assets) != 1 || artifact.Assets[0].MediaID != "img-1" {
		t.Fatalf("artifact assets = %#v", artifact.Assets)
	}
}

func TestImagePostServiceCreateUsesResolvedMarkdownAssetPaths(t *testing.T) {
	dir := t.TempDir()
	mdDir := filepath.Join(dir, "posts")
	if err := os.MkdirAll(filepath.Join(mdDir, "images"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	mdPath := filepath.Join(mdDir, "article.md")
	imagePath := filepath.Join(mdDir, "images", "a.png")
	if err := os.WriteFile(mdPath, []byte("![x](images/a.png)\n"), 0600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}
	if err := os.WriteFile(imagePath, []byte("img"), 0600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	assets := &fakeAssetProcessor{
		localResults: map[string]*image.UploadResult{
			imagePath: {MediaID: "img-1", WechatURL: "https://wechat.local/image"},
		},
	}
	creator := &fakeImagePostCreator{}
	svc := NewImagePostService(assets, creator)

	if _, err := svc.CreateImagePost(&ImagePostInput{
		Title:        "Title",
		FromMarkdown: mdPath,
	}); err != nil {
		t.Fatalf("CreateImagePost() error = %v", err)
	}

	if len(assets.localCalls) != 1 || assets.localCalls[0] != imagePath {
		t.Fatalf("local calls = %#v, want %q", assets.localCalls, imagePath)
	}
}
