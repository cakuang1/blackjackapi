package models

import "testing"

// TestNewPlayer tests the NewPlayer function
func TestNewPlayer(t *testing.T) {
    name := "testPlayer"
    player := NewPlayer(name)

    if player.Name != name {
        t.Errorf("Expected player name to be %s, got %s", name, player.Name)
    }
    if player.Wins != 0 {
        t.Errorf("Expected player wins to be 0, got %d", player.Wins)
    }
}

// TestAddWin tests the AddWin method
func TestAddWin(t *testing.T) {
    player := NewPlayer("testPlayer")

    player.AddWin()
    if player.Wins != 1 {
        t.Errorf("Expected player wins to be 1 after adding a win, got %d", player.Wins)
    }

    player.AddWin()
    if player.Wins != 2 {
        t.Errorf("Expected player wins to be 2 after adding another win, got %d", player.Wins)
    }
}
