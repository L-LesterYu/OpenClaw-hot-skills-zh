package wechat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var errCodePattern = regexp.MustCompile(`\b([0-9]{5})\b`)

// ExplainDraftAPIError converts known WeChat draft API errors into actionable hints.
func ExplainDraftAPIError(code int, msg string) string {
	base := fmt.Sprintf("wechat api error: %d - %s", code, msg)
	hint := draftAPIErrorHint(code)
	if hint == "" {
		return base
	}
	return base + "\nhint: " + hint
}

// ExplainDraftError inspects a raw error and enriches it when a known WeChat errcode is present.
func ExplainDraftError(err error) error {
	if err == nil {
		return nil
	}

	message := err.Error()
	matches := errCodePattern.FindStringSubmatch(message)
	if len(matches) < 2 {
		return err
	}
	code, convErr := strconv.Atoi(matches[1])
	if convErr != nil {
		return err
	}
	hint := draftAPIErrorHint(code)
	if hint == "" || strings.Contains(message, "\nhint: ") {
		return err
	}
	return fmt.Errorf("%s\nhint: %s", message, hint)
}

func draftAPIErrorHint(code int) string {
	switch code {
	case 45002:
		return "draft content exceeds the WeChat limit; shorten the article body or reduce oversized embedded HTML."
	case 45003:
		return "draft title exceeds the WeChat limit; shorten --title or frontmatter.title to 32 characters or fewer."
	case 45004:
		return "draft digest/description exceeds the WeChat limit; shorten --digest or frontmatter digest/summary/description to 128 characters or fewer."
	case 45005:
		return "a link field in the draft payload is invalid or exceeds the WeChat limit; check content_source_url or any generated external link fields."
	default:
		return ""
	}
}
