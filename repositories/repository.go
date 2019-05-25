package repositories

import (
	"database/sql"

	"../configs"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New() *Repository {
	return &Repository{db: configs.DatabaseSetup()}
}
