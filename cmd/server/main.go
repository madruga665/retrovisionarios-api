package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"retrovisionarios-api/config/env"

	v1 "retrovisionarios-api/internal/app/v1"
	"retrovisionarios-api/internal/app/v1/events/controllers"
	"retrovisionarios-api/internal/app/v1/events/repositories"
	"retrovisionarios-api/internal/app/v1/events/services"
	postgres "retrovisionarios-api/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := env.Load()

	// ... (logger configuration)
	var handler slog.Handler
	if os.Getenv("GIN_MODE") == "release" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(handler))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configuração do CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = strings.Split(cfg.AllowedOrigins, ",")
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Configuração de Trusted Proxies via Config
	router.SetTrustedProxies(strings.Split(cfg.TrustedProxies, ","))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool, err := postgres.DbPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Erro fatal ao conectar ao banco", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	eventRepository := repositories.NewEventRepository(dbPool)
	eventService := services.NewEventService(eventRepository)
	eventController := controllers.NewEventController(eventService)

	v1.EventRoutes(router, eventController)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Executa o servidor em uma goroutine
	go func() {
		slog.Info("Iniciando servidor", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Erro no servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	// Aguarda sinal de interrupção
	<-ctx.Done()
	slog.Info("Desligando servidor graciosamente...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Forçando desligamento do servidor", "error", err)
		os.Exit(1)
	}

	slog.Info("Servidor finalizado com sucesso.")
}
