package handlers

import (
	"log"
	"net/http"
)

func StartListening() {
	mux := http.NewServeMux()
	mux.Handle("/")

}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.Header.Get("X-Real-Ip")
		if ipAddress == "" {
			ipAddress = r.Header.Get("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		log.Printf("%v %v\n", r.URL.Path, ipAddress)
		next.ServeHTTP(w, r)
	})

}
