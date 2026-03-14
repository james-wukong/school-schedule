package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/james-wukong/school-schedule/internal/application/usecase"
	"github.com/james-wukong/school-schedule/internal/config"
	"github.com/james-wukong/school-schedule/internal/infrastructure/cache/redis"
	"github.com/james-wukong/school-schedule/internal/infrastructure/messaging/kafka"
	"github.com/james-wukong/school-schedule/internal/infrastructure/persistence/postgres"
	httpInterface "github.com/james-wukong/school-schedule/internal/interface/http"
)

func main() {
	// ==========================================
	// Root Context (for graceful shutdown)
	// ==========================================
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleShutdown(cancel)

	// ==========================================
	// Load Config
	// ==========================================
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ==========================================
	// Initialize Infrastructure
	// ==========================================

	// PostgreSQL
	pg, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer pg.Close()

	// Redis
	redisClient, err := redis.New(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer redisClient.Close()

	// Kafka Producer
	producer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatalf("kafka producer failed: %v", err)
	}
	defer producer.Close()

	// ==========================================
	// Build Use Cases
	// ==========================================
	scheduleUseCase := usecase.NewScheduleUseCase(
		pg,
		redisClient,
		producer,
	)

	// ==========================================
	// Setup Gin
	// ==========================================
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(httpInterface.TenantMiddleware())

	// Routes
	httpInterface.RegisterRoutes(router, scheduleUseCase)

	// ==========================================
	// HTTP Server Config
	// ==========================================
	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ==========================================
	// Start Server
	// ==========================================
	go func() {
		log.Printf("Scheduler service running on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	log.Println("Shutting down scheduler service...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("Scheduler service stopped gracefully")
}

// ==========================================
// Graceful Shutdown Handler
// ==========================================

func handleShutdown(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
}
