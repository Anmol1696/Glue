package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

func (a *APIServer) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		defer func() {
			a.logger.Info("client request",
				zap.Duration("latency", time.Since(start)),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.String("client_ip", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path))
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

// Extracts email for the user from a context
func getSubject(ctx context.Context) string {
	ui, ok := ctx.Value(IdentityCtxKey).(*UserIdentity)
	if !(ok) {
		return "identity-not-found"
	}
	return ui.Email
}

// Function for extracting user identity making request
// Current function simply extracts user identity from the headers
func (a *APIServer) NewUserIdentity(r *http.Request) (*UserIdentity, *ErrResponse) {
	subject := r.Header.Get("X-Auth-Subject")
	email := r.Header.Get("X-Auth-Email")
	roles := strings.Split(r.Header.Get("X-Auth-Roles"), ",")
	if email == "" {
		return nil, ErrIdentityNotFound
	}
	ui := &UserIdentity{Subject: subject, Roles: roles, Email: email}
	return ui, nil
}

func (a *APIServer) IdentityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdn, err := a.NewUserIdentity(r)
		if err != nil {
			render.Render(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), IdentityCtxKey, userIdn)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
