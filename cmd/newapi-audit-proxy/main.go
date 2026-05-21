package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"newapi-audit-proxy/internal/audit"
	"newapi-audit-proxy/internal/config"
	"newapi-audit-proxy/internal/proxy"
	"newapi-audit-proxy/internal/web"
	"newapi-audit-proxy/migrations"
)

func main() {
	configPath := flag.String("config", envOrDefault("CONFIG_PATH", "config.yaml"), "path to config file")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("load config failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Fatalf("connect postgres failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("ping postgres failed: %v", err)
	}
	if err := migrations.Apply(ctx, pool); err != nil {
		logger.Fatalf("apply migrations failed: %v", err)
	}

	store := audit.NewStore(pool, logger)

	proxyHandler, err := proxy.New(cfg, store, logger)
	if err != nil {
		logger.Fatalf("create proxy handler failed: %v", err)
	}

	webServer, err := web.New(cfg, store, logger)
	if err != nil {
		logger.Fatalf("create web server failed: %v", err)
	}

	rootMux := http.NewServeMux()
	rootMux.Handle(cfg.WebBasePath+"/", http.StripPrefix(cfg.WebBasePath, webServer.Handler()))
	rootMux.Handle("/", proxyHandler)

	httpServer := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           rootMux,
		ReadHeaderTimeout: 15 * time.Second,
	}

	errCh := make(chan error, 1)

	go serve(logger, "audit", httpServer, errCh)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalCh)

	select {
	case sig := <-signalCh:
		logger.Printf("received signal %s, shutting down", sig)
	case err := <-errCh:
		if err != nil {
			logger.Printf("server stopped with error: %v", err)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Printf("graceful shutdown timed out, forcing close")
			if closeErr := httpServer.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
				logger.Printf("force close audit server failed: %v", closeErr)
			}
		} else {
			logger.Printf("shutdown audit server failed: %v", err)
		}
	}
}

func serve(logger *log.Logger, name string, server *http.Server, errCh chan<- error) {
	logger.Printf("%s server listening on %s", name, server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- err
		return
	}
	errCh <- nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
