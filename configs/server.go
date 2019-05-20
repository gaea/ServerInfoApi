package configs

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ServerSetup(serverAddress string, router chi.Router) {
	log.Println("Starting server at:" + serverAddress)
	serverError := http.ListenAndServe(serverAddress, router)

	if serverError != nil {
		log.Println(serverError)
		os.Exit(1)
	}
}
