package main

import (
	"context"
	"fmt"
	"log"
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

	"github.com/gin-gonic/gin"
)

func main() {
	env.Load()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configuração de Trusted Proxies via Env
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		trustedProxies = "localhost"
	}
	router.SetTrustedProxies(strings.Split(trustedProxies, ","))

	dbPool, err := postgres.DbPool()

	if err != nil {
		log.Fatal("Erro ao conectar ao banco", err)
	}

	defer dbPool.Close()

	eventRepository := repositories.NewEventRepository(dbPool)
	eventService := services.NewEventService(eventRepository)
	eventController := controllers.NewEventController(eventService)

	v1.EventRoutes(router, eventController)

	// Configuração da Porta via Env
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "5000"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Executa o servidor em uma goroutine para não bloquear
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Aguarda sinal de interrupção para desligar graciosamente (timeout de 5s)
	quit := make(chan os.Signal, 1)
	// kill (sem param) envia syscall.SIGTERM
	// kill -2 envia syscall.SIGINT
	// kill -9 envia syscall.SIGKILL (não pode ser capturado)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Forçando desligamento do servidor:", err)
	}

	log.Println("Servidor finalizado.")
}
