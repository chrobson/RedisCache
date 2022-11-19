package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	database   Database
}

func NewAPIServer(listenAddr string, database Database) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		database:   database,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/{id}", makeHTTPHandleFunc(s.handleGetUserById))
	http.ListenAndServe(s.listenAddr, router)
}

func RenderJson(w http.ResponseWriter, val interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(val)
}

func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]

	redis, err := NewRedis()
	if err != nil {
		log.Fatalf("Could not initialize Redis client %s", err)
	}

	val, err := redis.GetName(r.Context(), id)
	if err == nil {
		return RenderJson(w, &val, http.StatusOK)
	}

	person, err := s.database.GetUserById(id)

	if err != nil {
		return RenderJson(w, &ApiError{Error: err.Error()}, http.StatusInternalServerError)

	}

	_ = redis.SetName(r.Context(), *person)

	return RenderJson(w, &person, http.StatusOK)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			RenderJson(w, ApiError{Error: err.Error()}, http.StatusBadRequest)
		}
	}
}
