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
	authGroup := e.Group("")
	authGroup.Use(auth.GetSessionManager().AuthenticatedMiddleware)

	for _, handler := range handlers {
		handler.SetAuthGroup(authGroup)
		handler.RegisterHandlers(e)
	}
}
