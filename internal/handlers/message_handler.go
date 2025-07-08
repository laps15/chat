package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/messages"
	"github.com/laps15/go-chat/internal/messages/components"
	"github.com/laps15/go-chat/internal/messages/requests"
	"github.com/laps15/go-chat/internal/users"
)

type MessageHandlers struct {
	authGroup       *echo.Group
	MessagesService *messages.MessagesService
}

func (h *MessageHandlers) SetAuthGroup(g *echo.Group) {
	h.authGroup = g
}

func (h *MessageHandlers) RegisterHandlers(e *echo.Echo) {

	chatGroup := h.authGroup.Group("/chat")
	chatGroup.GET("", h.handleHomePage)
	chatGroup.POST("/start", h.handleNewChat)
}

func (h *MessageHandlers) handleHomePage(c echo.Context) error {
	user := c.Get("user").(*users.User)
	chats, _ := h.MessagesService.GetChatsForUser(user)

	return components.Messages(components.MessageIndexProps{
		User:  user,
		Chats: chats,
	}).Render(context.Background(), c.Response().Writer)
}

func (h *MessageHandlers) handleNewChat(c echo.Context) error {
	chatReq := new(requests.NewMessageRequest)
	if err := c.Bind(chatReq); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input"})
	}

	user := c.Get("user").(*users.User)
	_, err := h.MessagesService.CreateMessage(user.ID, chatReq.ReceiverID, chatReq.Message)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create chat"})
	}

	return c.Redirect(http.StatusFound, "/chat")
}
