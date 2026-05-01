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
	eventCtrl "retrovisionarios-api/internal/app/v1/events/controller"
	eventRepo "retrovisionarios-api/internal/app/v1/events/repository"
	eventServ "retrovisionarios-api/internal/app/v1/events/service"
	videoCtrl "retrovisionarios-api/internal/app/v1/videos/controller"
	videoRepo "retrovisionarios-api/internal/app/v1/videos/repository"
	VideoServ "retrovisionarios-api/internal/app/v1/videos/service"
	postgres "retrovisionarios-api/internal/db"

	_ "retrovisionarios-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Retrovisionarios API
// @version         1.0
// @description     API para gerenciamento de eventos e shows da banda Retrovisionarios.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5000
// @BasePath  /v1

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

	// Events
	eventRepository := eventRepo.NewEventRepository(dbPool)
	eventService := eventServ.NewEventService(eventRepository)
	eventController := eventCtrl.NewEventController(eventService)

	// Videos
	videoRepository := videoRepo.NewVideoRepository(dbPool)
	videoService := VideoServ.NewVideoService(videoRepository)
	videoController := videoCtrl.NewVideoController(videoService)

	v1.SetupRoutes(router, v1.RouterConfig{
		EventController: eventController,
		VideoController: videoController,
	})

	// Rota do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/v1/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

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
