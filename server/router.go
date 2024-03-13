package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

// CREATE   /create
// JOIN  /{tableid}/{join}/{id}
// STATUS /{tableid}
// HIT  /{tableid}/{id}/hit
// STAND /{tableid}/{id}/stand
// DELETE /delete/{tableid}
// Handler holds the HTTP handlers for the API

// NewRouter initializes and returns the HTTP router
func NewRouter(handler *Handler) http.Handler {
	router := mux.NewRouter()
	// Create a new table
	router.HandleFunc("/create", handler.CreateTableHandler).Methods("GET")

	// Get table details
	router.HandleFunc("/delete/{tableID}", handler.DeleteTableHandler).Methods("GET")

	return router
}
