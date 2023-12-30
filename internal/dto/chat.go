package dto

type MessageDto struct {
	Role  string  `json:"role"`
	Text  string  `json:"content"`
	RunID *string `json:"runId,omitempty"`
	Model string  `json:"model,omitempty"`
}
