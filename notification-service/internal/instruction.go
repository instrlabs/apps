package internal

type InstructionNotification struct {
	UserID              *string `json:"user_id,omitempty"`
	GuestID             *string `json:"guest_id,omitempty"`
	InstructionID       string  `json:"instruction_id"`
	InstructionDetailID string  `json:"instruction_detail_id"`
}
