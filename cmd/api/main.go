package main

import (
	"os"
	"passenger_service_backend/internal/config"
	"passenger_service_backend/internal/db"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return logger
}
func main(){

	logger := initLogger()
cfg:=	config.Load()


	if err := db.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("✅ Database connected successfully")

	logger.Info("🌱 Running database seeders...")
	if err := db.SeedDatabase(); err != nil {
		logger.Fatal("Failed to seed database", zap.Error(err))
	}
	logger.Info("✅ Database seeded successfully")


}
