package publish

import (
	htmlpkg "html"
	"regexp"
	"strings"
)

// InsertAssetPlaceholders rewrites matching img tags to stable placeholders.
func InsertAssetPlaceholders(html string, assets []AssetRef) string {
	result := html
	inserted := make(map[int]bool, len(assets))
	for _, asset := range assets {
		if asset.Placeholder == "" {
			continue
		}

		escapedSource := htmlpkg.EscapeString(asset.Source)
		candidates := []string{asset.Source}
		if escapedSource != asset.Source {
			candidates = append(candidates, escapedSource)
		}

		for _, candidate := range candidates {
			doubleQuoted := regexp.MustCompile(`(?i)<img[^>]*src="` + regexp.QuoteMeta(candidate) + `"[^>]*>`)
			singleQuoted := regexp.MustCompile(`(?i)<img[^>]*src='` + regexp.QuoteMeta(candidate) + `'[^>]*>`)
			if doubleQuoted.MatchString(result) || singleQuoted.MatchString(result) {
				inserted[asset.Index] = true
			}
			result = doubleQuoted.ReplaceAllString(result, asset.Placeholder)
			result = singleQuoted.ReplaceAllString(result, asset.Placeholder)
		}
	}

	for _, asset := range assets {
		if inserted[asset.Index] || asset.Placeholder == "" {
			continue
		}

		imgTagPattern := regexp.MustCompile(`(?i)<img\b[^>]*>`)
		result = imgTagPattern.ReplaceAllStringFunc(result, func(tag string) string {
			if inserted[asset.Index] {
				return tag
			}
			inserted[asset.Index] = true
			return asset.Placeholder
		})
	}

	return result
}

// ReplaceAssetPlaceholders rewrites placeholders or original src values to published URLs.
func ReplaceAssetPlaceholders(html string, assets []AssetRef) string {
	result := html
	for _, asset := range assets {
		if asset.PublicURL == "" {
			continue
		}
		if asset.Placeholder != "" {
			imgTag := `<img src="` + asset.PublicURL + `" style="max-width:100%;height:auto;display:block;margin:20px auto;" />`
			result = strings.ReplaceAll(result, asset.Placeholder, imgTag)
		}

		escapedSource := htmlpkg.EscapeString(asset.Source)
		replacements := [][2]string{
			{`src="` + asset.Source + `"`, `src="` + asset.PublicURL + `"`},
			{`src='` + asset.Source + `'`, `src='` + asset.PublicURL + `'`},
		}
		if escapedSource != asset.Source {
			replacements = append(replacements,
				[2]string{`src="` + escapedSource + `"`, `src="` + asset.PublicURL + `"`},
				[2]string{`src='` + escapedSource + `'`, `src='` + asset.PublicURL + `'`},
			)
		}
		for _, replacement := range replacements {
			result = strings.ReplaceAll(result, replacement[0], replacement[1])
		}
	}
	return result
}
