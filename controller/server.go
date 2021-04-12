package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/gorilla/mux"
)

type Server struct {
	store      db.Store
	router     *mux.Router
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func getVillager(w http.ResponseWriter, r *http.Request) {
	ACNHClient := http.Client{ Timeout: time.Second * 10 }
	params := mux.Vars(r)
	url := "https://ac-vill.herokuapp.com/villagers?name=" + params["name"]

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := ACNHClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var villager []db.Villager
	jsonErr := json.Unmarshal(body, &villager)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(villager)
}

func (server *Server) setupRouter() {
	port := "7890"
	fmt.Printf("Starting server on port %q...\n", port)

	// Initialize router
	router := mux.NewRouter()
	
	// Route handlers / endpoints
	router.HandleFunc("/api/villagers/{name}", getVillager).Methods("GET")
	
	log.Fatal(http.ListenAndServe(":"+port, router))
}
