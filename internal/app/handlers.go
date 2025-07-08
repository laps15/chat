package app

import (
	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/auth"
)

type IHandlers interface {
	RegisterHandlers(*echo.Echo)
	SetAuthGroup(*echo.Group)
}

func RegisterHandlers(e *echo.Echo, handlers ...IHandlers) {
	auth_group := e.Group("")
	auth_group.Use(auth.GetSessionManager().AuthenticatedMiddleware)

	for _, handler := range handlers {
		handler.SetAuthGroup(auth_group)
		handler.RegisterHandlers(e)
	}
}
