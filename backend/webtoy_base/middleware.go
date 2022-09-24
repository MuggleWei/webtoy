package webtoy_base

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func MiddlewareTimeElapsed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("middleware time elapsed begin, url=%v", r.URL.Path)
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		log.Debugf("middleware time elapsed end, url=%v, elapsed=%v",
			r.URL.Path, elapsedTime)
	})
}
