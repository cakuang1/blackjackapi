package models

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

type Table struct {
	ID      string
	Players []*Player
}

// TABLE OPERATIONS

func NewTable(id string) *Table {
	return &Table{
		ID:      id,
		Players: []*Player{},
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

// GRAB AND SAVE TABLES

func GetTable(ctx context.Context, id string, client *redis.Client) (*Table, error) {
	val, err := client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	// Deserialize the JSON-encoded value into a slice of *Player structs
	var players []*Player
	if err := json.Unmarshal([]byte(val), &players); err != nil {
		return nil, err
	}
	// Create a new Table instance with the retrieved data
	return &Table{
		ID:      id,
		Players: players,
	}, nil
}

func SaveTable(ctx context.Context, table *Table, client *redis.Client) error {
	data, err := json.Marshal(table.Players)
	if err != nil {
		return err
	}
	// Save the JSON-encoded data to Redis
	if err := client.Set(ctx, table.ID, data, 0).Err(); err != nil {
		return err
	}
	return nil
}
