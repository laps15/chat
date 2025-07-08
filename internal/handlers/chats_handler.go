package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/chats"
	"github.com/laps15/go-chat/internal/chats/components"
	"github.com/laps15/go-chat/internal/chats/requests"
	"github.com/laps15/go-chat/internal/users"
)

type ChatHandlers struct {
	authGroup    *echo.Group
	ChatsService *chats.ChatsService
}

func (h *ChatHandlers) SetAuthGroup(g *echo.Group) {
	h.authGroup = g
}

func (h *ChatHandlers) RegisterHandlers(e *echo.Echo) {
	chatGroup := h.authGroup.Group("/chat")
	chatGroup.GET("", h.handleHomePage)
	chatGroup.POST("/start", h.handleNewChat)
}

func (h *ChatHandlers) handleHomePage(c echo.Context) error {
	user := c.Get("user").(*users.User)
	chats, _ := h.ChatsService.GetChatsForUser(user)

	return components.Index(components.MessageIndexProps{
		User:  user,
		Chats: chats,
	}).Render(context.Background(), c.Response().Writer)
}

func (h *ChatHandlers) handleNewChat(c echo.Context) error {
	chatReq := new(requests.NewMessageRequest)
	if err := c.Bind(chatReq); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input"})
	}

	user := c.Get("user").(*users.User)
	_, err := h.ChatsService.CreateMessage(user.ID, chatReq.ReceiverID, chatReq.Message)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create chat"})
	}

	return c.Redirect(http.StatusFound, "/chat")
}
