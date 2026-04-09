package transporthttp

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	ips map[string]time.Time
	mu   sync.RWMutex
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		ips: make(map[string]time.Time),
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Очистка старых записей (старше 1 минуты)
	for ip, lastSeen := range rl.ips {
		if time.Since(lastSeen) > time.Minute {
			delete(rl.ips, ip)
		}
	}

	// Проверка, сделал ли IP более 100 запросов за последнюю минуту
	count := 0
	for _, lastSeen := range rl.ips {
		if time.Since(lastSeen) <= time.Minute {
			count++
		}
	}

	if count >= 100 {
		return false
	}

	// Запись текущего запроса
	rl.ips[ip] = time.Now()
	return true
}

// RateLimitMiddleware ограничивает запросы для предотвращения злоупотреблений
func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := newRateLimiter()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получение IP клиента
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		if !limiter.allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}