package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/service"
	"github.com/zenpk/my-oauth/token"
	"github.com/zenpk/my-oauth/util"
)

type Handler struct {
	server *http.Server
	db     *dal.Database
	conf   *util.Configuration
	logger util.ILogger
	// TODO sv       service.IService
	sv       *service.Service
	authInfo *util.AuthorizationInfo
	// TODO tk     token.IToken
	tk *token.Token
}

func (h *Handler) Init(conf *util.Configuration, logger *util.Logger, sv *service.Service, authInfo *util.AuthorizationInfo, tk *token.Token) {
	h.conf = conf
	h.logger = logger
	h.sv = sv
	h.authInfo = authInfo
	h.tk = tk
}

func (h *Handler) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.Handle("/setup/register", middlewares(http.MethodPost, h.register))
	mux.Handle("/setup/client-list", middlewares(http.MethodGet, h.clientList))
	mux.Handle("/setup/client-create", middlewares(http.MethodPost, h.clientCreate))
	mux.Handle("/setup/client-delete", middlewares(http.MethodDelete, h.clientDelete))
	mux.Handle("/setup/public-key", middlewares(http.MethodGet, h.publicKey))
	mux.Handle("/auth/login", middlewares(http.MethodPost, h.login))
	mux.Handle("/auth/authorize", middlewares(http.MethodPost, h.authorize))
	mux.Handle("/auth/refresh", middlewares(http.MethodPost, h.refresh))
	mux.Handle("/auth/verify", middlewares(http.MethodPost, h.verify))
	h.server = &http.Server{
		Addr:    h.conf.HttpAddress,
		Handler: mux,
	}
	h.logger.Printf("start listening at %v\n", h.server.Addr)
	return h.server.ListenAndServe()
}

func (h *Handler) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if r.Method == "OPTIONS" {
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
			responseJson(w, data, http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
