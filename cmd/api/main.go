package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"gamemasterweb.net/internal/logger"
	_ "github.com/lib/pq"
)

func main() {
	var cfg application.Config
	var DB application.DB

	zeroLog := logger.NewLogger()

	flag.IntVar(&cfg.Port, "port", 4000, "API Server PORT")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(DB)
	if err != nil {
		zeroLog.Fatal().Msg("Database connection is not established ")
	}
	defer db.Close()

	zeroLog.Info().Msg("Database connection pool established")

	templateCache, err := application.ReadTemplates()
	if err != nil {
		zeroLog.Err(err).Msg("Template reading error")
	}

	app := application.Application{
		Config:    cfg,
		Logger:    logger,
		Storage:   data.NewStorage(db),
		Templates: templateCache,
		Response:  application.Response{},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes(&app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	zeroLog.Info().Msgf("Starting server on %s", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		zeroLog.Fatal().Msg("Error in server startup")
	}

}

func openDB(DB application.DB) (*sql.DB, error) {

	application.LoadEnv()
	DB.DSN = os.Getenv("DSN_DB")

	zeroLog := logger.NewLogger()

	db, err := sql.Open("postgres", DB.DSN)
	if err != nil {
		zeroLog.Err(err).Msg("Unknown driver db")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		zeroLog.Err(err).Msg("Problem with PingContext")
		return nil, err
	}

	return db, nil
}
