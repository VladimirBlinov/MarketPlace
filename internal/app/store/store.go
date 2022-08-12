package store

type Store interface {
	User() UserRepo
	Product() ProductRepo
}
