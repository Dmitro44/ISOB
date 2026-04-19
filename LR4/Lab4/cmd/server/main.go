package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(2, 3)
		visitors[ip] = limiter
	}

	return limiter
}

func secureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(host)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	host, port, _ := net.SplitHostPort(r.RemoteAddr)

	_, err := io.ReadAll(r.Body)
	if err != nil {

		fmt.Printf("Connection dropped for %s. Reason: %v\n", host, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Printf("Connection established for %s:%s\n", host, port)
}

func readBodyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {

		fmt.Printf("Attack prevented. Error: %v\n", err)
		http.Error(w, "Payload too large", http.StatusRequestEntityTooLarge)
		return
	}

	fmt.Printf("Request was processed successfully\n")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/readBody", readBodyHandler)

	secureHandler := secureMiddleware(mux)

	server := &http.Server{
		Addr:         ":8089",
		Handler:      secureHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	server.ListenAndServe()
}
