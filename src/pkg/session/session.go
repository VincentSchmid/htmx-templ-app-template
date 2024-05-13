package session

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store *sessions.CookieStore

const (
	userSessionKey   = "user_session"
	dbAccessTokenKey = "db_accessToken"
)

func Init(sessionConfig *appconfig.Session) error {
	authKeyBase64 := sessionConfig.AuthKey
	encryptionKeyBase64 := sessionConfig.EncryptionKey

	// Decode the base64 auth key
	authKey, err := base64.StdEncoding.DecodeString(authKeyBase64)
	if err != nil {
		return fmt.Errorf("failed to decode auth key: %w", err)
	}

	// Decode the base64 encryption key
	encryptionKey, err := base64.StdEncoding.DecodeString(encryptionKeyBase64)
	if err != nil {
		return fmt.Errorf("failed to decode encryption key: %w", err)
	}

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return nil
}

func StoreAccessToken(c echo.Context, accessToken string) error {
	session, err := store.Get(c.Request(), userSessionKey)
	if err != nil {
		return fmt.Errorf("failed to get session when storing accessToken: %w", err)
	}

	session.Values[dbAccessTokenKey] = accessToken
	return session.Save(c.Request(), c.Response())
}

func GetAccessToken(c echo.Context) (string, error) {
	session, err := store.Get(c.Request(), userSessionKey)
	if err != nil {
		return "", err
	}

	accessToken, ok := session.Values[dbAccessTokenKey].(string)
	if !ok {
		return "", nil
	}

	return accessToken, nil
}

func ClearAccessToken(c echo.Context) error {
	session, err := store.Get(c.Request(), userSessionKey)
	if err != nil {
		return fmt.Errorf("failed to get session when clearing accessToken: %w", err)
	}

	delete(session.Values, dbAccessTokenKey)
	return session.Save(c.Request(), c.Response())
}
