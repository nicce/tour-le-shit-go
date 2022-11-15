package logger

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("[%s]\t%s\t%s", r.Method, r.URL, time.Since(start))
	})
}
