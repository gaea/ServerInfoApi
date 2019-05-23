package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"../helpers"
	"../models"
	"../repositories"
)

// Based on: https://github.com/cockroachdb/examples-orms/blob/master/go/gorm/server.go
func jsonResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func jsonErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func searchAnalize(host string) (models.HostInfo, error) {
	sslLabsResponseObject := helpers.AnalizeSslLabs(host)
	numServers := len(sslLabsResponseObject.Servers)
	scrapedTitle, scrapedLogo := helpers.AnalizeScraping(host)

	if sslLabsResponseObject.Status != "READY" {
		return models.HostInfo{}, errors.New("Analysis is not ready, please try later")
	}

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

		hostInfo.Servers[i] = serverInfo
	}

	hostInfo.SetSslGrade()

	return hostInfo, nil
}

func Search(w http.ResponseWriter, r *http.Request) {
	hostInfoRepository := repositories.HostInfoRepository()
	serverInfoRepository := repositories.ServerInfoRepository()

	host := r.URL.Query().Get("host")

	hostInfo, err := hostInfoRepository.DetailHostInfo(host)

	if err == nil {
		if hostInfo.Host == "" {
			hostInfo, err = searchAnalize(host)

			if err == nil {
				err = serverInfoRepository.CreateHostInfo(&hostInfo)
			}
		} else {
			duration := time.Since(hostInfo.LastChecked)

			if duration.Hours() >= 1 {
				var hostInfoUpdated models.HostInfo
				hostInfoUpdated, err = searchAnalize(host)

				if err == nil {
					err = hostInfoRepository.CheckHostInfo(hostInfo, &hostInfoUpdated)

					if err == nil {
						hostInfo = hostInfoUpdated
					}
				}
			}
		}
	}

	if err == nil {
		jsonResponse(w, hostInfo)
	} else {
		jsonErrorResponse(w, err)
	}
}

func History(w http.ResponseWriter, r *http.Request) {
	hostInfoRepository := repositories.HostInfoRepository()

	hostInfos, err := hostInfoRepository.ListHostInfo()

	if err != nil {
		jsonErrorResponse(w, err)
	}

	hostInfoItems := struct {
		Items []models.HostInfo `json:"items"`
	}{
		Items: hostInfos,
	}

	jsonResponse(w, hostInfoItems)
}
