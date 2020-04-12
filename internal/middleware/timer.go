package middleware

import (
	"log"
	"net/http"
	"time"
)

func Timer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("Redirect took: %s", time.Since(start))
		})
	}
}
