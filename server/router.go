package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

// maybe tableid first makes more sense

// CREATE   /create
// DELETE /{tableid}/delete/
// STATUS /{tableid}
// JOIN  /{tableid}/{join}/{id}
// HIT  /{tableid}/{id}/hit
// STAND /{tableid}/{id}/stand

// Handler holds the HTTP handlers for the API

// NewRouter initializes and returns the HTTP router
func NewRouter(handler *Handler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/create", handler.CreateTableHandler).Methods("GET")
	router.HandleFunc("/{tableID}/delete/", handler.DeleteTableHandler).Methods("GET")
	router.HandleFunc("/{tableID}/status", handler.GetTableDetailsHandler).Methods("GET")

	return router
}
