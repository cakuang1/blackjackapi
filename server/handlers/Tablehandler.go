package handlers

import (
	"blackjackapi/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Creates a connect 4 table and sends the client the board id.
// CreateTableHandler creates a Connect 4 table and sends the client the board ID.
func (h *Handler) CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new Connect 4 table
	// Generate a unique ID for the table (you can use UUID or any other method)
	// For simplicity, let's assume the table ID is an integer incremented for each new table
	tableID := uuid.New().String()
	table := models.NewSession(tableID)
	// Set the table ID
	table.ID = tableID
	// Save the table to Redis
	err := models.SaveSession(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, "Trouble saving table. Please try again.", http.StatusInternalServerError)
		return
	}
	// Respond to the client with the table ID
	response := tableID
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

// DeleteTableHandler
// DeleteTableHandler deletes a Connect 4 table.
func (h *Handler) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]

	// Delete the table from Redis
	err := models.DeleteSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Failed to delete table from Redis", http.StatusInternalServerError)
		return
	}
	// Respond to the client
	response := fmt.Sprintf("Connect 4 table with ID %s deleted successfully. Thank you for deleting the table.", tableID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// JoinTableHandler handles requests to join a Connect 4 table
func (h *Handler) JoinTableHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]
	playerName := vars["name"]

	// Retrieve table from Redis
	table, err := models.GetSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Failed to retrieve table from Redis", http.StatusInternalServerError)
		return
	}
	// Check if the table is already full
	if len(table.Players) >= 2 {
		http.Error(w, "Table is already full", http.StatusConflict)
		return
	}
	PlayerIsin := false

	for _, v := range table.Players {
		if v.Name == playerName {
			PlayerIsin = true
		}
	}

	if PlayerIsin {
		http.Error(w, fmt.Sprintf("Name %s has already been taken", playerName), http.StatusConflict)
		return
	}
	// Create a new player
	player := models.NewPlayer(playerName)
	// Add the player to the table
	err = table.AddPlayer(player)
	if err != nil {
		http.Error(w, "Failed to add player to the table", http.StatusInternalServerError)
		return
	}
	// Save the updated table to Redis
	err = models.SaveSession(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, "Failed to save table to Redis", http.StatusInternalServerError)
		return
	}

	// Produce message to Kafka topic
	message := fmt.Sprintf("Player %s joined table", player.Name)
	message = table.StatusBoard(message) + table.StringBoard()

	err = models.ProduceMessage(h.KAFKAADDRESS, table.ID, message, h.KAFKAUSERNAME, h.KAFKAPASSWORD)
	if err != nil {
		http.Error(w, "Failed to produce message to Kafka topic", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// StartGameHandler handles requests to start a Connect 4 game
func (h *Handler) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]
	// Retrieve table from Redis
	table, err := models.GetSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Table does not exist. Please make sure your table id is correct.", http.StatusInternalServerError)
		return
	}
	// Check if there are exactly two players on the table
	if len(table.Players) != 2 {
		http.Error(w, "Need exactly two players to start the game", http.StatusBadRequest)
		return
	}
	// Start the game (for example, set the status to true)
	table.Status = true
	table.Starts++
	table.Turn = table.Starts % 2
	table.ClearBoard()
	// Start
	err = models.SaveSession(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, "Failed to save table to Redis", http.StatusInternalServerError)
		return
	}
	// Produce message to Kafka topic
	message := "Game has been started"

	message = table.StatusBoard(message) + table.StringBoard()
	err = models.ProduceMessage(h.KAFKAADDRESS, tableID, message, h.KAFKAUSERNAME, h.KAFKAPASSWORD)
	if err != nil {
		http.Error(w, "Failed to produce message to Kafka topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DropPieceHandler handles requests to drop a piece in the Connect 4 game
func (h *Handler) DropPieceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]
	playerName := vars["name"]
	columnStr := vars["column"]
	column, err := strconv.Atoi(columnStr)
	if err != nil {
		http.Error(w, "Invalid column number", http.StatusBadRequest)
		return
	}
	// Retrieve table from Redis
	table, err := models.GetSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Failed to retrieve table from Redis", http.StatusInternalServerError)
		return
	}
	// Find the player
	var currentPlayer *models.Player
	for _, player := range table.Players {
		if player.Name == playerName {
			currentPlayer = player
			break
		}
	}
	if currentPlayer == nil {
		http.Error(w, "Player not found in the table", http.StatusBadRequest)
		return
	}

	turnname := table.GetPlayersTurn()

	if turnname != playerName {
		http.Error(w, fmt.Sprintf("It is %s's turn. Please wait until %s plays their move.", turnname, turnname), http.StatusBadRequest)
		return
	}
	// Get symbol for the current player based on their position in the Players array
	playerIndex := 0
	for i, player := range table.Players {
		if player == currentPlayer {
			playerIndex = i
			break
		}
	}
	playerSymbol := ""
	if playerIndex == 0 {
		playerSymbol = "X"
	} else if playerIndex == 1 {
		playerSymbol = "O"
	} else {
		http.Error(w, "Unexpected player index", http.StatusInternalServerError)
		return
	}

	err = table.DropPiece(column, playerSymbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if table.Turn == 1 {
		table.Turn = 0
	} else {
		table.Turn = 1
	}
	// Produce message to Kafka topic

	message := fmt.Sprintf("Player %s dropped piece in column %d", currentPlayer.Name, column)
	if table.CheckWin(playerSymbol) {
		table.Players[playerIndex].AddWin()
		table.Status = false
	} else if table.IsBoardFull() {
		table.Status = false
	}

	models.SaveSession(h.Context, table, h.Client)
	message = table.StatusBoard(message) + table.StringBoard()
	err = models.ProduceMessage(h.KAFKAADDRESS, tableID, message, h.KAFKAUSERNAME, h.KAFKAPASSWORD)
	if err != nil {
		http.Error(w, "Failed to produce message to Kafka topic", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// LeaveTableHandler handles requests from players who want to leave the Connect 4 table
func (h *Handler) LeaveTableHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID := vars["tableID"]
	playerName := vars["name"]

	// Retrieve the table from Redis
	table, err := models.GetSession(h.Context, tableID, h.Client)
	if err != nil {
		http.Error(w, "Failed to retrieve table from Redis", http.StatusInternalServerError)
		return
	}
	if table.Status {
		http.Error(w, "Game is currently in progress. Please wait until the game is over", http.StatusBadRequest)
		return
	}
	// Find the player in the table
	playerIndex := -1
	for i, player := range table.Players {
		if player.Name == playerName {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		http.Error(w, "Player not found in the table", http.StatusBadRequest)
		return
	}
	// Remove the player from the table
	table.Players = append(table.Players[:playerIndex], table.Players[playerIndex+1:]...)
	// Save the updated table to Redis
	err = models.SaveSession(h.Context, table, h.Client)
	if err != nil {
		http.Error(w, "Failed to save table to Redis", http.StatusInternalServerError)
		return
	}
	// Produce a message to the Kafka topic indicating the player has left the table
	message := fmt.Sprintf("Player %s left the table ", playerName)
	message = table.StatusBoard(message) + table.StringBoard()
	err = models.ProduceMessage(h.KAFKAADDRESS, tableID, message, h.KAFKAUSERNAME, h.KAFKAPASSWORD)
	if err != nil {
		http.Error(w, "Failed to produce message to Kafka topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
