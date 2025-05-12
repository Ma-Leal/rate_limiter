package middleware

import (
	"net"
	"net/http"

	"github.com/Ma-Leal/rate-limiter/internal/usecase"
)

func RateLimiterMiddleware(useCase *usecase.RateLimiterUseCase, ipCfg, tokenCfg usecase.RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("API_KEY")
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			var key string
			var cfg usecase.RateLimiterConfig

			if token != "" {
				key = "token:" + token
				cfg = tokenCfg
			} else {
				key = "ip:" + ip
				cfg = ipCfg
			}

			allowed, err := useCase.Allow(key, cfg)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
