package requests

type NewChatRequest struct {
	ReceiverID int64 `json:"receiver_id" form:"receiver_id"`
}
