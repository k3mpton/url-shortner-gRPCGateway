package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
)

type Shortner struct {
	log          *slog.Logger
	linkSaver    LinkSaver
	linkProvider LinkProvider
}

type LinkSaver interface {
	LinkSave(
		ctx context.Context,
		originalLink string,
		shortLik string,
	) error
}

type LinkProvider interface {
	GetLink(
		ctx context.Context,
		shortLink string,
	) (string, error)
}

func NewService(
	log *slog.Logger,
	linksave LinkSaver,
	linkProvider LinkProvider,
) *Shortner {
	return &Shortner{
		log:          log,
		linkSaver:    linksave,
		linkProvider: linkProvider,
	}
}

func (s *Shortner) CreateShortLink(
	ctx context.Context,
	originalLink string,
) (string, error) {
	const op = "service.CreateShortLink"

	log := s.log.With(
		"op", op,
	)

	log.Info("short link generate...")

	_, err := url.ParseRequestURI(originalLink)
	if err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}

	hash := sha256.Sum256([]byte(originalLink))
	urlshort := base64.URLEncoding.EncodeToString(hash[:5])

	err = s.linkSaver.LinkSave(ctx, originalLink, urlshort)
	if err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}

	return "http://localhost:8080/" + urlshort, nil
}

func (s *Shortner) GetOrigLink(
	ctx context.Context,
	shortLink string,
) (string, error) {
	const op = "service.GetOriginalLink"

	log := s.log.With(
		"op", op,
	)

	log.Info("Getting original link...")

	origLink, err := s.linkProvider.GetLink(ctx, shortLink)
	if err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}

	return origLink, nil
}
