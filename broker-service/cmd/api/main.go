package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	port int
}

type application struct {
	config Config
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 80, "API Server port")
	flag.Parse()

	app := &application{
		config: cfg,
	}

	log.Printf("Stating broker server on port %d\n", app.config.port)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err := svr.ListenAndServe()

	if err != nil {
		log.Fatal("Server could not start")
	}
}
