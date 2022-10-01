package apiserver

import (
	"github.com/VladimirBlinov/MarketPlace/internal/handler"
	"github.com/VladimirBlinov/MarketPlace/internal/service"
	"github.com/VladimirBlinov/MarketPlace/internal/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
)

//const (
//	SessionName            = "MarketPlace"
//	ctxKeyUser      ctxKey = iota
//	ctxKeyRequestID ctxKey = iota
//)
//
//var (
//	errIncorectEmailOrPassword = errors.New("Incorrect email or password")
//	errNotAuthenticated        = errors.New("Not authenticated")
//)
//
//type ctxKey int8

type server struct {
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	service      service.Service
	handler      handler.Handler
}

func newServer(store store.Store, sessionStore sessions.Store, service service.Service, handler handler.Handler) *server {
	s := &server{
		logger:  logrus.New(),
		handler: handler,
	}

	s.handler.InitHandler()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.Router.ServeHTTP(w, r)
}
