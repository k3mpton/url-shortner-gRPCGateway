package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	urlMap  = make(map[string]string) // Хранит соответствие short → long URL
	urlLock sync.RWMutex              // Защита от гонок данных
)

func main() {
	// Пример: генерируем короткую ссылку
	originalURL := "https://yandex.ru/pogoda/ru/38/details?lat=48.672256&lon=44.437309&utm_source=serp&utm_campaign=helper&utm_medium=desktop&utm_content=helper_desktop_main&utm_term=10_days"
	shortURL, err := shortenURL(originalURL)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Оригинальная ссылка: %s\n", originalURL)
	fmt.Printf("Сокращенная ссылка: %s\n", shortURL)

	// Запускаем веб-сервер
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/alo", aloHand)

	// fmt.Println("Сервер запущен на http://192.168.0.101:8080")
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8083", nil)
}

func aloHand(w http.ResponseWriter, r *http.Request) {
	fmt.Println("aloooo")
}

// Генерация короткой ссылки (как у вас)
func shortenURL(longURL string) (string, error) {
	log.Println("shortenUrl")

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return "", fmt.Errorf("неверный URL: %v", err)
	}

	hash := sha256.Sum256([]byte(longURL))
	shortHash := base64.URLEncoding.EncodeToString(hash[:32])
	// shortHash = strings.TrimRight(shortHash, "=")

	// shortURL := "http://192.168.0.101:8080/" + shortHash // Теперь ссылки ведут на наш сервер
	shortURL := "http://localhost:8083/" + shortHash // Теперь ссылки ведут на наш сервер

	// Сохраняем в память
	urlLock.Lock()
	urlMap[shortHash] = longURL
	urlLock.Unlock()

	return shortURL, nil
}

// Обработчик перенаправления (при переходе по короткой ссылке)
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/")
	log.Println("RedirectHandler")
	log.Println(shortKey, r.URL.Path, r.URL)

	urlLock.RLock()
	longURL, exists := urlMap[shortKey]
	urlLock.RUnlock()

	if !exists {
		http.Error(w, "Ссылка не найдена", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound) // 302 редирект
}

// Обработчик создания короткой ссылки (через веб)
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	log.Println("shortenHandler", r.FormValue("url"))

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL не указан", http.StatusBadRequest)
		return
	}

	shortURL, err := shortenURL(longURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Короткая ссылка: <a href='%s'>%s</a>", "бе", shortURL)
}
