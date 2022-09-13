package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/VladimirBlinov/MarketPlace/internal/app/store/sqlstore"
	"github.com/VladimirBlinov/MarketPlace/internal/service"
	"github.com/gorilla/sessions"
)

func Start(config *Config) error {
	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	service := service.NewService(store)
	srv := newServer(store, sessionStore, *service)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
