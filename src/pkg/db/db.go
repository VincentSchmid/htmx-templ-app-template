package database

import (
	"database/sql"
	"runtime"
	"time"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

type DbProvider struct {
	config     appconfig.Database
	Db         *sql.DB
	Bun        *bun.DB
	SqlAdapter *sqladapter.Adapter
}

func NewDbProvider(config appconfig.Database, debug bool) *DbProvider {
	db := initializeDb(config)
	bunDb := initializeBunDb(db, debug)
	return &DbProvider{
		config: config,
		Db:     db,
		Bun:    bunDb,
	}
}

func initializeDb(config appconfig.Database) *sql.DB {
	db, err := sql.Open(config.GetDriver(), config.GetConnectionString())
	if err != nil {
		logger.Log.Error("Error opening database: %v", zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		logger.Log.Error("Error pinging database: %v", zap.Error(err))
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	runtime.SetFinalizer(db, finalizer)

	return db
}

func initializeBunDb(db *sql.DB, debug bool) *bun.DB {
	bunDb := bun.NewDB(db, pgdialect.New())
	if debug {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(false)))
	}
	return bunDb
}

func finalizer(db *sql.DB) {
	if db == nil {
		return
	}

	err := db.Close()
	if err != nil {
		panic(err)
	}
}
