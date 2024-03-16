package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strings"
)

type Table struct {
	ID      string    `json:"id"`
	Turn    int       `json:"turn"`
	Status  bool      `json:"status"`
	Players []*Player `json:"players"`
	Dealer  *Player   `json:"dealer"`
}

// TABLE OPERATIONS

func NewTable(id string) *Table {
	return &Table{
		ID:      id,
		Players: []*Player{},
		Dealer:  NewPlayer("Dealer"),
		Status:  false,
		Turn:    0,
	}
}

func (t *Table) AddPlayer(player *Player) error {
	// Check if the table is full
	if len(t.Players) >= 5 {
		return errors.New("Table is full")
	}
	t.Players = append(t.Players, player)
	return nil
}

func (t *Table) DeletePlayer(name string) error {
	// Find the index of the player with the given name
	index := -1
	for i, player := range t.Players {
		if player.Name == name {
			index = i
			break
		}
	}

	// If the player is not found, return an error
	if index == -1 {
		return errors.New("Player not found")
	}
	// If the player is the last element, remove it without appending
	fmt.Println(t.Players)

	t.Players = remove(t.Players, index)
	fmt.Println(t.Players)

	return nil
}

func (t *Table) PlayerIsin(name string) bool {
	listofplayer := t.Players
	for _, v := range listofplayer {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (t *Table) TableClear() {
	for _, v := range t.Players {
		v.Hand = []string{}
		v.Value = 0
	}
	t.Turn = 0
	t.Dealer.Hand = []string{}
	t.Dealer.Value = 0
}

func (t *Table) StartTable() {
	t.Dealer.Hand = append(t.Dealer.Hand, GetRandomCard())
	t.Dealer.Hand = append(t.Dealer.Hand, GetRandomCard())
	t.Dealer.Value = CalculateHandTotal(t.Dealer.Hand)
	for _, v := range t.Players {
		v.Hand = append(v.Hand, GetRandomCard())
		v.Hand = append(v.Hand, GetRandomCard())
		v.Value = CalculateHandTotal(v.Hand)
	}
	t.Status = true
}

// Helper function to format a hand into text
func (table *Table) GetBoardText() string {
	var builder strings.Builder
	maxWidth := 50
	// Calculate the maximum hand length among players and dealer
	maxHandLength := len(table.Dealer.Hand)
	for _, player := range table.Players {
		if len(player.Hand) > maxHandLength {
			maxHandLength = len(player.Hand)
		}
	}
	// Calculate the total width based on the longest hand length
	totalWidth := maxWidth + maxHandLength*5 + 15
	// Print the board header
	builder.WriteString(strings.Repeat("-", totalWidth) + "\n")
	builder.WriteString("|" + centerText("BLACKJACK", totalWidth-2) + "|\n")
	builder.WriteString(strings.Repeat("-", totalWidth) + "\n")
	tableid := fmt.Sprintf("| Table ID:%s ", table.ID)
	builder.WriteString(tableid + strings.Repeat(" ", totalWidth-len(tableid)-1) + "|\n")
	status := fmt.Sprintf("| Status:%s ", GetStatus(table.Status))
	builder.WriteString(status + strings.Repeat(" ", totalWidth-len(status)-1) + "|\n")
	turn := fmt.Sprintf("| Turn:%d ", table.Turn)
	builder.WriteString(turn + strings.Repeat(" ", totalWidth-len(turn)-1) + "|\n")
	playercount := fmt.Sprintf("| Number of Players:%d", len(table.Players))
	builder.WriteString(playercount + strings.Repeat(" ", totalWidth-len(playercount)-1) + "|\n")
	// need to first convert
	var listofplayers []string
	for _, v := range table.Players {
		listofplayers = append(listofplayers, v.Name)
	}
	playerlist := fmt.Sprintf("| Players in lobby : %s", listofplayers)
	builder.WriteString(playerlist + strings.Repeat(" ", totalWidth-len(playerlist)-1) + "|\n")
	// Print dealer's hand
	if table.Status {
		// what are the conditions here
		deal := ""
		for i := 1; i <= len(table.Dealer.Hand)-1; i++ {
			deal += fmt.Sprintf("[ %s ]", table.Dealer.Hand[i])
		}
		dealerhand1 := "| Dealer:[ ? ] " + deal
		builder.WriteString(dealerhand1 + strings.Repeat(" ", totalWidth-len(dealerhand1)-1) + "|\n")
		for _, player := range table.Players {
			playerstring := ""
			for _, v := range player.Hand {
				playerstring += fmt.Sprintf("[ %s ]", v)
			}
			new := fmt.Sprintf("| %s: %s", player.Name, playerstring)
			builder.WriteString(new + strings.Repeat(" ", totalWidth-len(new)-1) + "|\n")
		}
	}

	builder.WriteString(strings.Repeat("-", totalWidth) + "\n")
	return builder.String()
}

// GRAB AND SAVE TABLES

func GetStatus(status bool) string {
	if status {
		return "Game is in progress. Player deletion "
	}
	return "Game has not started. You are allowed to join "
}

func GetTable(ctx context.Context, id string, client *redis.Client) (*Table, error) {
	val, err := client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize the JSON-encoded value into a Table struct
	var table Table
	if err := json.Unmarshal([]byte(val), &table); err != nil {
		return nil, err
	}

	return &table, nil
}

func SaveTable(ctx context.Context, table *Table, client *redis.Client) error {
	data, err := json.Marshal(table)
	if err != nil {
		return err
	}
	// Save the JSON-encoded data to Redis
	if err := client.Set(ctx, table.ID, data, 0).Err(); err != nil {
		return err
	}

	return nil
}

func CalculateHandTotal(hand []string) int {
	total := 0
	aces := 0
	for _, card := range hand {
		switch card {
		case "J", "Q", "K":
			total += 10
		case "A":
			aces++
		default:
			total += CardValue(card)
		}
	}
	for aces > 0 && total+11 <= 21 {
		total += 11
		aces--
	}
	for aces > 0 {
		total++
		aces--
	}
	return total
}

func CardValue(card string) int {
	switch card {
	case "2", "3", "4", "5", "6", "7", "8", "9", "10":
		return int(card[0] - '0')
	case "J", "Q", "K":
		return 10
	case "A":
		return 11
	default:
		return 0
	}
}

func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

func RandomRank() string {
	// Define the possible ranks
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	// Generate a random index for the rank
	rankIndex := rand.Intn(len(ranks))

	// Return the randomly selected rank
	return ranks[rankIndex]
}

func remove(slice []*Player, i int) []*Player {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
