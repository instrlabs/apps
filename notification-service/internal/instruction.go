package internal

type InstructionNotification struct {
	UserID              string `json:"user_id"`
	InstructionID       string `json:"instruction_id"`
	InstructionDetailID string `json:"instruction_detail_id"`
}
