package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
var clients []Client

func main() {
	log.Print("Init server...")
	database = dbInit()
	router := mux.NewRouter()
	router.HandleFunc("/ws", clientServeWs)
	router.HandleFunc("/login", clientLogin).Methods("POST")
	router.HandleFunc("/api", apiHomeEndpoint).Methods("GET")
	router.HandleFunc("/api/postData", apiRecieveDataEndpoint).Methods("POST")
	router.HandleFunc("/", clientServeLogin).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}
