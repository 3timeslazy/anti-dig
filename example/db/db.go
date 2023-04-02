package db

type DB struct{}

func NewDB() (DB, error) {
	return DB{}, nil
}
