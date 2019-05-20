package main

import (
	"./configs"
	"./routers"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	router := configs.ChiRouterSetup()

	routers.AnalizeServerApi(router)
	routers.Pages(router)

	serverAddress := "localhost:3000"
	configs.ServerSetup(serverAddress, router)
}
