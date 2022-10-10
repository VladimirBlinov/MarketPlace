package handler

import (
	"encoding/json"
	"errors"
	"github.com/VladimirBlinov/MarketPlace/internal/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
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
	Router       *mux.Router
	logger       *logrus.Logger
}

func NewHandler(service *service.Service, sessionStore sessions.Store) *Handler {
	return &Handler{
		service:      service,
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
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"})))
	h.Router.HandleFunc("/register", h.handleRegister()).Methods("POST")
	h.Router.HandleFunc("/signin", h.handleSignIn()).Methods("POST")

	private := h.Router.PathPrefix("/private").Subrouter()
	private.Use(h.AuthenticateUser)
	private.HandleFunc("/whoami", h.handleWhoami()).Methods("GET")
	private.HandleFunc("/signout", h.handleSignOut()).Methods("GET")

	product := private.PathPrefix("/product").Subrouter()
	product.HandleFunc("/create", h.handleProductCreate()).Methods("POST")
	product.HandleFunc("/product/{id}", h.handleProductGetProductById()).Methods("GET")
	product.HandleFunc("/update/{id}", h.handleProductUpdate()).Methods("POST")
	product.HandleFunc("/all", h.handleProductList()).Methods("GET")
	product.HandleFunc("/category/get_categories", h.handleProductCategoryGetAll()).Methods("GET")
	product.HandleFunc("/material/get_materials", h.handleProductGetMaterials()).Methods("GET")
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
