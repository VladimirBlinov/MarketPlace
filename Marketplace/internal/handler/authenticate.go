package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirBlinov/AuthService/pkg/authservice"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/service"
)

func (h *Handler) handleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &service.InputUser{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.service.AuthService.SignIn(req)
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		sessionS, err := h.sessionManager.Create(context.Background(), &authservice.Session{
			UserID: int32(u.ID),
		})
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    SessionIDKey,
			Value:   sessionS.ID,
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)

		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := h.sessionStore.Save(r, w, session); err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &service.InputUser{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.service.AuthService.Register(req)
		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusCreated, u)
	}
}

func (h *Handler) handleSignOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, SessionName)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		delete(session.Values, "user_id")
		_ = session.Save(r, w)

		cookieSessionID, err := r.Cookie(SessionIDKey)
		if err == http.ErrNoCookie {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		} else if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		deleted, err := h.sessionManager.Delete(context.Background(), &authservice.SessionID{
			ID: cookieSessionID.Value,
		})
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if !deleted.Dummy {
			h.error(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Session %s not deleted", cookieSessionID.Value)))
			return
		}

		cookieSessionID.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cookieSessionID)

		h.respond(w, r.WithContext(context.Background()), http.StatusOK, nil)
	}
}

func (h *Handler) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, r.Context().Value(CtxKeyUser).(*model.User))
	}
}
