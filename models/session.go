package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
)

type Session struct {
	ID            string     `json:"id"`
	Turn          int        `json:"turn"`
	Status        bool       `json:"status"`
	Players       []*Player  `json:"players"`
	Grid          [][]string `json:"grid"` // Representing the Connect Four grid
	Starts        int
	OccupiedSlots int `json:"occupied_slots"` // Counter for the number of occupied slots
}

const (
	EmptySlot     = " "
	Player1Symbol = "X" // Symbol for player 1
	Player2Symbol = "O" // Symbol for player 2
)

func (s *Session) StatusBoard(announcement string) string {
	var sb strings.Builder

	// Calculate the length of the longest string in each section
	var maxLengthID, maxLengthStatus, maxLengthPlayers int
	for _, player := range s.Players {
		if len(player.Name) > maxLengthPlayers {
			maxLengthPlayers = len(player.Name)
		}
	}
	maxLengthID = len(s.ID)
	maxLengthStatus = len("Game is in Progress")
	if !s.Status {
		maxLengthStatus = len("Game has not started")
	}
	// Define the fixed width of the box
	boxWidth := 60 // Adjust this value as needed
	// Top line of the box
	sb.WriteString("┌")
	sb.WriteString(strings.Repeat("─", boxWidth-2))
	sb.WriteString("┐\n")

	// Session ID
	sessionid := fmt.Sprintf("│ Session ID: %-"+fmt.Sprintf("%d", maxLengthID)+"s ", s.ID)
	sb.WriteString(sessionid)
	sb.WriteString(strings.Repeat(" ", boxWidth-len(sessionid)+1))
	sb.WriteString("│\n")

	// Game status
	status := "Game is in Progress"
	if !s.Status {
		status = "Game has not started"
	}

	stat := fmt.Sprintf("│ Status: %-"+fmt.Sprintf("%d", maxLengthStatus)+"s ", status)
	sb.WriteString(stat)
	sb.WriteString(strings.Repeat(" ", boxWidth-len(stat)+1))
	sb.WriteString("│\n")

	// Player list
	playerList := strings.Join(func() []string {
		var names []string
		for _, v := range s.Players {
			names = append(names, v.Name)
		}
		return names
	}(), ", ")
	players := fmt.Sprintf("│ Player List: %-"+fmt.Sprintf("%d", maxLengthPlayers)+"s ", playerList)
	sb.WriteString(players)
	sb.WriteString(strings.Repeat(" ", boxWidth-len(players)+1))
	sb.WriteString("│\n")

	// Current turn
	if s.Status {
		currentplayer := s.Players[s.Turn].Name
		turn := fmt.Sprintf("│ Current Turn: %s", currentplayer)
		sb.WriteString(turn)
		sb.WriteString(strings.Repeat(" ", boxWidth-len(turn)+1))
		sb.WriteString("│\n")
	}

	// Players and their symbols
	for i, player := range s.Players {
		symbol := ""
		if i == 0 {
			symbol = Player1Symbol // Assign symbol "X" to the first player
		} else if i == 1 {
			symbol = Player2Symbol // Assign symbol "O" to the second player
		}

		playerInfo := fmt.Sprintf("| Player: %s - Symbol: %s - Wins: %d", player.Name, symbol, player.Wins)
		sb.WriteString(playerInfo)
		sb.WriteString(strings.Repeat(" ", boxWidth-len(playerInfo)-1))
		sb.WriteString("│\n")
	}

	// Middle line of the box

	// Announcement
	sb.WriteString(fmt.Sprintf("│ %-"+fmt.Sprintf("%d", boxWidth-4)+"s │\n", announcement))

	// Bottom line of the box
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", boxWidth-2))
	sb.WriteString("┘\n")

	return sb.String()
}

// Initialize the Connect Four grid
func NewSession(id string) *Session {
	width := 7  // Typical width for Connect Four
	height := 6 // Typical height for Connect Four
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = EmptySlot
		}
	}
	return &Session{
		ID:            id,
		Players:       []*Player{},
		Status:        false,
		Turn:          0,
		Grid:          grid,
		Starts:        0,
		OccupiedSlots: 0,
	}
}

func (s *Session) ClearBoard() {
	s.OccupiedSlots = 0
	for r := 0; r <= len(s.Grid)-1; r++ {
		for c := 0; c <= len(s.Grid[0])-1; c++ {
			s.Grid[r][c] = EmptySlot
		}
	}
}

func (s *Session) AddPlayer(player *Player) error {
	s.Players = append(s.Players, player)
	return nil
}

