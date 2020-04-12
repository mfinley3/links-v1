package middleware

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/metrics/influx"
)

func RedirectCounter(redirectCounter *influx.Counter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if w.Header().Get("location") != "" {
				routeContext := chi.RouteContext(r.Context())
				redirectCounter.With("short_url", routeContext.URLParam("short_url")).Add(1)
			}
		})
	}
}
