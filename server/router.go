package server

import (
	"blackjackapi/server/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// maybe tableid first makes more sense

// CREATE   /create
// DELETE /{tableid}/delete/
// START /{tableid}
// JOIN  /{tableid}/join/{id}
// DROP  /{tableid}/{id}/{column}/drop
// CONNECT /{tableid}/connect

// Handler holds the HTTP handlers for the API

// NewRouter initializes and returns the HTTP router
func NewRouter(handler *handlers.Handler) http.Handler {
	router := mux.NewRouter()
	// STATIC
	router.Handle("/", http.FileServer(http.Dir("./static")))
	//CREATE
	router.HandleFunc("/create", handler.CreateTableHandler).Methods("GET")
	//DELETE
	router.HandleFunc("/{tableID}/delete", handler.DeleteTableHandler).Methods("GET")
	//START
	router.HandleFunc("/{tableID}/start", handler.StartGameHandler)
	// JOIN
	router.HandleFunc("/{tableID}/{name}/join", handler.JoinTableHandler).Methods("GET")
	// LEAVE
	router.HandleFunc("/{tableID}/{name}/leave", handler.LeaveTableHandler).Methods("GET")
	// DROP
	router.HandleFunc("/{tableID}/{name}/{column}/drop", handler.DropPieceHandler).Methods("GET")
	// CONNECT
	router.HandleFunc("/{tableID}/connect", handler.KafkaSSEHandler).Methods("GET")

	return router
}
