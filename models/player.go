package models

import (
	"math/rand"
)

var (
	Ranks = map[string]int{
		"2":  2,
		"3":  3,
		"4":  4,
		"5":  5,
		"6":  6,
		"7":  7,
		"8":  8,
		"9":  9,
		"10": 10,
		"J":  10,
		"Q":  10,
		"K":  10,
		"A":  11, // Assuming Ace initially counts as 11, it can be 1 later if needed
	}
)

// Player represents a player in the Blackjack game
type Player struct {
	Name  string   `json:"name"`
	Hand  []string `json:"hand"`
	Value int      `json:"value"`
}

// NewPlayer creates a new player with the given ID and name
func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Hand:  []string{},
		Value: 0,
	}
}

// AddCard adds a specific card to the player's hand
func (p *Player) AddCard(rank string) {
	p.Hand = append(p.Hand, rank)
}

// AddRandomCard adds a random card to the player's hand
func (p *Player) AddRandomCard() {
	rank := GetRandomCard()
	p.Hand = append(p.Hand, rank)
}



func (p *Player) SetScore() {
	score := 0
	numAces := 0 // Count of Aces in the hand
	for _, rank := range p.Hand {
		score += Ranks[rank]
		if rank == "A" {
			numAces++
		}
	}
	// Adjust score if necessary due to Aces
	for numAces > 0 && score > 21 {
		score -= 10 // Subtract 10 for each Ace until the score is below or equal to 21
		numAces--
	}
	p.Value = score
}

// GetRandomCard returns a random card rank
func GetRandomCard() string {
	// Shuffle the ranks and select a random one
	var ranks []string
	for rank := range Ranks {
		ranks = append(ranks, rank)
	}
	// Select a random rank
	return ranks[rand.Intn(len(ranks))]
}
