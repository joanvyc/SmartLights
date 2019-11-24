package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Content string `json:"flow_label"`
}

type WsData struct {
	Timestamp string `json:"timestamp"`
	Value int `json:"value"`
}
type ApiData struct {
	TimeStamp int `json:"timestamp"`
	Samples []int `json:"samples"`
}

func apiResponse(w http.ResponseWriter, code int, content string) error {
	w.Header().Set("Content-Type", "application/json")
	resp := ApiResponse{
		Code:    code,
		Content: content,
	}
	return json.NewEncoder(w).Encode(resp)
}

func apiHomeEndpoint(w http.ResponseWriter, r *http.Request) {
	if err := apiResponse(w, 0, "OK"); err != nil {
		log.Print("Error apiHomeEndpoint failed to respond.")
		log.Print("Error was: ", err.Error())
	}
}

func apiRecieveDataEndpoint(w http.ResponseWriter, r *http.Request) {
	//current time stamp
	t := time.Now()
	//parse response into string
	bytesToSend, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Could not read request body")
		log.Print("Error was: ", err.Error())
		return
	}
	log.Print(string(bytesToSend))
	// parse into ApiData
	var newData ApiData
	err = json.Unmarshal(bytesToSend, &newData)
	if err != nil {
		log.Print("Error could not parse json ApiData")
		log.Print("Error was: ", err.Error())
		return
	}
	var entry WsData
	//for each sample, insert into database, notify current clients of new samples
	var entryArray []WsData
	for i, v := range newData.Samples {
		var toSend []byte
		entry.Timestamp = t.Add(time.Duration(20-i)*-4*time.Second).String()
		entry.Value = v
		toSend, _ = json.Marshal(&entry)
		log.Print(string(toSend))
		for _, cl := range clients{
			cl.writeChannel <- string(toSend)
		}
		entryArray = append(entryArray, entry)
	}
	dbEntryInsert(database, entryArray)

	//if newData.Samples[0] < 500 {
	//	openPowerSocket()
	//}else{
	//	closePowerSocket()
	//}

}
