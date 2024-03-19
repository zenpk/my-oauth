package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (h *Handler) Init(conf *util.Configuration, logger *util.Logger, db *dal.Database, sv *service.Service, authInfo *util.AuthorizationInfo, tk *token.Token) {
	h.conf = conf
	h.logger = logger
	h.db = db
	h.sv = sv
	h.authInfo = authInfo
	h.tk = tk
}

func (h *Handler) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.Handle("/setup/register", h.middlewares(http.MethodPost, h.register))
	mux.Handle("/setup/client-list", h.middlewares(http.MethodGet, h.clientList))
	mux.Handle("/setup/client-create", h.middlewares(http.MethodPost, h.clientCreate))
	mux.Handle("/setup/client-delete", h.middlewares(http.MethodDelete, h.clientDelete))
	mux.Handle("/setup/public-key", h.middlewares(http.MethodGet, h.publicKey))
	mux.Handle("/auth/login", h.middlewares(http.MethodPost, h.login))
	mux.Handle("/auth/authorize", h.middlewares(http.MethodPost, h.authorize))
	mux.Handle("/auth/refresh", h.middlewares(http.MethodPost, h.refresh))
	mux.Handle("/auth/verify", h.middlewares(http.MethodPost, h.verify))
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

func (h *Handler) middlewares(method string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return h.logMiddleware(h.corsMiddleware(h.methodMiddleware(method, http.HandlerFunc(handler))))
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (h *Handler) logMiddleware(next http.Handler) http.Handler {
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
		timeNow := time.Now().UTC().Format("2006-01-02 15:04:05") // UTC guaranteed
		h.logger.Printf("%s | %v | %-7s | %v | %v\n", timeNow, sw.statusCode, r.Method, r.URL.Path, ipAddress)
	})
}

func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
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

func (h *Handler) methodMiddleware(method string, next http.Handler) http.Handler {
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
