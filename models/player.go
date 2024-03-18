package models

type Player struct {
	Name string `json:"name"`
	Wins int    `json:"wins"`
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Wins: 0,
	}
}

func (p *Player) AddWin() {
	p.Wins++
}
