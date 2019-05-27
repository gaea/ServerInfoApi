package routers

import (
	"../controllers"

	"github.com/go-chi/chi"
)

func AnalizeServerApi(router chi.Router) {
	router.Get("/analize", controllers.Analize)
	router.Get("/searchhistory", controllers.History)
	router.Get("/detail", controllers.HostDetail)
}
