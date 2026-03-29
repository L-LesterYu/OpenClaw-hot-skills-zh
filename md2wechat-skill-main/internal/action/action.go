package action

// SchemaVersion identifies the machine-readable response schema exposed by the CLI.
const SchemaVersion = "v1"

// Status describes the lifecycle state of a command or internal action result.
type Status string

const (
	StatusCompleted      Status = "completed"
	StatusActionRequired Status = "action_required"
	StatusFailed         Status = "failed"
)

// Common action names shared across packages.
const (
	ActionAIRequest     = "ai_request"
	ActionHumanize      = "humanize_request"
	ActionWrite         = "write_request"
	ActionConvert       = "convert_request"
	ActionImageUpload   = "image_upload"
	ActionImageGenerate = "image_generate"
	ActionDraftCreate   = "draft_create"
)

// Result captures the stable machine-readable lifecycle fields shared by commands
// and internal modules that need external continuation.
type Result struct {
	Status    Status `json:"status"`
	Action    string `json:"action,omitempty"`
	Retryable bool   `json:"retryable,omitempty"`
}

// CompletedResult returns the common success lifecycle fields.
func CompletedResult(action string) Result {
	return Result{
		Status: StatusCompleted,
		Action: action,
	}
}

// ActionRequiredResult returns the lifecycle fields for a request that needs an
// external actor to continue, such as an AI prompt.
func ActionRequiredResult(action string) Result {
	return Result{
		Status: StatusActionRequired,
		Action: action,
	}
}

// FailedResult returns the common failure lifecycle fields.
func FailedResult(action string, retryable bool) Result {
	return Result{
		Status:    StatusFailed,
		Action:    action,
		Retryable: retryable,
	}
}
