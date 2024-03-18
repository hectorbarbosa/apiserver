package main

import (
	"context"
	"errors"
	"filmoteka/apirouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	// routing
	ar := apirouter.NewApiRouter()
	err = ar.Start()
	if err != nil {
		log.Fatal("Router start: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/actors/", ar)
	mux.Handle("/films/", ar)

	// server
	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	// close db connection
	ar.Stop()
}
