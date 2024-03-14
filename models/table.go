package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
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

	// Print dealer's hand
	var dealerHand string
	if table.Status {
		dealerHand = fmt.Sprintf("[ ? ] [ %s ]%s", table.Dealer.Hand[1], strings.Repeat(" ", totalWidth-50))
	} else {
		dealerHand = " "
	}

	builder.WriteString(fmt.Sprintf("| Dealer:   %s|\n", dealerHand))
	// Print players' hands
	for _, player := range table.Players {
		playerHand := fmt.Sprintf("[ %s ] ", strings.Join(player.Hand, "] [ "))
		builder.WriteString(fmt.Sprintf("| Player %s: %s", player.Name, playerHand))
		totalStr := fmt.Sprintf("(Total: %d)", CalculateHandTotal(player.Hand))
		builder.WriteString(fmt.Sprintf("%s%s|\n", totalStr, strings.Repeat(" ", totalWidth-len(totalStr)-len(player.Hand)*5-15)))
	}
	// Print bottom of the board
	builder.WriteString(strings.Repeat("-", totalWidth) + "\n")
	return builder.String()
}

// GRAB AND SAVE TABLES

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

func getStatusText(status bool) string {
	if status {
		return "In Play"
	}
	return "Not Started"
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

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}
