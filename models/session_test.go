package models

import (
	"testing"
)

// TestNewSession tests the NewSession function
func TestNewSession(t *testing.T) {
	id := "testSession"
	session := NewSession(id)

	if session.ID != id {
		t.Errorf("Expected session ID to be %s, got %s", id, session.ID)
	}
	if len(session.Players) != 0 {
		t.Errorf("Expected no players in the session, got %d", len(session.Players))
	}
	if session.Status != false {
		t.Errorf("Expected session status to be false, got %v", session.Status)
	}
	if session.Turn != 0 {
		t.Errorf("Expected session turn to be 0, got %d", session.Turn)
	}
	if len(session.Grid) != 6 {
		t.Errorf("Expected grid to have 6 rows, got %d", len(session.Grid))
	}
	if len(session.Grid[0]) != 7 {
		t.Errorf("Expected grid to have 7 columns, got %d", len(session.Grid[0]))
	}
}

// TestAddPlayer tests the AddPlayer method
func TestAddPlayer(t *testing.T) {
	session := NewSession("testSession")
	player := &Player{Name: "Player1"}

	err := session.AddPlayer(player)

	if err != nil {
		t.Errorf("Error adding player: %v", err)
	}
	if len(session.Players) != 1 {
		t.Errorf("Expected 1 player in the session, got %d", len(session.Players))
	}
	if session.Players[0] != player {
		t.Errorf("Expected player to be added to the session")
	}
}

// TestDropPiece tests the DropPiece method
func TestDropPiece(t *testing.T) {
	session := NewSession("testSession")
	session.DropPiece(0, Player1Symbol)

	if session.Grid[5][0] != Player1Symbol {
		t.Errorf("Expected player's piece to be dropped at the bottom of column 0")
	}
}

// TestCheckWin tests the CheckWin method
func TestCheckWin(t *testing.T) {
	session := NewSession("testSession")
	// Set up a winning condition horizontally
	session.Grid[5][0] = Player1Symbol
	session.Grid[5][1] = Player1Symbol
	session.Grid[5][2] = Player1Symbol
	session.Grid[5][3] = Player1Symbol

	if !session.CheckWin(Player1Symbol) {
		t.Errorf("Expected a win for player 1")
	}
}
