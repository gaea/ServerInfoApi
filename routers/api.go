package routers

import (
	"../configs"
	"../controllers"

	"github.com/go-chi/chi"
)

func AnalizeServerApi(router chi.Router) {
	dbAddr := "postgresql://gaea@localhost:26257/servers?sslmode=disable"
	db := configs.DatabaseSetup(dbAddr)
	//defer db.Close()

	apiController := controllers.ApiController(db)

	router.Get("/servers", apiController.ServerAnalize)
	router.Get("/searchhistory", apiController.History)
}
