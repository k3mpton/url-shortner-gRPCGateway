package grpcApp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/k3mpton/shortner-project/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log      *slog.Logger
	grpcServ *grpc.Server
	port     int
}

func NewRegApp(log *slog.Logger, port int, shortner server.Shortner) *App {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	server.NewRegServer(grpcServer, shortner)

	return &App{
		log:      log,
		grpcServ: grpcServer,
		port:     port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	op := "appgrpc.Run"

	log := a.log.With(
		"op", op,
		"port", a.port,
	)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%v", a.port))
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Start server: ",
		slog.String("addr", conn.Addr().String()))

	if err := a.grpcServ.Serve(conn); err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (a *App) StopApp() {
	const op = "grpcApp.StopApp"

	a.log.With(
		"op", op,
	).Info("stopped grpc server...", slog.Int("port", a.port))

	a.grpcServ.GracefulStop()
}
