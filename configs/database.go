package configs

import (
	"log"

	"../models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DatabaseSetup(dbAddr string) *gorm.DB {
	log.Println("Connecting to database:" + dbAddr)
	db, err := gorm.Open("postgres", dbAddr)

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.HostInfo{}, &models.ServerInfo{}, &models.History{})

	return db
}
