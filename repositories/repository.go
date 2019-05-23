package repositories

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var initialized uint32

type Repository struct {
	db *sql.DB
}
