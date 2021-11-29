package main

import (
	"errors"

	"github.com/google/uuid"
)

type Room struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Client1       *Client
	Client2       *Client
	Player1Online bool `json:"player1Online"`
	Player2Online bool `json:"player2Online"`
	register      chan *Client
	unregister    chan *Client
	broadcast     chan *MessageFromServer
	Turn          int32     `json:"turn"`
	Board         [100]byte `json:"board"`
}

// NewRoom creates a new Room
func NewRoom(name string) *Room {
	var b [100]byte
	for i := 0; i < 100; i++ {
		b[i] = '0'
	}
	return &Room{
		ID:            uuid.New(),
		Name:          name,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan *MessageFromServer),
		Turn:          0,
		Player1Online: false,
		Player2Online: false,
		Board:         b,
	}

}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
		}

	}
}
func (room *Room) registerClientInRoom(client *Client) error {
	if !room.Player1Online {
		room.Player1Online = true
		room.Client1 = client
		return nil
	}
	if !room.Player2Online {
		room.Player2Online = true
		room.Client2 = client
		return nil
	}
	return errors.New("no space for new player.")

}

func (room *Room) unregisterClientInRoom(client *Client) {
	if client == room.Client1 {
		room.Player1Online = false
		return
	}
	if client == room.Client2 {
		room.Player2Online = false
		return
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	if room.Player1Online {
		room.Client1.send <- message
	}
	if room.Player2Online {
		room.Client2.send <- message
	}

}

func (room *Room) putPawn(x, y, role int) {
	if role == 1 {
		room.Board[10*x+y] = '1'
		return
	}
	if role == 2 {
		room.Board[10*x+y] = '2'
		return
	}
}
