package publish

import "testing"

func TestAssetRefUploaded(t *testing.T) {
	tests := []struct {
		name  string
		asset AssetRef
		want  bool
	}{
		{
			name:  "not uploaded",
			asset: AssetRef{},
			want:  false,
		},
		{
			name: "uploaded by media id",
			asset: AssetRef{
				MediaID: "media-1",
			},
			want: true,
		},
		{
			name: "uploaded by public url",
			asset: AssetRef{
				PublicURL: "https://example.com/image.png",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.asset.Uploaded(); got != tt.want {
				t.Fatalf("Uploaded() = %v, want %v", got, tt.want)
			}
		})
	}
}
