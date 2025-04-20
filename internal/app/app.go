package app

import (
	"log/slog"

	grpcapp "github.com/Gergenus/AuthService/internal/app/grpc"
	"github.com/Gergenus/AuthService/internal/grpc/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort string, auth auth.Auth) *App {
	grpcApp := grpcapp.NewApp(log, grpcPort, auth)
	return &App{
		GRPCSrv: grpcApp,
	}
}
