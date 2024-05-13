package test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/constants"
	"github.com/VincentSchmid/htmx-templ-app-template/internal/model"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	database "github.com/VincentSchmid/htmx-templ-app-template/pkg/db"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type RequestBodyType int

const (
	RequestBodyTypeForm RequestBodyType = iota
	RequestBodyTypeHTML
)

type RequestHelper struct {
	CurrentAccount *model.Account
	RequestBody    string
	Params         map[string]string
	Form           map[string]string
	BodyType       RequestBodyType
}

func CallHandler(helper RequestHelper, handlerFunc echo.HandlerFunc) (*httptest.ResponseRecorder, *goquery.Document, error) {
	e := echo.New()

	var request *http.Request

	if helper.BodyType == RequestBodyTypeForm {
		form := make(url.Values)
		for key, value := range helper.Form {
			form.Set(key, value)
		}

		request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(helper.RequestBody))
		request.Header.Set("Content-Type", "text/html")
	}

	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	authenticatedUser := *model.NewAuthenticatedUser()
	if helper.CurrentAccount != nil {
		authenticatedUser = model.AuthenticatedUser{
			Uuid:       uuid.New(),
			Email:      fmt.Sprintf("%s@email.com", helper.CurrentAccount.Username),
			IsLoggedIn: true,
			Account:    helper.CurrentAccount,
		}
	}

	tmpCtx := context.WithValue(ctx.Request().Context(), constants.USER_CONTEXT_KEY, authenticatedUser)
	ctx.SetRequest(ctx.Request().WithContext(tmpCtx))

	for key, value := range helper.Params {
		ctx.SetParamNames(key)
		ctx.SetParamValues(value)
	}

	handlerFunc(ctx)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(recorder.Body.String()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create goquery document: %w", err)
	}

	return recorder, doc, nil
}

func initializeMigrations(db *sql.DB, config appconfig.Database) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		config.GetMigrationDirectory(),
		config.GetDriver(),
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migration, nil
}

func MigrationReset(dbProvider *database.DbProvider, config appconfig.Database) error {
	migration, err := initializeMigrations(dbProvider.Db, config)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	logger.Log.Info("Dropping database")

	err = migration.Drop()
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	migration, err = initializeMigrations(dbProvider.Db, config)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	logger.Log.Info("Applying migrations")

	err = migration.Up()
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
