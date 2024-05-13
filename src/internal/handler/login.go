package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	authView "github.com/VincentSchmid/htmx-templ-app-template/internal/view/auth"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/session"
	"github.com/labstack/echo/v4"
)

const (
	homePath         = "/"
	authCallbackPath = "/v1/auth/login/callback"
)

type LoginHandler struct {
	authnProvider authn.AuthenticationProvider
}

var _ IHandler = (*LoginHandler)(nil)

func NewLoginHandler(authnProvider authn.AuthenticationProvider) *LoginHandler {
	return &LoginHandler{
		authnProvider: authnProvider,
	}
}

func (lh *LoginHandler) loginIndex(c echo.Context) error {
	return render(c, authView.LoginIndex())
}

func (lh *LoginHandler) signIn(c echo.Context) error {
	redirect, err := GetRedirectUrl(c)
	if err != nil {
		redirect = homePath
	}

	credentials := authn.UserCredentials{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	accessToken, err := lh.authnProvider.SignInWithCredentials(c, credentials)
	if err != nil {
		slog.Error("failed to sign in", "err", err)
		return render(c, authView.LoginForm(credentials.Email, "invalid credentials"))
	}

	if err := session.StoreAccessToken(c, accessToken); err != nil {
		return fmt.Errorf("failed to store access token while trying to create account: %w", err)
	}

	return hxRedirect(c, redirect)
}

func (lh *LoginHandler) signInProvider(c echo.Context) error {
	provider := c.Param("provider")
	redirectUrl := appconfig.Config.Webserver.Url + authCallbackPath

	oAuthUrl, err := lh.authnProvider.SignInWithProvider(provider, redirectUrl)
	if err != nil {
		return fmt.Errorf("failed to sign in with provider: %w", err)
	}

	return c.Redirect(http.StatusSeeOther, string(oAuthUrl))
}

func (lh *LoginHandler) loginCallback(c echo.Context) error {
	accessToken := c.Request().URL.Query().Get("access_token")
	if accessToken == "" {
		return render(c, authView.CallbackScript())
	}

	err := session.StoreAccessToken(c, accessToken)
	if err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}

	return c.Redirect(http.StatusSeeOther, homePath)
}

// RegisterRoutes registers the login-related routes to an Echo instance.
func (lh *LoginHandler) RegisterRoutes(basePath string, e *echo.Group) {
	group := e.Group(basePath)
	group.GET("", lh.loginIndex)
	group.POST("", lh.signIn)
	group.GET("/provider/:provider", lh.signInProvider)
	group.GET("/callback", lh.loginCallback)
}
