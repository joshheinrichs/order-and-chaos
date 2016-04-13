package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"sync"
)

var ErrMissingRoom = errors.New("A room with that ID does not exist")
var ErrDuplicateRoom = errors.New("A room with that ID already exists")

var rooms = NewRooms()

type Rooms struct {
	rwMutex sync.RWMutex
	rooms   map[string]*Room
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: make(map[string]*Room),
	}
}

// NewID generates a random string ID. The length of the string roughly equal to
// the given length but due to rounding errors is not guaranteed to be correct
// for all lengths.
func NewID(length int) string {
	b := make([]byte, length*6/8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func (rooms *Rooms) NewRoom(user *User) *Room {
	rooms.rwMutex.Lock()
	defer rooms.rwMutex.Unlock()
	// Generates IDs until it finds one that doesn't conflict. This is
	// potentially scary, but in practice conflicts should basically never occur
	var id string
	for {
		id = NewID(6)
		_, exists := rooms.rooms[id]
		if !exists {
			break
		}
	}
	room := &Room{
		id:      id,
		users:   make(map[*User]interface{}),
		queue:   make([]*User, 0),
		game:    nil,
		deleted: false,
	}
	room.users[user] = nil
	rooms.rooms[id] = room
	log.Printf("created a new room with id %s\n", id)
	return room
}

func (rooms *Rooms) GetRoom(id string) (*Room, error) {
	rooms.rwMutex.RLock()
	defer rooms.rwMutex.RUnlock()
	room, exists := rooms.rooms[id]
	if !exists {
		return nil, ErrMissingRoom
	}
	return room, nil
}

func (rooms *Rooms) DeleteRoom(id string) error {
	rooms.rwMutex.Lock()
	defer rooms.rwMutex.Unlock()
	_, exists := rooms.rooms[id]
	if !exists {
		return ErrMissingRoom
	}
	delete(rooms.rooms, id)
	log.Printf("deleted room with id %s\n", id)
	return nil
}

type Room struct {
	rwMutex sync.RWMutex
	// id is immutable and should not be changed after creation
	id      string
	users   map[*User]interface{}
	queue   []*User
	game    *Game
	deleted bool
}

// AddUser adds a user to the room, returning an error if unsuccessful.
func (room *Room) AddUser(user *User) error {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	if room.deleted {
		return ErrMissingRoom
	}
	room.users[user] = nil
	return nil
}

// RemoveUser removes a user from the room, return an error if unsuccessful.
func (room *Room) RemoveUser(user *User) error {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	delete(room.users, user)
	if len(room.users) == 0 {
		room.deleted = true
		rooms.DeleteRoom(room.id)
	}
	// TODO: Remove from queue and game
	return nil
}

// func (room *Room) BroadcastMessage(user *User, chatMessage *ChatMessage) error {
// 	return nil
// }
