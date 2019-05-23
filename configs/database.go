package configs

import (
	"database/sql"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DatabaseSetup() *sql.DB {
	dbAddr := "postgresql://gaea@localhost:26257/servers?sslmode=disable"
	log.Println("Connecting to database:" + dbAddr)
	dbConnection, err := sql.Open("postgres", dbAddr)

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return dbConnection
}
