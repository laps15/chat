package main

import (
	"database/sql"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/app"
	"github.com/laps15/go-chat/internal/auth"
	"github.com/laps15/go-chat/internal/handlers"
	"github.com/laps15/go-chat/internal/messages"
	"github.com/laps15/go-chat/internal/users"
	_ "github.com/mattn/go-sqlite3"
)

type AppConfig struct {
	DatabaseFile       string
	SessionStoreSecret string
	SessionFieldName   string
	SessionMaxAge      int
}

func initDb(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	return db
}

func main() {
	appConfig := AppConfig{
		DatabaseFile:       "data/chat.db",
		SessionStoreSecret: "secret",
		SessionFieldName:   "session",
		SessionMaxAge:      (int)(12 * time.Hour.Seconds()),
	}

	e := echo.New()
	e.Debug = true

	e.Static("/assets", "assets")

	db := initDb(appConfig.DatabaseFile)
	defer db.Close()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(appConfig.SessionStoreSecret))))

	usersRepo := users.NewUsersRepository(db)
	messagesRepo := messages.NewMessagesRepository(db)

	usersService := users.NewUsersService(usersRepo)
	messagesService := messages.NewMessagesService(messagesRepo)
	authService := auth.NewAuthService(usersRepo)

	auth.InitSessionManager(&auth.SessionManager{
		SessionFieldName: appConfig.SessionFieldName,
		SessionMaxAge:    appConfig.SessionMaxAge,
		UsersService:     usersService,
	})
	authHandlers := &handlers.AuthHandlers{
		AuthService: authService,
	}
	messagesHandlers := &handlers.MessageHandlers{
		MessagesService: messagesService,
	}

	app.RegisterHandlers(e,
		authHandlers,
		messagesHandlers)

	// handlers.RegisterHandlers(e)

	e.Logger.Fatal(e.Start(":8090"))
}
