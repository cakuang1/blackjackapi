package server

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/wafregional/wafregionaliface"
	"github.com/gorilla/mux"
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
func NewRouter(handler *Handler) http.Handler {
	router := mux.NewRouter()
	//CREATE
	router.HandleFunc("/create", handler.CreateTableHandler).Methods("GET")
	//DELETE
	router.HandleFunc("/{tableID}/delete/", handler.DeleteTableHandler).Methods("GET")
	//STATUS
	router.HandleFunc("/{tableID}/status", handler.GetTableDetailsHandler).Methods("GET")
	// JOIN
	router.HandleFunc("/{tableID}/join/{name}", handler.GetTableDetailsHandler).Methods("GET")
	// LEAVE
	router.HandleFunc("/{tableID}/leave/{name}", handler.GetTableDetailsHandler).Methods("GET")
	// HIT
	router.HandleFunc("/{tableID}/hit/{name}", handler.GetTableDetailsHandler).Methods("GET")
	// STAND
	router.HandleFunc("/{tableID}/stand/{name}", handler.GetTableDetailsHandler).Methods("GET")

	return router
}
