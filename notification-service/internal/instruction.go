package internal

type InstructionNotification struct {
	UserID        string `json:"user_id"`
	InstructionID string `json:"instruction_id"`
	FileID        string `json:"file_id"`
}
