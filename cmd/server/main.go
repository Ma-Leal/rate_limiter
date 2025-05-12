package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ma-Leal/rate-limiter/config"
	"github.com/Ma-Leal/rate-limiter/internal/middleware"
	"github.com/Ma-Leal/rate-limiter/internal/repository"
	"github.com/Ma-Leal/rate-limiter/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	storage := repository.NewRedisStorage(cfg.RedisAddr)
	useCase := usecase.NewRateLimiterUseCase(storage)

	ipCfg := usecase.RateLimiterConfig{
		Limit:         cfg.RateLimitIP,
		Window:        cfg.Window,
		BlockDuration: cfg.BlockDuration,
	}

	tokenCfg := usecase.RateLimiterConfig{
		Limit:         cfg.RateLimitToken,
		Window:        cfg.Window,
		BlockDuration: cfg.BlockDuration,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	handler := middleware.RateLimiterMiddleware(useCase, ipCfg, tokenCfg)(mux)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
