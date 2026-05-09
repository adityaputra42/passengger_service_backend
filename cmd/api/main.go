package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"passenger_service_backend/cmd/routes"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/db"
	"passenger_service_backend/internal/injection"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		NameKey:      "logger",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func main() {
	logger := initLogger()
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────
	if err := db.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("✅ Database connected")

	if err := db.SeedDatabase(); err != nil {
		logger.Fatal("Failed to seed database", zap.Error(err))
	}
	logger.Info("✅ Database seeded")

	handler, err := injection.InitializeAllHandler(cfg, context.Background())
	if err != nil {
		logger.Fatal("Failed to initialise dependencies", zap.Error(err))
	}
	logger.Info("✅ Dependencies wired")

	router := routes.SetupRoutes(handler, logger, cfg.CORS)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("\n╔════════════════════════════════════════╗\n")
	fmt.Printf("║  Server running on port %-4s           ║\n", port)
	fmt.Printf("╚════════════════════════════════════════╝\n\n")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Forced shutdown", zap.Error(err))
	}
	logger.Info("Server stopped cleanly")
}
