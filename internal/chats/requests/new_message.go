package requests

type NewMessageRequest struct {
	ReceiverID int64  `json:"receiver_id" form:"receiver_id"`
	Message    string `json:"message" form:"message"`
}
