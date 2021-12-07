package main

import (
	"context"
	"forum/dbs"
	"forum/internal"
	"forum/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	addr := ":8080"

	templ, err := internal.GetTempl()
	if err != nil {
		log.Println(err)
		return
	}

	if err := dbs.NewConnect(); err != nil {
		log.Println(err)
		return
	}

	limiter := new(models.Limiter)
	limiter.IPs = make(map[string]*models.Counter)

	mux := internal.Register(templ, limiter)

	go internal.CleanupVisitors(limiter)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("main: running simple server on https://localhost%s", addr)

	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Printf("HTTP server ListenAndServe: %v", err)
		return
	}
	<-idleConnsClosed
}
