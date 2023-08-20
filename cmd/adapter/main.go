package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EgorHenek/go-oathkeeper-spicedb-adapter/configs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	ReadHeaderTimeout = 5 * time.Second
	ShutdownTimeOut   = 30 * time.Second
)

func main() {
	config := configs.NewConfig()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           run(),
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, ShutdownTimeOut)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out... forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func run() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			log.Fatal(err)
		}
	})

	return r
}
