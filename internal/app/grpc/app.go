package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/Gergenus/AuthService/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewApp(log *slog.Logger, port string, auth authgrpc.Auth) *App {
	Server := grpc.NewServer()
	authgrpc.Register(Server, auth)
	return &App{
		log:        log,
		gRPCServer: Server,
		port:       ":" + port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	a.log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpsapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
