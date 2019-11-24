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

func (cl Client) writeSocket(){

	defer func() {
		_ = cl.connection.Close()
	}()

	//wait for samples
	for {
		select {
		//on new sample
		case msg, ok := <- cl.writeChannel:
			if !ok {
				//channel closed by main
				return
			}
			//send it to the client
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
func (cl Client) readSocket(){
	for{
		select {
			default:
				var msg []byte
				_, msg, _ = cl.connection.ReadMessage()
				if string(msg) == "ping" {
					cl.writeChannel <- "pong"
				}
		}
	}
}

func clientServeWs(w http.ResponseWriter, r *http.Request) {
	//create the writeSocket
	conn, err := UPGRADER.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error on opening websocket connection.")
		log.Print("Error was: ", err.Error())
	}

	//new client object
	newClient := Client{
		conn,
		make(chan string),
	}
	clients = append(clients, newClient)
	go newClient.writeSocket()

	newClient.writeChannel <- ""
	//send history of samples to new client
	oldData := dbGetSamples(database)
	for _, s := range oldData {
		b, _ := json.Marshal(s)
		log.Print(string(b))
		newClient.writeChannel <- string(b)
	}

	//run concurrently the client to process samples!
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

