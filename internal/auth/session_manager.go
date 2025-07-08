package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/laps15/go-chat/internal/users"
)

type CookieOptions struct {
	Name   string
	Secret string
	MaxAge int
}

type SessionManager struct {
	SessionFieldName string
	SessionMaxAge    int
	UsersService     *users.UsersService
}

const (
	SessIDKey  = "uuid"
	SessExpKey = "expireAt"
)

var (
	mgr *SessionManager
)

func InitSessionManager(newMgr *SessionManager) {
	mgr = newMgr
}

func GetSessionManager() *SessionManager {
	return mgr
}

func GetSession(c echo.Context) *sessions.Session {
	s, err := session.Get(mgr.SessionFieldName, c)

	if err != nil {
		c.Logger().Errorf("Failed to get session: %v", err)
		return nil
	}

	if s.IsNew {
		s.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   mgr.SessionMaxAge,
			HttpOnly: true,
			Secure:   true,
		}
		s.Values = make(map[interface{}]interface{})
	}

	return s
}

func GetSessionValue(c echo.Context, key string) interface{} {
	s := GetSession(c)
	return s.Values[key]
}

func SetSessionValue(c echo.Context, key string, value interface{}) error {
	s := GetSession(c)
	s.Values[key] = value

	if err := s.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func DeleteSessionValue(c echo.Context, key string) error {
	s := GetSession(c)
	s.Options.MaxAge = -1
	delete(s.Values, key)

	if err := s.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func DeleteSession(c echo.Context) error {
	s := GetSession(c)
	s.Options.MaxAge = -1
	s.Values = make(map[interface{}]interface{})

	if err := s.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}

func GetAuthenticatedUserId(c echo.Context) (string, error) {
	s := GetSession(c)
	if s == nil {
		return "", fmt.Errorf("SessionNotFound")
	}

	uid, ok := s.Values[SessIDKey]
	if !ok || uid.(string) == "" {
		c.Logger().Info("Session ID not found or empty")
		return "", fmt.Errorf("SessionIdNotFound")
	}

	c.Logger().Infof("Session ID found: %s", uid)
	expireAt, ok := s.Values[SessExpKey]
	if !ok {
		c.Logger().Info("Session ExpireAt not found")
		return "", fmt.Errorf("SessionExpireDateNotFound")
	}

	c.Logger().Infof("Session ExpireAt found: %v", expireAt)
	if v := time.Unix(expireAt.(int64), 0); v.Before(time.Now()) {
		return "", fmt.Errorf("SessionExpired")
	}

	return uid.(string), nil
}

func (mgr *SessionManager) AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		strUserId, err := GetAuthenticatedUserId(c)
		if err != nil {
			return c.Redirect(http.StatusFound, "/signin")
		}

		userID, err := strconv.ParseInt(strUserId, 10, 64)
		if err != nil {
			DeleteSession(c)
			return c.Redirect(http.StatusFound, "/signin")
		}

		userFromDB, err := mgr.UsersService.GetUserByID(userID)
		if err != nil {
			DeleteSession(c)
			return c.Redirect(http.StatusFound, "/signin")
		}

		c.Set("user", userFromDB)

		return next(c)
	}
}
