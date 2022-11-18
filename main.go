package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type Person struct {
	Id         string `json:"id"`
	FirstName  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Mail       string `json:"mail"`
	Gender     string `json:"gender"`
}

type Error struct {
	Message string `json:"error"`
}

const (
	host     = "127.0.0.1"
	port     = 49153
	user     = "postgres"
	password = "postgrespw"
	dbname   = "people"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	renderJSON := func(w http.ResponseWriter, val interface{}, statusCode int) {
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(val)
	}

	id := html.EscapeString(strings.Split(r.URL.Path, "/")[1])

	redis, err := NewRedis()
	if err != nil {
		log.Fatalf("Could not initialize Redis client %s", err)
	}

	val, err := redis.GetName(r.Context(), id)
	if err == nil {
		renderJSON(w, &val, http.StatusOK)
		//w.Write(val)
		return
	}

	db := OpenConnection()
	row := db.QueryRow("SELECT * FROM people WHERE id=$1", id)

	var person Person
	err = row.Scan(&person.Id, &person.FirstName, &person.Secondname, &person.Mail, &person.Gender)

	defer db.Close()

	if err != nil {
		renderJSON(w, &Error{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	_ = redis.SetName(r.Context(), person)

	renderJSON(w, &person, http.StatusOK)
}

func main() {

	http.HandleFunc("/", GetUserById)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
