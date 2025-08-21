package server

import (
	"context"
	"fmt"
	"log/slog"

	shortener "github.com/k3mpton/shortner-project/pkg/protoc"
	"google.golang.org/grpc"
)

type Shortner interface {
	CreateShortLink(
		ctx context.Context,
		originalLink string,
	) (string, error)

	GetOrigLink(
		ctx context.Context,
		shortLink string,
	) (string, error)
}

type ServerApi struct {
	shortener.UnimplementedLinkShortenerServer
	shortner_serv Shortner
}

func NewRegServer(grpcServer *grpc.Server, shortner Shortner) {
	shortener.RegisterLinkShortenerServer(
		grpcServer,
		&ServerApi{
			shortner_serv: shortner,
		},
	)
}

func (s *ServerApi) CreateShortLink(
	ctx context.Context,
	req *shortener.CreateShortLinkRequest,
) (*shortener.CreateShortLinkResponse, error) {
	const op = "server.CreateShortLink"

	log := slog.With(
		"op", op,
	)

	log.Info("Start Create Short Link...")
	fmt.Println(req.Original)

	shortLink, err := s.shortner_serv.CreateShortLink(ctx, req.GetOriginal())
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("a short link has been created")

	return &shortener.CreateShortLinkResponse{
		Short: shortLink,
	}, nil
}

func (s *ServerApi) GetOriginalLink(
	ctx context.Context,
	req *shortener.GetOriginalLinkRequest,
) (*shortener.GetOriginalLinkResponse, error) {
	const op = "server.GetOriginalLink"

	log := slog.With(
		"op", op,
	)

	log.Info("Compose shortLink to original...")
	originalLink, err := s.shortner_serv.GetOrigLink(ctx, req.GetShort())
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &shortener.GetOriginalLinkResponse{
		Original: originalLink,
	}, nil
}
