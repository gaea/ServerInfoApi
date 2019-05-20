package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"../helpers"
	"../models"
	"github.com/jinzhu/gorm"
)

type Api struct {
	db *gorm.DB
}

func ApiController(db *gorm.DB) *Api {
	return &Api{db: db}
}

// Based on: https://github.com/cockroachdb/examples-orms/blob/master/go/gorm/server.go
func jsonResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func (api *Api) ServerAnalize(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")

	sslLabsResponseObject := helpers.AnalizeSslLabs(host)
	numServers := len(sslLabsResponseObject.Servers)
	scrapedTitle, scrapedLogo := helpers.AnalizeScraping(host)

	var hostInfo = models.HostInfo{
		Title:   scrapedTitle,
		Logo:    scrapedLogo,
		Host:    sslLabsResponseObject.Host,
		Servers: make([]models.ServerInfo, numServers),
	}

	hostInfo.SetStatus(sslLabsResponseObject.Status)

	for i := 0; i < numServers; i++ {
		sslLabServer := sslLabsResponseObject.Servers[i]
		whoisCountry, whoisOwner := helpers.AnalizeWhois(sslLabServer.Address)

		var serverInfo = models.ServerInfo{
			Address:    sslLabServer.Address,
			ServerName: sslLabServer.ServerName,
			SslGrade:   sslLabServer.SslGrade,
			Country:    whoisCountry,
			Owner:      whoisOwner,
		}

		if err := api.db.Create(&serverInfo).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		hostInfo.Servers[i] = serverInfo
	}

	hostInfo.SetSslGrade()

	if err := api.db.Create(&hostInfo).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		var history = models.History{
			Host:     hostInfo.Host,
			Date:     time.Now(),
			HostInfo: hostInfo,
		}

		api.db.Create(&history)
	}

	jsonResponse(w, hostInfo)
}

func (api *Api) History(w http.ResponseWriter, r *http.Request) {
	var histories []models.History
	if err := api.db.Preload("HostInfo.Servers").Find(&histories).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		jsonResponse(w, histories)
	}
}
