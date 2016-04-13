package main

import "errors"

const (
	MsgTyError       = "error"
	MsgTyEcho        = "echo"
	MsgTyCreateRoom  = "createRoom"
	MsgTyJoinRoom    = "joinRoom"
	MsgTyLeaveRoom   = "leaveRoom"
	MsgTyChatMessage = "chatMessage"
)

var ErrInavlidMsgTy = errors.New("Invalid message type")

type Message struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func NewErrorMessage(err error) *Message {
	return &Message{
		Type:  MsgTyError,
		Value: err.Error(),
	}
}
