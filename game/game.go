package game

import (
	"sync"
	"time"
)

type Player struct {
	Username string
	Score    int64
	Drawing  string
}

type GameState struct {
	players   []*Player
	round     int64
	maxRounds int64
	roundTime time.Time
	IsEnded   bool
	lock      sync.Mutex
}

func New(rounds int64, roundTime time.Time) *GameState {
	return &GameState{
		maxRounds: rounds,
		roundTime: roundTime,
	}
}

func (gs *GameState) StartRound() {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	gs.round++
}

func (gs *GameState) AddPlayer(username string) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	for _, p := range gs.players {
		if p.Username == username {
			return
		}
	}

	gs.players = append(gs.players, &Player{
		Username: username,
	})
}
