package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const serverTimeout = 10 * time.Second

func main() {
	s, err := NewServer()
	if err != nil {
		log.Fatalf("NewServer: %v", err)
	}
	addr := "localhost:8080"
	srv := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Listening on http://%v", addr)
		srvErr <- srv.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-srvErr:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
