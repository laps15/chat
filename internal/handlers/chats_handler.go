package handlers

import (
	"context"
	"net/http"
	"strconv"

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
	chatGroup.POST("/start", h.handleNewChat)
	chatGroup.GET("", h.handleHomePage)
	chatGroup.POST("/send", h.handleSendMessage)
	chatGroup.GET("/:chatid", h.handleChat)
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
	chatReq := new(requests.NewChatRequest)
	if err := c.Bind(chatReq); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input"})
	}

	user := c.Get("user").(*users.User)

	existingChat, err := h.ChatsService.GetPrivateChatForUsers(user.ID, chatReq.ReceiverID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to check existing chat", "details": err.Error()})
	}

	if existingChat != nil {
		return c.Redirect(http.StatusFound, "/chat/"+strconv.FormatInt(existingChat.ID, 10))
	}

	_, err = h.ChatsService.CreatePrivateChat(user.ID, chatReq.ReceiverID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create chat"})
	}

	return c.Redirect(http.StatusFound, "/chat")
}

func (h *ChatHandlers) handleChat(c echo.Context) error {
	chatIdFromParams := c.Param("chatid")
	if chatIdFromParams == "" {
		return c.JSON(400, map[string]string{"error": "Chat ID is required"})
	}
	chatId, err := strconv.ParseInt(chatIdFromParams, 10, 64)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid Chat ID"})
	}

	user := c.Get("user").(*users.User)
	chat, err := h.ChatsService.GetChatById(chatId)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve chat"})
	}

	if chat == nil || chat.Participants == nil {
		return c.JSON(404, map[string]string{"error": "Chat not found"})
	}

	if _, ok := chat.Participants[user.ID]; !ok {
		return c.JSON(404, map[string]string{"error": "Chat not found"})
	}

	return components.Chat(components.ChatProps{
		User:     user,
		Chat:     chat,
		Messages: h.ChatsService.GetMessagesForChat(chat),
	}).Render(context.Background(), c.Response().Writer)
}

func (h *ChatHandlers) handleSendMessage(c echo.Context) error {
	user := c.Get("user").(*users.User)
	messageReq := new(requests.SendMessageRequest)
	if err := c.Bind(messageReq); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid input", "details": err.Error()})
	}

	chat, err := h.ChatsService.GetChatById(messageReq.ChatID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve chat"})
	}

	if chat == nil || chat.Participants == nil {
		return c.JSON(404, map[string]string{"error": "Chat not found"})
	}

	if _, ok := chat.Participants[user.ID]; !ok {
		return c.JSON(403, map[string]string{"error": "You are not a participant in this chat"})
	}

	message, err := h.ChatsService.SendMessage(user.ID, *chat, messageReq.Message)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to send message", "details": err.Error()})
	}

	return components.Message(components.MessageProps{
		User:    user,
		Message: *message,
	}).Render(context.Background(), c.Response().Writer)
}
