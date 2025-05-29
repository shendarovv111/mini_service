package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anastasiakormilina/quotes-service/config"
	"github.com/anastasiakormilina/quotes-service/internal/domain/quote"
	"github.com/anastasiakormilina/quotes-service/internal/infrastructure/memory"
	"github.com/anastasiakormilina/quotes-service/internal/interfaces/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewServerConfig()

	router := bootstrap(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Настраиваем обработку сигналов остановки в отдельной горутине
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	// Запуск сервера в основной горутине
	log.Printf("Server starting on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

func bootstrap(cfg *config.ServerConfig) *chi.Mux {
	repo := memory.NewRepository()
	service := quote.NewService(repo)
	handler := api.NewHandler(service)
	return api.NewRouter(handler)
}
