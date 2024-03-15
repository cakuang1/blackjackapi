package server

import (
	"blackjackapi/server/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// maybe tableid first makes more sense

// CREATE   /create
// DELETE /{tableid}/delete/
// STATUS /{tableid}
// JOIN  /{tableid}/join/{id}
// LEAVE /{tableid}/leave/{id}
// HIT  /{tableid}/{id}/hit
// STAND /{tableid}/{id}/stand

// Handler holds the HTTP handlers for the API

// NewRouter initializes and returns the HTTP router
func NewRouter(handler *handlers.Handler) http.Handler {
	router := mux.NewRouter()
	//CREATE
	router.HandleFunc("/create", handler.CreateTableHandler).Methods("GET")
	//DELETE
	router.HandleFunc("/{tableID}/delete", handler.DeleteTableHandler).Methods("GET")
	//STATUS
	router.HandleFunc("/{tableID}/status", handler.GetTableDetailsHandler).Methods("GET")
	//START
	router.HandleFunc("/{tableID}/start", handler.StartTableHandler)
	// JOIN
	router.HandleFunc("/{tableID}/join/{name}", handler.AddPlayerHandler).Methods("GET")
	// LEAVE
	router.HandleFunc("/{tableID}/leave/{name}", handler.DeletePlayerHandler).Methods("GET")
	// HIT
	router.HandleFunc("/{tableID}/hit/{name}", handler.HitPlayerHandler).Methods("GET")
	// STAND
	router.HandleFunc("/{tableID}/stand/{name}", handler.StandPlayerHandler).Methods("GET")

	return router
}
