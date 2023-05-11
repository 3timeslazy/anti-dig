package db

import "github.com/3timeslazy/anti-dig/example/config"

type DB struct{}

func NewDB(config.Config) (DB, error) {
	return DB{}, nil
}
