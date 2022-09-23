package webtoy_base

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func MiddlewareTimeElapsed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("recv url=%v, params=%v", r.URL.Path, r.URL.RawQuery)
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		log.Debugf("url=%v, params=%v, elapsed=%v",
			r.URL.Path, r.URL.RawQuery, elapsedTime)
	})
}
