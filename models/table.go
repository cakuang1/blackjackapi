package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Table struct {
	ID              string    `json:"id"`
	Turn            int       `json:"turn"`
	Status          bool      `json:"status"`
	Players         []*Player `json:"players"`
	NumberOfPlayers int       `json:"number_of_players"`
	Dealer          *Player   `json:"dealer"`
}

// TABLE OPERATIONS

func NewTable(id string) *Table {
	return &Table{
		ID:              id,
		Players:         []*Player{},
		Dealer:          NewPlayer("Dealer"),
		Status:          false,
		Turn:            0,
		NumberOfPlayers: 0,
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
