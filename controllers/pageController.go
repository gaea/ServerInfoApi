package controllers

import "net/http"

func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Server Analize API"))
}
