package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)
type Client struct {
	connection *websocket.Conn
	writeChannel chan string
}

func (cl Client) socket(){
	defer func() {
		_ = cl.connection.Close()
	}()

	for {
		select {
			case msg, ok := <- cl.writeChannel:
				if !ok {
					//channel closed by main
					return
				}

				w, err := cl.connection.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				_, err = w.Write([]byte(msg))
				if err != nil {
					return
				}

		}
	}
}

func clientServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error on opening websocket connection.")
		log.Print("Error was: ", err.Error())
	}

	newClient := Client{
		conn,
		make(chan string),
	}
	oldData := dbGetSamples(database)
	for _, s := range oldData {
		b, _ := json.Marshal(s)
		newClient.writeChannel <- string(b)
	}

	clients = append(clients, newClient)
	go newClient.socket()
}

func clientLogin(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Print("Can't parse login.")
		return
	}

	if userExists(r.PostFormValue("userid"), r.PostFormValue("pswrd")) {
		http.ServeFile(w, r, PROJECT_FOLDER + "/frontend/graph.html")

	}else{
		http.ServeFile(w, r, PROJECT_FOLDER + "/frontend/index.html")
	}

}

func clientServeLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, PROJECT_FOLDER + "/frontend//index.html")
}

