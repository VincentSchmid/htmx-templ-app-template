package cmd

import (
	"fmt"
	"net/http"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/handler"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/repository"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/service"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authn"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/authz"
	database "github.com/VincentSchmid/htmx-templ-app-template/pkg/db"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/events"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/session"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var PublicHandler http.Handler

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the webserver",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initConfigurableDependecies(); err != nil {
			logger.Log.Error("failed to initialize: ", zap.Error(err))
		}

		dbProvider := database.NewDbProvider(appconfig.Config.Database, appconfig.Config.Debug)
		authnProvider := authn.NewSupabaseProvider(appconfig.Config.Supabase)
		authzProvider := authz.NewCasbinAuthzProvider(appconfig.Config.Authz, handler.GetUserIdFunc(), dbProvider.Db)
		eventManager := events.NewEventManager()

		repositories := repository.NewRepositories(dbProvider.Bun)
		services := service.NewServices(repositories, authzProvider, eventManager)
		handlers := handler.NewHandlers(services, authnProvider)

		app := echo.New()
		appNoUser := app.Group("")
		withUser := app.Group("")
		withUser.Use(handler.GetWithUserMiddleware(authnProvider))
		withUser.GET("/public/*", echo.WrapHandler(PublicHandler))
		withUser.GET("/", handler.Make(handler.HandleHomeIndex))

		v1 := withUser.Group("/v1")
		auth := v1.Group("/auth")
		authNoUser := appNoUser.Group("v1/auth")

		handlers.LoginHandler.RegisterRoutes("/login", authNoUser)
		handlers.SignUpHandler.RegisterRoutes("/signup", authNoUser)
		handlers.PasswordHandler.RegisterRoutes("/password", authNoUser)
		handlers.LogoutHandler.RegisterRoutes("/logout", auth)

		user := v1.Group("/user")
		user.Use(handler.WithAuthn)
		handlers.UserHandler.RegisterRoutes("", user)

		app.Logger.Fatal(app.Start(appconfig.Config.Webserver.Host))
	},
}

func initConfigurableDependecies() error {
	if err := logger.Init(appconfig.Config.Debug); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := session.Init(appconfig.Config.Session); err != nil {
		return fmt.Errorf("failed to initialize session: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
