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

	"github.com/rabin-nyaundi/authentication-service/internal/data"

	_ "github.com/lib/pq"
)

type JSONResponse struct {
	Error   bool        `json:"error,omitempty"`
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// type envelope map[string]interface{

// }

type Config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config Config
	models data.Models
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 80, "Authentication server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_DSN"), "Database connection string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL maximum open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL maximum idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "10m", "PostgreSQL maximum idle time")
	flag.Parse()

	db, err := OpenDB(cfg)
	if err != nil {
		log.Panic(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: data.NewModel(db),
	}

	log.Printf("Starting server at port:%d", cfg.port)
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err = svr.ListenAndServe()

	if err != nil {
		log.Panic(err)
		return
	}

	log.Printf("server started at port:%d", cfg.port)
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		log.Panic(err)
		return nil, err
	}
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return db, nil
}
