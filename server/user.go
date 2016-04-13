package main

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var ErrInsideRoom = errors.New("You must leave your current room to perform that action")
var ErrOutsideRoom = errors.New("You must be in a room to perform that action")

type User struct {
	rwMutex  sync.RWMutex
	ws       *websocket.Conn
	outgoing chan *Message
	room     *Room
}

func NewUser(ws *websocket.Conn) *User {
	user := &User{
		ws:       ws,
		outgoing: make(chan *Message),
	}
	go user.listen()
	return user
}

func (user *User) listen() {
	go user.read()
	go user.write()
}

// read reads messages from the user's websocket
func (user *User) read() {
	for {
		var message Message
		err := user.ws.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				break
			} else {
				user.SendMessage(NewErrorMessage(err))
				continue
			}
		}
		err = user.handleMessage(&message)
		if err != nil {
			user.SendMessage(NewErrorMessage(err))
			continue
		}
	}
	// user disconnected
	user.LeaveRoom()
}

func (user *User) handleMessage(message *Message) error {
	switch message.Type {
	case MsgTyEcho:
		user.SendMessage(message)
	case MsgTyCreateRoom:
		user.CreateRoom()
	case MsgTyJoinRoom:
		user.JoinRoom(message.Value.(string))
	case MsgTyLeaveRoom:
		user.LeaveRoom()
	default:
		user.SendMessage(NewErrorMessage(ErrInavlidMsgTy))
	}
	return nil
}

// write writes messages to the user's websocket from their outgoing channel
func (user *User) write() {
	for message := range user.outgoing {
		err := user.ws.WriteJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				break
			} else {
				log.Println("unexpected error:", err)
				user.ws.Close()
				break
			}
		}
	}
	// user disconnected
}

func (user *User) SendMessage(message *Message) {
	user.outgoing <- message
}

func (user *User) CreateRoom() {
	user.rwMutex.Lock()
	defer user.rwMutex.Unlock()
	if user.room != nil {
		user.SendMessage(NewErrorMessage(ErrInsideRoom))
		return
	}
	room := rooms.NewRoom(user)
	user.room = room
	// TODO
	user.SendMessage(&Message{"joinRoom", room.id})
}

func (user *User) JoinRoom(id string) {
	user.rwMutex.Lock()
	defer user.rwMutex.Unlock()
	if user.room != nil {
		user.SendMessage(NewErrorMessage(ErrInsideRoom))
		return
	}
	room, err := rooms.GetRoom(id)
	if err != nil {
		user.SendMessage(NewErrorMessage(err))
		return
	}
	err = room.AddUser(user)
	if err != nil {
		user.SendMessage(NewErrorMessage(err))
		return
	}
	user.room = room
	// TODO
	user.SendMessage(&Message{"joinRoom", room.id})
}

func (user *User) LeaveRoom() {
	user.rwMutex.Lock()
	defer user.rwMutex.Unlock()
	if user.room == nil {
		user.SendMessage(NewErrorMessage(ErrOutsideRoom))
		return
	}
	err := user.room.RemoveUser(user)
	if err != nil {
		user.SendMessage(NewErrorMessage(err))
		return
	}
	user.room = nil
	user.SendMessage(&Message{"leaveRoom", nil})
}
