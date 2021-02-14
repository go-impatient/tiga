package limiter

import (
	"net/http"

	"moocss.com/tiga/pkg/rate"
	"moocss.com/tiga/pkg/server/middleware"
)

func Limiter(l *rate.Limiter) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.Allow() {
				http.Error(w,
					http.StatusText(http.StatusForbidden),
					http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
			return
		})
	}
}
