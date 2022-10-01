package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/VladimirBlinov/MarketPlace/internal/model"
	"github.com/VladimirBlinov/MarketPlace/internal/service"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	SessionName            = "MarketPlace"
	CtxKeyUser      ctxKey = iota
	ctxKeyRequestID ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type Handler struct {
	service      *service.Service
	sessionStore sessions.Store
	store        store.Store
	Router       *mux.Router
	logger       *logrus.Logger
}

func NewHandler(store store.Store, service *service.Service, sessionStore sessions.Store) *Handler {
	return &Handler{
		service:      service,
		store:        store,
		sessionStore: sessionStore,
		Router:       mux.NewRouter(),
		logger:       logrus.New(),
	}
}

func (h *Handler) InitHandler() {
	h.Router.Use(h.setRequestID)
	h.Router.Use(h.logRequest)
	h.Router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.ExposedHeaders([]string{"Set-Cookie"}),
		handlers.AllowCredentials()))
	h.Router.HandleFunc("/users", h.handleUsersCreate()).Methods("POST")
	h.Router.HandleFunc("/sessions", h.handleSessionsCreate()).Methods("POST")

	private := h.Router.PathPrefix("/private").Subrouter()
	private.Use(h.AuthenticateUser)
	private.HandleFunc("/whoami", h.handleWhoami()).Methods("GET")
	private.HandleFunc("/signout", h.handleSignOut()).Methods("GET")

	product := private.PathPrefix("/product").Subrouter()
	product.HandleFunc("/create", h.handleProductCreate()).Methods("POST")
	product.HandleFunc("/product/{id}", h.handleProductGetProductById()).Methods("GET")
	product.HandleFunc("/all", h.handleProductList()).Methods("GET")
	product.HandleFunc("/category/get_categories", h.handleProductCategoryGetAll()).Methods("GET")
	product.HandleFunc("/material/get_materials", h.handleProductGetMaterials()).Methods("GET")
}

func (h *Handler) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := h.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
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

		u, err := h.store.User().FindById(id.(int))
		if err != nil {
			h.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, u)))
	})
}

func (h *Handler) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.respond(w, r, http.StatusOK, r.Context().Value(CtxKeyUser).(*model.User))
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

func (h *Handler) handleUsersCreate() http.HandlerFunc {
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

func (h *Handler) handleSessionsCreate() http.HandlerFunc {
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

func (h *Handler) handleProductCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//req := &service.InputProduct{}
		req := &model.Product{}
		req.UserID = r.Context().Value(CtxKeyUser).(*model.User).ID
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.service.ProductService.CreateProduct(req); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusCreated, nil)
	}
}

func (h *Handler) handleProductGetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqVars := mux.Vars(r)
		productId, err := strconv.Atoi(reqVars["id"])
		if err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		product, err := h.service.ProductService.GetProductById(productId)
		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, product)
	}
}

func (h *Handler) handleProductList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(CtxKeyUser).(*model.User)

		products, err := h.service.ProductService.GetProductsByUserId(u.ID)
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, products)
	}
}

func (h *Handler) handleProductCategoryGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.service.ProductService.GetProductCategories()
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, categories)
	}
}

func (h *Handler) handleProductGetMaterials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		materials, err := h.service.ProductService.GetProductMaterials()
		if err != nil {
			h.error(w, r, http.StatusInternalServerError, err)
			return
		}

		h.respond(w, r, http.StatusOK, materials)
	}
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *Handler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
