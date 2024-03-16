package handlers

import (
	"blackjackapi/models"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// HIT
func (h *Handler) HitPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// Get the tableID from the URL path parameters
	vars := mux.Vars(r)
	id, ok := vars["tableID"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	name, ok := vars["name"]
	if !ok {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	table, err := models.GetTable(h.Context, id, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if table == nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}
	// Needs to be an actual player
	if !table.PlayerIsin(name) {
		http.Error(w, "Player is not registered for table {table name here}", http.StatusBadRequest)
	}
	// Game has not started !
	if !table.Status {
		http.Error(w, "Game has not yet started. Please start the game before making an player moves.", http.StatusBadRequest)
	}
	// Not your turn,then you can't go !
	if table.Players[table.Turn].Name != name {
		http.Error(w, fmt.Sprintf("The current  %s", table.ID), http.StatusBadRequest)
	}
	// Checks have been made. The turn can now move on. Slight edge case to consider if you are the last player.
	if table.Turn == len(table.Players) {
		table.Status = false
	} else {
		table.Turn++
	}
	models.SaveTable(h.Context, table, h.Client)
	w.WriteHeader(http.StatusOK)
	turnmessage := fmt.Sprintf("Player %s has decided to stand\n", name)
	final := turnmessage + table.GetBoardText()
	fmt.Fprint(w, final)
	models
}

// STAND
func (h *Handler) StandPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// Get the tableID from the URL path parameters
	vars := mux.Vars(r)
	id, ok := vars["tableID"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	name, ok := vars["name"]
	if !ok {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	table, err := models.GetTable(h.Context, id, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if table == nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}
	if !table.PlayerIsin(name) {
		http.Error(w, "Player is not registered for table {table name here}", http.StatusBadRequest)
	}
	// Game has not started !
	if !table.Status {
		http.Error(w, "Game has not yet started. Please start the game before making an player moves.", http.StatusBadRequest)
	}
	// Not your turn,then you can't go !
	if table.Players[table.Turn].Name != name {
		http.Error(w, fmt.Sprintf("Player is not registered for table %s", table.ID), http.StatusBadRequest)
	}

	// Checks have been made. The turn can now move on. Slight edge case to consider if you are the last player.
	if table.Turn == len(table.Players)-1 {
		table.Status = false
	} else {
		table.Turn++
	}

	models.SaveTable(h.Context, table, h.Client)
	w.WriteHeader(http.StatusOK)
	turnmessage := fmt.Sprintf("%s has decided to stand", name)
	final := turnmessage + table.GetBoardText()
	fmt.Fprint(w, final)

}
