package routers

import (
	"../controllers"
	"github.com/go-chi/chi"
)

func Pages(router chi.Router) {
	router.Get("/", controllers.MainPage)
}
