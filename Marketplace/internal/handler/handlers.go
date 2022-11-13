package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VladimirBlinov/AuthService/pkg/authservice"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	SessionName            = "MarketPlace"
	SessionIDKey           = "session_id"
	CtxKeyUser      ctxKey = iota
	ctxKeyRequestID ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type Handler struct {
	service        *service.Service
	sessionStore   sessions.Store
	Router         *mux.Router
	logger         *logrus.Logger
	sessionManager authservice.AuthServiceClient
}

func NewHandler(service *service.Service, sessionStore sessions.Store, sessionManager authservice.AuthServiceClient) *Handler {
	return &Handler{
		service:        service,
		sessionStore:   sessionStore,
		Router:         mux.NewRouter(),
		logger:         logrus.New(),
		sessionManager: sessionManager,
	}
}

func (h *Handler) InitHandler() {
	api := h.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(handlers.CORS(
		handlers.ExposedHeaders([]string{"Set-Cookie"}),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "content-type", "Origin", "Accept", "X-Requested-With"}),
		handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost"}),
	))
	api.Use(h.setRequestID)
	api.Use(h.logRequest)
	api.HandleFunc("/register", h.handleRegister()).Methods("POST")
	api.HandleFunc("/signin", h.handleSignIn()).Methods("POST")

	private := api.PathPrefix("/private").Subrouter()
	private.Use(h.AuthenticateUser)
	private.HandleFunc("/whoami", h.handleWhoami()).Methods("GET")
	private.HandleFunc("/signout", h.handleSignOut()).Methods("GET")

	product := private.PathPrefix("/product").Subrouter()
	product.HandleFunc("/product", h.handleProductCreate()).Methods("POST")
	product.HandleFunc("/product", h.handleProductList()).Methods("GET")
	product.HandleFunc("/product/{id}", h.handleProductOptions()).Methods("OPTIONS")
	product.HandleFunc("/product/{id}", h.handleProductGet()).Methods("GET")
	product.HandleFunc("/product/{id}", h.handleProductUpdate()).Methods("PUT")
	product.HandleFunc("/product/{id}", h.handleProductDelete()).Methods("DELETE")
	product.HandleFunc("/category/get_categories", h.handleProductCategoryGet()).Methods("GET")
	product.HandleFunc("/material/get_materials", h.handleProductMaterialGet()).Methods("GET")
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
