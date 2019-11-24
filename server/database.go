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
	//open the DB
	data, err := sql.Open("sqlite3", PROJECT_FOLDER + "/server/data/database.db")
	if err != nil {
		log.Print("Error loading database.")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	//create tables if they dont exist
	statement, err := data.Prepare(
		"CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, uname TEXT, name TEXT, pass TEXT)")
	if err != nil {
		log.Print("Error preapring init users")
		log.Fatal("Error was: ", err.Error())
		return nil
	}
	_, err = statement.Exec()
	if err != nil {
		log.Print("Error initialising users table.")
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
	//return the DB controller
	return data
}

func dbUsersLog( db *sql.DB ) {
	log.Print("Logging all users...")
	//get all users
	data, err := db.Query("SELECT * from users")
	if err != nil {
		log.Print("Error logging database.")
		log.Print("Error was: ", err.Error())
	}

	var userTmp DbUser

	//for each user dump it's information
	for data.Next() {
		if data == nil{
			return
		}
		_ = data.Scan(&userTmp.id, &userTmp.name, &userTmp.uname, &userTmp.pass)
		log.Print(strconv.Itoa(userTmp.id) + " : " + userTmp.name + " : " + userTmp.uname + " " + userTmp.pass)
	}
}

func dbSamplesLog( db *sql.DB ) {
	log.Print("Logging all samples...")
	//get all users
	data, err := db.Query("SELECT * from samples")
	if err != nil {
		log.Print("Error logging database.")
		log.Print("Error was: ", err.Error())
	}

	var sampleTmp WsData
	var id int
	//for each user dump it's information
	for data.Next() {
		if data == nil{
			return
		}
		_ = data.Scan(&id, &sampleTmp.Timestamp, &sampleTmp.Value)
		log.Print(strconv.Itoa(id) + " : " + sampleTmp.Timestamp + " : " + strconv.Itoa(sampleTmp.Value))
	}
}

func dbUserInsert(db *sql.DB, u DbUser){
	//insert a new user
	st, _ := db.Prepare("INSERT INTO users (uname, name, pass) values (?, ?, ?)")
	_, _ = st.Exec(u.uname, u.name, u.pass)
}

func dbEntryInsert(db *sql.DB, u []WsData){
	//insert a new sample
	qr := "INSERT INTO samples (timestamp, value) values "
	for i, l := range u {
		if i != 0 {
			qr += ","
		}
		qr = qr + "('"+l.Timestamp+"','"+strconv.Itoa(l.Value)+"')"
	}
	log.Print(qr)
	st, _ := db.Prepare(qr)
	_, _ = st.Exec()
}

func dbGetSamples(db *sql.DB) []WsData {

	var ret []WsData
	//get all samples
	data, err := database.Query(
		"select distinct * from samples where timestamp!='' and value!=0")
	if err != nil {
		log.Print("Could not get samples.")
		log.Print("Error was: ", err.Error())
		return nil
	}

	//for each sample save it to slice
	for data.Next(){
		var toIns WsData
		var dummy int
		err = data.Scan(&dummy, &toIns.Timestamp, &toIns.Value)
		if err != nil{
			log.Print("Could not parse row of samples.")
			log.Print("Error was: ", err.Error())
			return nil
		}
		ret = append(ret, toIns)
	}
	//return sample slice
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

