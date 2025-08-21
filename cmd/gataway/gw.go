package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	shortener "github.com/k3mpton/shortner-project/pkg/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPC-Gateway мультиплексор
	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем gRPC-сервис
	err := shortener.RegisterLinkShortenerHandlerFromEndpoint(
		context.Background(),
		gwmux,
		"localhost:8080", // Адрес gRPC-сервера
		opts,
	)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Основной мультиплексор HTTP
	mux := http.NewServeMux()
	mux.Handle("/v1/", gwmux) // Все /v1/... идут в gRPC-Gateway

	// Обработчик для коротких ссылок (например /abc123)
	mux.HandleFunc("/", redirectHandler)

	log.Println("Gateway server started on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// Обработка главной страницы (можно вернуть HTML-форму)
		w.Write([]byte("Welcome to URL shortener! Use /v1/links to create short URLs."))
		return
	}

	// Извлекаем короткий код из URL (например "/abc123" → "abc123")
	shortCode := strings.TrimPrefix(r.URL.Path, "/") // Убираем ведущий слэш

	// Здесь нужно вызвать gRPC-метод GetOriginalLink
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Println(shortCode)

	client := shortener.NewLinkShortenerClient(conn)
	resp, err := client.GetOriginalLink(context.Background(), &shortener.GetOriginalLinkRequest{
		Short: shortCode,
	})
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Делаем редирект на оригинальный URL
	http.Redirect(w, r, resp.Original, http.StatusFound)

	fmt.Println("ALOOOOOOOOOOOOOOOOOO")
}
