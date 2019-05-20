package helpers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"../models"
)

func AnalizeSslLabs(host string) models.SslLabsResponse {
	var sslLabsResponseObject models.SslLabsResponse

	response, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + host)

	if err != nil {
		log.Println(err)
	} else {
		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Println(err)
		} else {
			json.Unmarshal(responseData, &sslLabsResponseObject)
		}
	}

	return sslLabsResponseObject
}
