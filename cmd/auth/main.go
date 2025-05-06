package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gergenus/AuthService/internal/app"
	"github.com/Gergenus/AuthService/internal/config"
	"github.com/Gergenus/AuthService/internal/pkg/database"
	"github.com/Gergenus/AuthService/internal/repository"
	"github.com/Gergenus/AuthService/internal/services/auth"
)

func main() {
	config.InitConfig()
	log := setupLogger()
	log.Info("Starting application", slog.String("port", os.Getenv("GRPC_PORT")))
	database := database.InitDB(os.Getenv("POSTRGES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DBNAME"))
	repository := repository.NewPostgresRepository(log, database)
	service := auth.NewAuth(log, repository, repository)
	application := app.NewApp(log, os.Getenv("GRPC_PORT"), service)
	go application.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}
