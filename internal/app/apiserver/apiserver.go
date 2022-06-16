package apiserver

import (
	"io"
	"net/http"

	"github.com/VladimirBlinov/MarketPlace/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	serverConfig *Config
	logger       *logrus.Logger
	router       *mux.Router
	store        *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		serverConfig: config,
		logger:       logrus.New(),
		router:       mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.ConfigureLogger(); err != nil {
		return err
	}

	s.ConfigureRouter()

	if err := s.ConfigureStore(); err != nil {
		return err
	}

	s.logger.Info("Starting API server")
	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)
}

func (s *APIServer) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.serverConfig.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)

	return nil
}

// Config Router ...
func (s *APIServer) ConfigureRouter() {
	s.router.HandleFunc("/allproducts", s.HandleGetAllProducts())
	s.router.HandleFunc("/product", s.HandleGetProductByID())
}

// Config Store ...
func (s *APIServer) ConfigureStore() error {
	serverStore := store.New(s.serverConfig.Store)

	if err := serverStore.Open(); err != nil {
		return err
	}

	s.store = serverStore

	return nil
}

func (s *APIServer) HandleGetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "All Products here")
	}
}

func (s *APIServer) HandleGetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Products by ID")
	}
}
