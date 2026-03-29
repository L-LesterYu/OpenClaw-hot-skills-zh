package writer

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
	"github.com/geekjourneyx/md2wechat-skill/internal/promptcatalog"
)

func newTestAssistant() *Assistant {
	style := &WriterStyle{
		Name:             "Dan Koe",
		EnglishName:      DefaultStyleName,
		WritingPrompt:    "Write clearly and directly.",
		CoverPrompt:      "Illustrate: {article_content}",
		CoverStyle:       "editorial",
		CoverMood:        "contemplative",
		CoverColorScheme: []string{"black", "white"},
		TitleFormulas:    []TitleFormula{{Examples: []string{"Example Title"}}},
	}

	return &Assistant{
		styleManager: &StyleManager{
			styles: map[string]*WriterStyle{
				DefaultStyleName: style,
			},
			initialized: true,
		},
		generator: NewGenerator(),
	}
}

func TestValidateWriteRequestNormalizesDefaults(t *testing.T) {
	asst := newTestAssistant()
	req := &WriteRequest{
		Input:       "这是一个足够长的测试输入，用来验证默认值规范化。",
		InputType:   InputType("unknown"),
		ArticleType: ArticleType("invalid"),
		Length:      Length("invalid"),
	}

	if err := asst.ValidateWriteRequest(req); err != nil {
		t.Fatalf("ValidateWriteRequest() error = %v", err)
	}
	if req.StyleName != DefaultStyleName {
		t.Fatalf("StyleName = %q", req.StyleName)
	}
	if req.InputType != InputTypeIdea {
		t.Fatalf("InputType = %q", req.InputType)
	}
	if req.ArticleType != ArticleTypeEssay {
		t.Fatalf("ArticleType = %q", req.ArticleType)
	}
	if req.Length != LengthMedium {
		t.Fatalf("Length = %q", req.Length)
	}
}

func TestWriteUsesNormalizedRequestAndReturnsAIRequest(t *testing.T) {
	asst := newTestAssistant()
	req := &WriteRequest{
		Input:       "这是一个足够长的测试输入，用来生成 AI 写作请求并验证默认 contract。",
		InputType:   InputType("unknown"),
		ArticleType: ArticleType("invalid"),
		Length:      Length("invalid"),
	}

	result := asst.Write(req)
	if !result.Success {
		t.Fatalf("Write() failed: %+v", result)
	}
	if !result.IsAIRequest {
		t.Fatalf("expected AI request result: %+v", result)
	}
	if result.Status != action.StatusActionRequired || result.Action != action.ActionWrite {
		t.Fatalf("unexpected typed state: %+v", result)
	}
	if result.Style == nil || result.Style.EnglishName != DefaultStyleName {
		t.Fatalf("unexpected style: %+v", result.Style)
	}
	if !strings.Contains(result.Prompt, "文章类型: essay") {
		t.Fatalf("prompt missing normalized article type: %q", result.Prompt)
	}
	if !strings.Contains(result.Prompt, "期望长度: medium") {
		t.Fatalf("prompt missing normalized length: %q", result.Prompt)
	}
}

func TestGeneratePromptReturnsPromptForAIRequests(t *testing.T) {
	asst := newTestAssistant()
	prompt := asst.GeneratePrompt(&WriteRequest{
		Input:     "这是一个足够长的测试输入，用来验证提示词生成。",
		StyleName: DefaultStyleName,
	})

	if prompt == "" {
		t.Fatal("GeneratePrompt() returned empty prompt")
	}
	if !strings.Contains(prompt, "请根据以上要求，生成符合该风格的文章。") {
		t.Fatalf("unexpected prompt: %q", prompt)
	}
}

func TestRefineReturnsExtractableAIRequest(t *testing.T) {
	promptcatalog.ResetDefaultCatalogForTests()

	asst := newTestAssistant()
	result := asst.Refine(&RefineRequest{
		Content:   "原始内容",
		StyleName: DefaultStyleName,
		Feedback:  "更直接一点",
	})

	if !result.Success {
		t.Fatalf("Refine() failed: %+v", result)
	}
	if !IsRefineRequest(result) {
		t.Fatalf("expected refine AI request: %+v", result)
	}
	if result.Status != action.StatusActionRequired || result.Action != action.ActionWrite {
		t.Fatalf("unexpected typed state: %+v", result)
	}
	if !strings.Contains(ExtractRefineRequest(result), "用户反馈") {
		t.Fatalf("unexpected refine prompt: %q", ExtractRefineRequest(result))
	}
	if !strings.Contains(ExtractRefineRequest(result), "Write clearly and directly.") {
		t.Fatalf("expected style prompt in refine prompt: %q", ExtractRefineRequest(result))
	}
}

func TestCompleteAIRequestMarksCompletedState(t *testing.T) {
	gen := NewGenerator()
	result := gen.Generate(&GenerateRequest{
		Style:       newTestAssistant().styleManager.styles[DefaultStyleName],
		UserInput:   "这是一个足够长的测试输入，用来生成 AI 写作请求并验证完成态。",
		InputType:   InputTypeIdea,
		ArticleType: ArticleTypeEssay,
		Length:      LengthMedium,
	})

	if !IsAIRequest(result) {
		t.Fatalf("expected AI request result: %+v", result)
	}

	completed := CompleteAIRequest("finished article", result)
	if completed.Status != action.StatusCompleted {
		t.Fatalf("Status = %q", completed.Status)
	}
	if completed.Action != action.ActionWrite {
		t.Fatalf("Action = %q", completed.Action)
	}
	if IsAIRequest(completed) {
		t.Fatalf("completed result should not require AI: %+v", completed)
	}
	if ExtractAIRequest(completed) != "" {
		t.Fatalf("ExtractAIRequest() = %q", ExtractAIRequest(completed))
	}
}

func TestCoverGeneratorBuildsPromptAndMetadata(t *testing.T) {
	asst := newTestAssistant()
	gen := NewCoverGenerator(asst.GetStyleManager())

	result, err := gen.GeneratePrompt(&GenerateCoverRequest{
		ArticleTitle:   "成长突破",
		ArticleContent: "我认为成长的本质，是持续打破旧的自己。",
		StyleName:      DefaultStyleName,
	})
	if err != nil {
		t.Fatalf("GeneratePrompt() error = %v", err)
	}
	if !result.Success {
		t.Fatalf("GeneratePrompt() failed: %+v", result)
	}
	if !strings.Contains(result.Prompt, "成长突破") {
		t.Fatalf("prompt = %q", result.Prompt)
	}
	if result.MetaData.CoreTheme == "" || result.MetaData.VisualAnchor == "" {
		t.Fatalf("metadata = %+v", result.MetaData)
	}
}

func TestParseStyleInputTrimsPrefixes(t *testing.T) {
	if got := ParseStyleInput(" --style=dan-koe "); got != DefaultStyleName {
		t.Fatalf("ParseStyleInput() = %q", got)
	}
	if got := ParseStyleInput("风格:dan-koe"); got != DefaultStyleName {
		t.Fatalf("ParseStyleInput() = %q", got)
	}
}
