package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	shortener "github.com/k3mpton/shortner-project/pkg/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcServerAddr := os.Getenv("GRPC_SERVER_ADDR")
	fmt.Println(grpcServerAddr)
	if grpcServerAddr == "" { // если не находим в environment (docker) то используем localhost:8080
		grpcServerAddr = "localhost:8080" // По умолчанию для локальной разработки
	}

	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := shortener.RegisterLinkShortenerHandlerFromEndpoint(
		context.Background(),
		gwmux,
		grpcServerAddr,
		opts,
	)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/v1/", gwmux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		redirectHandler(w, r, grpcServerAddr)
	})

	log.Println("Gateway server started on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request, grpcServerAddr string) {
	if r.URL.Path == "/" {
		w.Write([]byte("Welcome to URL shortener! Use /v1/links to create short URLs."))
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	http.Redirect(w, r, resp.Original, http.StatusFound)

	fmt.Println("pong")
}
