package handler

import (
	"context"
	"encoding/json"
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"net/http"
)

func (h *Handler) handleSignIn() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := h.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			h.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

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
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			UserRole: 2,
			Active:   true,
		}
		if err := h.store.User().Create(u); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
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

		h.respond(w, r.WithContext(context.Background()), http.StatusOK, nil)
	}
}

func (h *Handler) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, r.Context().Value(CtxKeyUser).(*model.User))
	}
}
