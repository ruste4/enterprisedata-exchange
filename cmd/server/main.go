package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"enterprisedata-exchange/internal/config"
	"enterprisedata-exchange/internal/handler"
	logMiddleware "enterprisedata-exchange/internal/rest/middleware/logger"
	"enterprisedata-exchange/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	log := logger.SetupLogger(cfg.Env)

	exchangeUseCase, eucCleanUp, err := InitUseCase()
	if err != nil {
		log.Error("Failed to initialize use case", "error", err)
		os.Exit(1)
	}
	defer eucCleanUp()

	log.Info("Starting server", "env", cfg.Env)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logMiddleware.NewLogMiddleware(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/test/hs/exchange_dsl_1_0_0_1", func(r chi.Router) {
		r.Get("/version", handler.VersionHandler(log))
		r.Route("/v1", func(r chi.Router) {
			r.Get("/GetIBParameters", handler.GetIbParams(ctx, log))
			r.Post("/CreateExchangeNode", handler.CreateExchangeNode(ctx, log, exchangeUseCase))
			r.Post("/PutFilePart", handler.PutFile(ctx, log, exchangeUseCase))
		})
	})

	srv := &http.Server{
		Addr:         cfg.HTTPService.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPService.Timeout,
		WriteTimeout: cfg.HTTPService.Timeout,
		IdleTimeout:  cfg.HTTPService.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	log.Info("Server started", "address", cfg.HTTPService.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
