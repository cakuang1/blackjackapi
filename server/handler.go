package server

import (
	"blackjackapi/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
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

	// Get the ID of the table to delete from the request parameters or body
	id := r.FormValue("tableID")
	if id == "" {
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

// JOIN TABLE
