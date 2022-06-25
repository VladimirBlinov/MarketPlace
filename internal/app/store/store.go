package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Store
type Store struct {
	config      *Config
	db          *sql.DB
	logger      *logrus.Logger
	userRepo    *UserRepo
	productRepo *ProductRepo
}

// New...
func New(config *Config) *Store {
	return &Store{
		config: config,
		logger: logrus.New(),
	}
}

// Open
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.logger.Info("DataBase connected")
	s.db = db

	return nil
}

// Close
func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{
		store: s,
	}
	return s.userRepo
}

func (s *Store) Product() *ProductRepo {
	if s.productRepo != nil {
		return s.productRepo
	}

	s.productRepo = &ProductRepo{
		store: s,
	}
	return s.productRepo
}

func (s *Store) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)

	return nil
}
