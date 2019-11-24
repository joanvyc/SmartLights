package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

type DbUser struct {
	id int
	uname string
	name string
	pass string
}

var database *sql.DB

func dbInit() *sql.DB{
	data, err := sql.Open("sqlite3", PROJECT_FOLDER + "/server/data/database.db")
	if err != nil {
		log.Print("Error loading database.")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	statement, err := data.Prepare(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uname TEXT, name TEXT, pass TEXT)")
	if err != nil {
		log.Print("Error init users")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	_, err = statement.Exec()
	if err != nil {
		log.Print("Error initialising  database.")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	statement, _ = data.Prepare(
		"CREATE TABLE IF NOT EXISTS samples (id INTEGER PRIMARY KEY, timestamp TEXT, value INTEGER)")
	_, err = statement.Exec()
	if err != nil {
		log.Print("Error initialising  database.")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	return data
}

func dbLog( db *sql.DB ) {

	data, err := db.Query("SELECT * from users")
	if err != nil {
		log.Print("Error logging database.")
		log.Print("Error was: ", err.Error())
	}

	var userTmp DbUser

	for data.Next() {
		if data == nil{
			return
		}
		_ = data.Scan(&userTmp.id, &userTmp.name, &userTmp.uname, &userTmp.pass)
		log.Print(strconv.Itoa(userTmp.id) + " : " + userTmp.name + " : " + userTmp.uname + " " + userTmp.pass)
	}
}

func dbUserInsert(db *sql.DB, u DbUser){
	st, _ := db.Prepare("INSERT INTO users (uname, name, pass) values (?, ?, ?)")
	_, _ = st.Exec(u.uname, u.name, u.pass)
}

func dbEntryInsert(db *sql.DB, u WsData){
	st, _ := db.Prepare("INSERT INTO samples (timestamp, value) values (?, ?)")
	_, _ = st.Exec(u.Timestamp, u.Value)
}

func dbGetSamples(db *sql.DB) []WsData {
	var ret []WsData

	data, err := database.Query(
		"select * from samples")
	if err != nil {
		log.Print("Could not get samples.")
		log.Print("Error was: ", err.Error())
		return nil
	}

	for data.Next(){
		var toIns WsData
		err = data.Scan(toIns.Timestamp, toIns.Value)
		if err != nil{
			log.Print("Could not parse row of samples.")
			log.Print("Error was: ", err.Error())
			return nil
		}
		ret = append(ret, toIns)
	}
	return ret
}

func userExists(uid string, password string) bool{
	st := "select * from users where uname='" + uid + "' and pass='" + password + "';"
	data, err := database.Query(st)
	if err != nil {
		log.Print("Error querying db.")
		log.Print("Error was: ", err.Error())
		return false
	}
	return data.Next()

}

