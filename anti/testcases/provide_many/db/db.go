package db

import "github.com/3timeslazy/anti-dig/anti/testcases/provide_many/config"

type DB struct{}

func New(config config.Config) DB {
	return DB{}
}
