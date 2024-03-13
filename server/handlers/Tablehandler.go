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
	w.Header().Set("Content-Type", "application/json")
	id := uuid.New().String()
	table := models.NewTable(id)
	err := models.SaveTable(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Message string `json : "message"`
		ID      string `json :  "id"`
	}{
		Message: "Table created successfully. Use the ID to join the table and play",
		ID:      id,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the response JSON to the client
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// DELETE TABLE
func (h *Handler) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
		http.Error(w, "Game is in play. Please wait till the game is over", http.StatusBadRequest)
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
}