// DropPiece updates the Connect Four grid with the player's piece and increments the occupied slots counter
func (s *Session) DropPiece(column int, playerSymbol string) error {
	// Check if the column is out of range
	if column < 0 || column >= len(s.Grid[0]) {
		return errors.New("column index out of range")
	}
	for i := len(s.Grid) - 1; i >= 0; i-- {
		if s.Grid[i][column] == EmptySlot {
			s.Grid[i][column] = playerSymbol
			s.OccupiedSlots++ // Increment the counter
			return nil
		}
	}

	// If all slots in the column are occupied
	return errors.New("no more slots available in the column")
}

// IsBoardFull checks if the Connect Four board is completely filled
func (s *Session) IsBoardFull() bool {
	return s.OccupiedSlots == len(s.Grid)*len(s.Grid[0]) // Check if all slots are occupied
}

// CheckWin checks if a player has won the game
func (s *Session) CheckWin(playerSymbol string) bool {
	// Check horizontally
	for _, row := range s.Grid {
		for j := 0; j <= len(row)-4; j++ {
			if strings.Join(row[j:j+4], "") == strings.Repeat(playerSymbol, 4) {
				return true
			}
		}
	}

	// Check vertically
	for i := 0; i <= len(s.Grid)-4; i++ {
		for j := 0; j < len(s.Grid[i]); j++ {
			if s.Grid[i][j] == playerSymbol && s.Grid[i+1][j] == playerSymbol && s.Grid[i+2][j] == playerSymbol && s.Grid[i+3][j] == playerSymbol {
				return true
			}
		}
	}

	// Check diagonally (from bottom-left to top-right)
	for i := 3; i < len(s.Grid); i++ {
		for j := 0; j <= len(s.Grid[i])-4; j++ {
			if s.Grid[i][j] == playerSymbol && s.Grid[i-1][j+1] == playerSymbol && s.Grid[i-2][j+2] == playerSymbol && s.Grid[i-3][j+3] == playerSymbol {
				return true
			}
		}
	}
	// Check diagonally (from top-left to bottom-right)
	for i := 0; i <= len(s.Grid)-4; i++ {
		for j := 0; j <= len(s.Grid[i])-4; j++ {
			if s.Grid[i][j] == playerSymbol && s.Grid[i+1][j+1] == playerSymbol && s.Grid[i+2][j+2] == playerSymbol && s.Grid[i+3][j+3] == playerSymbol {
				return true
			}
		}
	}
	return false
}

// StringBoard returns a string representation of the Connect Four board
func (s *Session) StringBoard() string {
	var sb strings.Builder
	// Determine the number of rows and columns in the grid
	numCols := len(s.Grid[0]) // Assuming all rows have the same number of columns
	liveBoardString := "Live board"
	padding := (numCols*4 + (numCols - 1)) / 2                            // Calculate padding to center the string
	sb.WriteString(strings.Repeat(" ", padding-len(liveBoardString)/2-3)) // Add left padding
	sb.WriteString(liveBoardString)
	sb.WriteString("\n\n")
	// Print grid
	for _, row := range s.Grid {
		sb.WriteString("|")
		for _, slot := range row {
			sb.WriteString(fmt.Sprintf(" %s |", slot)) // Ensure each slot is separated by lines
		}
		sb.WriteString("\n")
		sb.WriteString("-")
		sb.WriteString(strings.Repeat("----", numCols)) // Add line separator between rows
		sb.WriteString("\n")
	}

	// Print column numbers centered
	sb.WriteString(" ")
	for col := 0; col < numCols; col++ { // Start from 0 to shift one position to the left
		sb.WriteString(fmt.Sprintf("%2d  ", col+1)) // Add 1 to col to start from 1 instead of 0
	}
	sb.WriteString("\n")

	return sb.String()
}

func (s *Session) GetPlayersTurn() string {
	if len(s.Players) == 0 {
		return "No players in the session"
	}

	currentPlayerIndex := s.Turn % len(s.Players)
	currentPlayer := s.Players[currentPlayerIndex]
	return currentPlayer.Name
}

// GRAB AND SAVE SESSIONS

func GetSession(ctx context.Context, id string, client *redis.Client) (*Session, error) {
	val, err := client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize the JSON-encoded value into a Session struct
	var session Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func SaveSession(ctx context.Context, session *Session, client *redis.Client) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	// Save the JSON-encoded data to Redis
	if err := client.Set(ctx, session.ID, data, 0).Err(); err != nil {
		return err
	}

	return nil
}

// DeleteSession deletes a Connect 4 session from Redis.
func DeleteSession(ctx context.Context, sessionID string, client *redis.Client) error {
	// Delete the session from Redis
	if err := client.Del(ctx, sessionID).Err(); err != nil {
		return err
	}
	return nil
}
