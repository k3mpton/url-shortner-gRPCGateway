package app

import (
	"log/slog"

	grpcApp "github.com/k3mpton/shortner-project/internal/app/grpc"
	"github.com/k3mpton/shortner-project/internal/service"
	"github.com/k3mpton/shortner-project/internal/storage/postgres"
)

type App struct {
	Grpc *grpcApp.App
}

func NewApp(log *slog.Logger, port int) *App {
	storage := postgres.NewStorage()

	shortener := service.NewService(log, storage, storage)

	app := grpcApp.NewRegApp(log, port, shortener)

	return &App{
		Grpc: app,
	}
}
