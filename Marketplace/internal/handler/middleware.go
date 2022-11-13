package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/VladimirBlinov/AuthService/pkg/authservice"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Add("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessionStore.Get(r, SessionName)
		id, ok := session.Values["user_id"]
		if !ok {
			id = 0
		}

		logger := h.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
			"method":      r.Method,
			"user_id":     id,
			"url":         r.URL.Path,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof("completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (h *Handler) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := h.service.AuthService.Authenticate(id.(int))
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		cookieSessionID, err := r.Cookie(SessionIDKey)
		if err == http.ErrNoCookie {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		} else if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		sessionS, err := h.sessionManager.Check(context.Background(), &authservice.SessionID{
			ID: cookieSessionID.Value,
		})

		if err != nil {
			h.error(w, r, http.StatusUnauthorized, err)
			return
		}

		if sessionS == nil {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, u)))
	})
}
