package main

import (
	"errors"
	"sync"
)

const (
	turnOrder = "order"
	turnChaos = "chaos"
)

var ErrOutOfTurn = errors.New("It is not currently your turn")
var ErrMissingGame = errors.New("You're not currently in a game")
var ErrIllegalPosition = errors.New("Illegal position")

type Game struct {
	rwMutex sync.RWMutex
	board   *Board `json:"board"`
	turn    string `json:"turn"`
	order   *User  `json:"order"`
	chaos   *User  `json:"chaos"`
}

func NewGame(order, chaos *User) *Game {
	return &Game{
		board: NewBoard(),
		turn:  "order",
		order: order,
		chaos: chaos,
	}
}

func (game *Game) Move(user *User, move *Move) error {
	if game == nil {
		return ErrMissingGame
	}
	game.rwMutex.Lock()
	defer game.rwMutex.Unlock()
	if (user == game.order && game.turn != turnOrder) ||
		(user == game.chaos && game.turn != turnChaos) {
		return ErrOutOfTurn
	}
	if !game.board.IsLegalPosition(&move.Position) {
		return ErrIllegalPosition
	}
	return nil
}

func (game *Game) Concede(user *User) error {
	if game == nil {
		return ErrMissingGame
	}
	return nil
}
