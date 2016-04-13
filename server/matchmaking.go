package main

import (
	"sync"
)

var matchmaking Matchmaking

type Matchmaking struct {
	rwMutex sync.RWMutex
	Waiting *User
}

func Find(user *User) {

}
