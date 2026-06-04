package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mechaniel-coder/netwatch/server/internal/api"
	"github.com/mechaniel-coder/netwatch/server/internal/config"
	"github.com/mechaniel-coder/netwatch/server/internal/db"
	"github.com/mechaniel-coder/netwatch/server/internal/ws"
)

func main() {
	cfg := config.MustLoad()

	pool, err := db.Connect(cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("db: failed to connect: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatalf("db: migrations failed: %v", err)
	}

	hub := ws.NewHub()
	go hub.Run()

	router := api.NewRouter(cfg, pool, hub)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("netwatch server listening on :%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	log.Println("shutdown complete")
}
