package handlers

import (
	"encoding/json"
	"github.com/zenpk/my-oauth/utils"
	"log"
	"net/http"
)

func StartListening() error {
	mux := http.NewServeMux()
	mux.Handle("/user/register", middlewares(register))
	//mux.Handle("/user/login", middlewares(login))
	//mux.Handle("/client/create")
	//mux.Handle("/client/")
	log.Printf("start listening at %v\n", utils.HttpAddress)
	return http.ListenAndServe(utils.HttpAddress, mux)
}

func middlewares(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return logMiddleware(corsMiddleware(http.HandlerFunc(handler)))
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
		log.Printf("%v|%-7s|%v|%v\n", sw.statusCode, r.Method, r.URL.Path, ipAddress)
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

func responseMethodUnsupported(w http.ResponseWriter) {
	data := commonResp{
		Ok:  false,
		Msg: "HTTP method not supported",
	}
	responseJson(w, data, http.StatusBadRequest)
}

func responseNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
