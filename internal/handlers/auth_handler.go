package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/auth"
	"github.com/laps15/go-chat/internal/auth/components"
	"github.com/laps15/go-chat/internal/auth/requests"
	"github.com/laps15/go-chat/internal/users"
)

type AuthHandlers struct {
	authGroup   *echo.Group
	AuthService *auth.AuthService
}

func (h *AuthHandlers) SetAuthGroup(g *echo.Group) {
	h.authGroup = g
}

func (h *AuthHandlers) RegisterHandlers(e *echo.Echo) {
	e.GET("/signin", h.handleLoginPage)
	e.POST("/signin", h.handleLogin)
	e.GET("/signup", h.handleSignupPage)
	e.POST("/signup", h.handleSignup)

	h.authGroup.GET("/logout", h.handleLogout)
}

func (h *AuthHandlers) handleLoginPage(c echo.Context) error {
	return components.Login().Render(context.Background(), c.Response().Writer)
}

func (h *AuthHandlers) handleLogin(c echo.Context) error {
	req := new(requests.AuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	user, err := h.AuthService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication failed"})
	}

	if err := auth.SetSessionValue(c, auth.SessIDKey, strconv.FormatInt(user.ID, 10)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set session ID"})
	}

	sv := time.Now().Add(12 * time.Hour).Unix()
	if err := auth.SetSessionValue(c, auth.SessExpKey, sv); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set session ExpKey"})
	}

	return c.Redirect(http.StatusFound, "/chat")
}

func (h *AuthHandlers) handleSignupPage(c echo.Context) error {
	return components.Signup().Render(context.Background(), c.Response().Writer)
}

func (h *AuthHandlers) handleSignup(c echo.Context) error {
	user := new(users.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	_, err := h.AuthService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user", "details": err.Error()})
	}

	return c.JSON(http.StatusCreated, "user created successfully")
}

func (h *AuthHandlers) handleLogout(c echo.Context) error {
	if err := auth.DeleteSession(c); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete session ID"})
	}

	return c.Redirect(http.StatusFound, "/signin")
}
