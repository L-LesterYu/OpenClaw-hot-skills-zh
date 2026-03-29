package action

import "testing"

func TestCompletedResult(t *testing.T) {
	result := CompletedResult(ActionConvert)
	if result.Status != StatusCompleted {
		t.Fatalf("expected completed status, got %q", result.Status)
	}
	if result.Action != ActionConvert {
		t.Fatalf("expected action %q, got %q", ActionConvert, result.Action)
	}
	if result.Retryable {
		t.Fatal("completed result should not be retryable")
	}
}

func TestActionRequiredResult(t *testing.T) {
	result := ActionRequiredResult(ActionAIRequest)
	if result.Status != StatusActionRequired {
		t.Fatalf("expected action_required status, got %q", result.Status)
	}
	if result.Action != ActionAIRequest {
		t.Fatalf("expected action %q, got %q", ActionAIRequest, result.Action)
	}
	if result.Retryable {
		t.Fatal("action-required result should not be retryable")
	}
}

func TestFailedResult(t *testing.T) {
	result := FailedResult(ActionDraftCreate, true)
	if result.Status != StatusFailed {
		t.Fatalf("expected failed status, got %q", result.Status)
	}
	if result.Action != ActionDraftCreate {
		t.Fatalf("expected action %q, got %q", ActionDraftCreate, result.Action)
	}
	if !result.Retryable {
		t.Fatal("failed result should preserve retryable=true")
	}
}
