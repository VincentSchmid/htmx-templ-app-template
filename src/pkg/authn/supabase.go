package authn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
)

const (
	recoverEndpoint = "/auth/v1/recover"
)

type SupabaseProvider struct {
	client      *supabase.Client
	supabaseUrl string
	supabaseKey string
}

var _ AuthenticationProvider = (*SupabaseProvider)(nil)

func NewSupabaseProvider(supabaseConfig *appconfig.Supabase) *SupabaseProvider {
	return &SupabaseProvider{
		client:      supabase.CreateClient(supabaseConfig.Url, supabaseConfig.Key),
		supabaseUrl: supabaseConfig.Url,
		supabaseKey: supabaseConfig.Key,
	}
}

func (sp *SupabaseProvider) SignUpWithCredentials(e echo.Context, credentials UserCredentials) (User, error) {
	if credentials.Password != credentials.ConfirmPassword {
		return User{}, fmt.Errorf("passwords do not match")
	}

	spCredentials := supabase.UserCredentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	}

	spUser, err := sp.client.Auth.SignUp(e.Request().Context(), spCredentials)
	if err != nil {
		return User{}, fmt.Errorf("failed to sign in: %w", err)
	}

	user := User{
		Uuid:  uuid.MustParse(spUser.ID),
		Email: spUser.Email,
	}

	return user, nil
}

func (sp *SupabaseProvider) SignInWithCredentials(e echo.Context, credentials UserCredentials) (string, error) {
	sbCredentials := supabase.UserCredentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	}

	resp, err := sp.client.Auth.SignIn(e.Request().Context(), sbCredentials)
	if err != nil {
		return "", fmt.Errorf("failed to sign in: %w", err)
	}

	return resp.AccessToken, nil
}

func (sp *SupabaseProvider) SignInWithProvider(provider string, redirectUrl string) (OAuthUrl, error) {
	resp, err := sp.client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   provider,
		RedirectTo: redirectUrl,
	})

	if err != nil {
		return "", fmt.Errorf("failed to sign in with provider: %w", err)
	}

	return OAuthUrl(resp.URL), nil
}

func (sp *SupabaseProvider) AuthenticateWithToken(e echo.Context, token string) (model.AuthenticatedUser, error) {
	resp, err := sp.client.Auth.User(e.Request().Context(), string(token))
	if err != nil {
		return model.AuthenticatedUser{}, fmt.Errorf("failed to authenticate with token: %w", err)
	}

	user := model.AuthenticatedUser{
		Uuid:        uuid.MustParse(resp.ID),
		Email:       resp.Email,
		AccessToken: token,
		Account:     &model.Account{},
		IsLoggedIn:  true,
	}

	return user, nil
}

func (sp *SupabaseProvider) ChangePassword(e echo.Context, token string, newCredentials UserCredentials) error {
	if newCredentials.Password != newCredentials.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	data := map[string]any{
		"password": newCredentials.Password,
	}

	_, err := sp.client.Auth.UpdateUser(e.Request().Context(), string(token), data)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (sp *SupabaseProvider) ResetPassword(e echo.Context, email string, redirectUrl string) error {
	recoverUrl := sp.supabaseUrl + recoverEndpoint
	params := map[string]any{
		"email":       email,
		"redirect_to": redirectUrl,
	}

	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal reset password params: %w", err)
	}

	req, err := http.NewRequest("POST", recoverUrl, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request to reset password: %w", err)
	}
	req.Header.Set("apikey", sp.supabaseKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to reset password: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to reset password: %s", body)
	}

	return nil
}
