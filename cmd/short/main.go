package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/k3mpton/shortner-project/internal/app"
	"github.com/k3mpton/shortner-project/internal/config"
	"github.com/k3mpton/shortner-project/pkg/logger"
)

func main() {
	cfg := config.MustReadCfg()
	log := logger.InitLogger(cfg.Env)

	application := app.NewApp(log, cfg.GRPC.Port)

	go application.Grpc.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Stopped server", slog.Int("port", cfg.GRPC.Port))

	application.Grpc.StopApp()

	log.Info("Server stop!")
}
