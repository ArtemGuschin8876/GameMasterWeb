package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"gamemasterweb.net/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	logger    *log.Logger
	config    config
	storage   data.Storage
	templates map[string]*template.Template
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API Server PORT")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Fatal(err)
	}

	app := application{
		config:    cfg,
		logger:    logger,
		storage:   data.NewStorage(db),
		templates: templateCache,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
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

func openDB(cfg config) (*sql.DB, error) {

	LoadEnv()
	cfg.db.dsn = os.Getenv("DSN_DB")

	db, err := sql.Open("postgres", cfg.db.dsn)
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
		"./static/ui/html/successfullCreatedUser.html",
		"./static/ui/html/404.html",
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
