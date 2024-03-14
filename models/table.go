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
const cardWidth = 6

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

func (t *Table) GetBoardText() string {
	var builder strings.Builder

	// Write table ID and status
	builder.WriteString(fmt.Sprintf("Table ID: %s\n", t.ID))
	builder.WriteString(fmt.Sprintf("Status: %s\n", getStatusText(t.Status)))

	// Write dealer's hand
	builder.WriteString("Dealer's Hand:\n")
	builder.WriteString(formatHandText(t.Dealer.Hand))

	// Write players' hands
	builder.WriteString("Players' Hands:\n")
	for _, player := range t.Players {
		builder.WriteString(fmt.Sprintf("%s's Hand:\n", player.Name))
		builder.WriteString(formatHandText(player.Hand))
	}

	return builder.String()
}

// Helper function to format a hand into text
func formatHandText(hand []string) string {
	var builder strings.Builder

	for _, card := range hand {
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", cardWidth)))
		builder.WriteString(fmt.Sprintf("| %-2s  |\n", card))
		builder.WriteString(fmt.Sprintf("|     |\n"))
		builder.WriteString(fmt.Sprintf("|  %-2s |\n", card))
		builder.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", cardWidth)))
	}
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
