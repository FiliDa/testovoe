package transporthttp

import (
	"net/http"
	"strings"
)

// SecurityHeadersMiddleware добавляет security headers ко всем ответам
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Content Security Policy
		csp := strings.Join([]string{
			"default-src 'self'",
			"script-src 'self'",
			"style-src 'self'",
			"img-src 'self'",
			"font-src 'self'",
			"connect-src 'self'",
			"frame-ancestors 'none'",
		}, "; ")
		w.Header().Set("Content-Security-Policy", csp)

		next.ServeHTTP(w, r)
	})
}

// InputValidationMiddleware валидирует и санитизирует входные данные
func InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Валидация content type для POST/PUT запросов
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}
		}

		// Предотвращение потенциальных NoSQL инъекций (хотя используется PostgreSQL)
		// Это общая мера безопасности
		for key, values := range r.URL.Query() {
			for _, value := range values {
				if strings.Contains(strings.ToLower(value), "javascript:") ||
				   strings.Contains(strings.ToLower(value), "onerror=") ||
				   strings.Contains(strings.ToLower(value), "onload=") {
					http.Error(w, "Invalid input detected", http.StatusBadRequest)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}