package main

import (
	"log"
	"net/http"
	"path"
)

func openPowerSocket(){
	resp, err := http.Get(path.Join("https://maker.ifttt.com/trigger/open/with/key/", IFTT_KEY))
	if err != nil {
		log.Print("Error on openning power socket.")
		log.Print("Error was: ", err.Error())
		return
	}

	log.Print(resp.StatusCode)
}

func closePowerSocket(){
	resp, err := http.Get(path.Join("https://maker.ifttt.com/trigger/close/with/key/", IFTT_KEY))
	if err != nil {
		log.Print("Error on closing power socket.")
		log.Print("Error was: ", err.Error())
		return
	}

	log.Print(resp.StatusCode)
}

