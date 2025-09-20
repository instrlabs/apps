package internal

type InstructionNotification struct {
	InstructionID string `json:"instruction_id"`
	FileID        string `json:"file_id"`
	FileStatus    string `json:"file_status"`
}
