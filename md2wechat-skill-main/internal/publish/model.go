package publish

// AssetKind identifies how an asset entered the publish pipeline.
type AssetKind string

const (
	AssetKindLocal  AssetKind = "local"
	AssetKindRemote AssetKind = "remote"
	AssetKindAI     AssetKind = "ai"
)

// Metadata is the canonical article metadata carried across the publish pipeline.
type Metadata struct {
	Title            string `json:"title,omitempty"`
	Author           string `json:"author,omitempty"`
	Digest           string `json:"digest,omitempty"`
	ContentSourceURL string `json:"content_source_url,omitempty"`
}

// AssetRef is the canonical asset model for the publish pipeline.
type AssetRef struct {
	Index          int       `json:"index"`
	Kind           AssetKind `json:"kind"`
	Source         string    `json:"source"`
	ResolvedSource string    `json:"resolved_source,omitempty"`
	Prompt         string    `json:"prompt,omitempty"`
	Placeholder    string    `json:"placeholder,omitempty"`
	MediaID        string    `json:"media_id,omitempty"`
	PublicURL      string    `json:"public_url,omitempty"`
}

// Uploaded reports whether the asset has already been published to the target backend.
func (a AssetRef) Uploaded() bool {
	return a.MediaID != "" || a.PublicURL != ""
}

// ArticleSource is the normalized source document entering the publish pipeline.
type ArticleSource struct {
	Path     string     `json:"path,omitempty"`
	Markdown string     `json:"markdown"`
	Metadata Metadata   `json:"metadata,omitempty"`
	Assets   []AssetRef `json:"assets,omitempty"`
}

// PublishIntent captures the user-facing publish switches without binding them to a CLI command.
type PublishIntent struct {
	Mode        string `json:"mode,omitempty"`
	Preview     bool   `json:"preview,omitempty"`
	Upload      bool   `json:"upload,omitempty"`
	CreateDraft bool   `json:"create_draft,omitempty"`
	SaveDraft   bool   `json:"save_draft,omitempty"`
}

// Artifact is the canonical publish output shared by future orchestrators and adapters.
type Artifact struct {
	HTML         string     `json:"html,omitempty"`
	OutputFile   string     `json:"output_file,omitempty"`
	Metadata     Metadata   `json:"metadata,omitempty"`
	Assets       []AssetRef `json:"assets,omitempty"`
	CoverMediaID string     `json:"cover_media_id,omitempty"`
	DraftMediaID string     `json:"draft_media_id,omitempty"`
	DraftURL     string     `json:"draft_url,omitempty"`
}

// DraftResult is the canonical draft creation result returned by publish adapters.
type DraftResult struct {
	MediaID  string `json:"media_id"`
	DraftURL string `json:"draft_url,omitempty"`
}

// ImagePostSource is the normalized image-post request entering the publish pipeline.
type ImagePostSource struct {
	Title       string     `json:"title"`
	Content     string     `json:"content,omitempty"`
	Assets      []AssetRef `json:"assets"`
	OpenComment bool       `json:"open_comment,omitempty"`
	FansOnly    bool       `json:"fans_only,omitempty"`
}

// ImagePostArtifact is the canonical image-post artifact shared with publish adapters.
type ImagePostArtifact struct {
	Title       string     `json:"title"`
	Content     string     `json:"content,omitempty"`
	Assets      []AssetRef `json:"assets"`
	OpenComment bool       `json:"open_comment,omitempty"`
	FansOnly    bool       `json:"fans_only,omitempty"`
}

// ImagePostResult is the canonical image-post creation result.
type ImagePostResult struct {
	MediaID     string   `json:"media_id"`
	DraftURL    string   `json:"draft_url,omitempty"`
	ImageCount  int      `json:"image_count"`
	UploadedIDs []string `json:"uploaded_ids"`
}
