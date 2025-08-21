package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saurabhraut1212/students_api_golang/internal/config"
	"github.com/saurabhraut1212/students_api_golang/internal/http/handlers/student"
)

func main() {
	//load config

	cfg := config.MustLoad()
	//database setup

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started")

	done := make(chan os.Signal, 1) //this channel is for listen os signals

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //wacthes ctrl+c ,kill command

	go func() { //run server in a goroutine
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")

		}
	}()
	<-done //Main goroutine waits until a shutdown signal (Ctrl+C) is received.

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //Creates a context with timeout → gives the server 5 seconds to shutdown gracefully before force exit.
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")
}
