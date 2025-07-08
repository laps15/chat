package requests

type SendMessageRequest struct {
	ChatID  int64  `form:"chat_id"`
	Message string `form:"message"`
}
