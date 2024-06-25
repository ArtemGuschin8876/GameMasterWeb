package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	_ "github.com/lib/pq"
)

func main() {
	var cfg application.Config
	var DB application.DB

	flag.IntVar(&cfg.Port, "port", 4000, "API Server PORT")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(DB)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Fatal(err)
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

	logger.Printf("starting server on %s", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(DB application.DB) (*sql.DB, error) {

	application.LoadEnv()
	DB.DSN = os.Getenv("DSN_DB")

	db, err := sql.Open("postgres", DB.DSN)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"./static/ui/html/table.html",
		"./static/ui/html/tableAllUsers.html",
		"./static/ui/html/addUser.html",
		"./static/ui/html/404.html",
		"./static/ui/html/updateUserForms.html",
	}

	for _, page := range pages {
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[filepath.Base(page)] = ts
	}
	return cache, nil
}
