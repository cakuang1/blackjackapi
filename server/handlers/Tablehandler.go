package server

import (
	"blackjackapi/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Handler struct {
	Client  *redis.Client
	Context context.Context
}

// NewHandler initializes and returns a new Handler instance
func NewHandler(tableStore *redis.Client) *Handler {
	return &Handler{
		Client:  tableStore,
		Context: context.Background(),
	}
}

// CREATE TABLE
func (h *Handler) CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	id := uuid.New().String()
	table := models.NewTable(id)
	err := models.SaveTable(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	text := fmt.Sprintf("Table has been created. Table id : %s", id)
	fmt.Fprintf(w, text)
	fmt.Fprint(w, table.GetBoardText())
}

// DELETE TABLE
func (h *Handler) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// Get the tableID from the URL path parameters
	vars := mux.Vars(r)
	id, ok := vars["tableID"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	// If the table does not exist, return a 404 Not Found error
	err := h.Client.Del(h.Context, id).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Table deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// START HANDLER
func (h *Handler) StartTableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// Get the tableID from the URL path parameters
	vars := mux.Vars(r)
	id, ok := vars["tableID"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}
	table, err := models.GetTable(h.Context, id, h.Client)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	table.TableClear()
	models.SaveTable(h.Context, table, h.Client)
	fmt.Fprint(w, "Game has started")
	fmt.Fprint(w, table.GetBoardText())
}

// STATUS HANDLER

func (h *Handler) GetTableDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	vars := mux.Vars(r)
	id, ok := vars["tableID"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
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
	// Send the game board string to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, table.GetBoardText())
}

// ADD PLAYER HANDLER
func (h *Handler) AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(name) < 1 || len(name) > 10 {
		http.Error(w, "Please keep names within 1 to 10 characters please", http.StatusBadRequest)
		return
	}
	// Get the table details from the models package
	table, err := models.GetTable(h.Context, id, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if the game is in Play
	if !table.Status {
		http.Error(w, "Game is in play. Please wait till the game is over to add a new player", http.StatusBadRequest)
		return
	}
	// Check if the name is already used
	for _, v := range table.Players {
		if v.Name == name {
			http.Error(w, "Name has already been taken. Please choose another name.", http.StatusBadRequest)
		}
	}
	// Create the new player
	NewPlayer := models.NewPlayer(name)
	// Checks have been made, Players can start being added
	table.AddPlayer(NewPlayer)
	models.SaveTable(h.Context, table, h.Client)
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("New player %s has been added to table %s", name, table.ID)
	fmt.Fprintf(w, text)
	fmt.Fprintf(w, table.GetBoardText())
}

// DELETE PLAYER HANDLER
func (h *Handler) DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get the table details from the models package
	table, err := models.GetTable(h.Context, id, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if the game is in Play
	if !table.Status {
		http.Error(w, "Game is in play. Please wait till the game is over to delete a player", http.StatusBadRequest)
		return
	}
	if !table.PlayerIsin(name) {
		http.Error(w, fmt.Sprintf("Player %s is not in table %s", name, table.ID), http.StatusBadRequest)
	}
	// Check if the name is already used
	for i, v := range table.Players {
		if v.Name == name {
			// Delete the element at index i by slicing the slice
			table.Players = append(table.Players[:i], table.Players[i+1:]...)
			break // Exit the loop since the element is deleted
		}
	}
	w.WriteHeader(200)
	fmt.Fprintf("")
	fmt.Fprintf("")
}
