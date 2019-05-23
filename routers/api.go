package routers

import (
	"../controllers"

	"github.com/go-chi/chi"
)

func AnalizeServerApi(router chi.Router) {
	router.Get("/servers", controllers.Search)
	router.Get("/searchhistory", controllers.History)
}
