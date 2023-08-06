package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/zenpk/my-oauth/utils"
	"log"
	"net/http"
)

func StartListening() error {
	mux := http.NewServeMux()
	mux.Handle("/user/register", middlewares(http.MethodPost, register))
	mux.Handle("/user/login", middlewares(http.MethodPost, login))
	//mux.Handle("/client/create")
	//mux.Handle("/client/")
	log.Printf("start listening at %v\n", utils.Conf.HttpAddress)
	return http.ListenAndServe(utils.Conf.HttpAddress, mux)
}

func middlewares(method string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return logMiddleware(corsMiddleware(methodMiddleware(method, http.HandlerFunc(handler))))
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			w,
			http.StatusOK,
		}

		next.ServeHTTP(sw, r)

		ipAddress := r.Header.Get("X-Real-Ip")
		if ipAddress == "" {
			ipAddress = r.Header.Get("X-Forwarded-For")
		}
		if ipAddress == "" {
			ipAddress = r.RemoteAddr
		}
		log.Printf("| %v | %-7s | %v | %v\n", sw.statusCode, r.Method, r.URL.Path, ipAddress)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Request-Method", "GET, POST, PUT, PATCH, DELETE")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func methodMiddleware(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			data := commonResp{
				Ok:  false,
				Msg: fmt.Sprintf("HTTP method %v is not supported", r.Method),
			}
			responseJson(w, data, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func responseJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func responseError(w http.ResponseWriter, err error, statusCode int) {
	data := commonResp{
		Ok:  false,
		Msg: err.Error(),
	}
	responseJson(w, data, statusCode)
}

func responseNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
