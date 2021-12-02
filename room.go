package main

import (
	"log"
	"time"
)

type Pawn struct {
	X      int32
	Y      int32
	client *Client
}

const boardcastTime = 5000 * time.Millisecond

type Room struct {
	Name          string `json:"name"`
	wsServer      *WsServer
	Client1       *Client
	Client2       *Client
	Player1Online bool `json:"player1Online"`
	Player2Online bool `json:"player2Online"`
	register      chan *Client
	unregister    chan *Client
	over          chan int
	pawn          chan Pawn
	Turn          int32     `json:"turn"`
	Board         [100]byte `json:"board"`
}

// NewRoom creates a new Room
func NewRoom(name string, ws *WsServer) *Room {
	var b [100]byte
	for i := 0; i < 100; i++ {
		b[i] = '0'
	}
	return &Room{
		Name:          name,
		wsServer:      ws,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		pawn:          make(chan Pawn),
		over:          make(chan int),
		Turn:          0,
		Player1Online: false,
		Player2Online: false,
		Board:         b,
	}

}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	ticker := time.NewTicker(boardcastTime)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
		case p := <-room.pawn:
			room.putPawn(p.X, p.Y, p.client)
		case <-room.over:
			return
		case <-ticker.C:
			room.updateInfo()
		}

	}
}
func (room *Room) updateInfo() {

	msg := MessageFromServer{
		RoomName:      room.Name,
		Turn:          room.Turn,
		Player1Online: room.Player1Online,
		Player2Online: room.Player2Online,
		Board:         room.Board,
	}
	if room.Player1Online {
		msg.Player = 1
		room.Client1.send <- msg.encode()
	}
	if room.Player2Online {
		msg.Player = 2
		room.Client2.send <- msg.encode()
	}
}
func (room *Room) checkStart() {
	if room.Player1Online && room.Player2Online {
		room.Turn = 1
	} else {
		room.Turn = 0
	}
}
func (room *Room) checkEnd() {
	log.Println("check end, player1", room.Player1Online, "player2", room.Player2Online)
	if !room.Player1Online && !room.Player2Online {
		room.wsServer.delRoom <- room
	}
}
func (room *Room) registerClientInRoom(client *Client) {
	defer room.updateInfo()
	if !room.Player1Online {
		room.Player1Online = true
		room.Client1 = client
		room.checkStart()
		return
	}
	if !room.Player2Online {
		room.Player2Online = true
		room.Client2 = client
		room.checkStart()
		return
	}
	//notify full
	var b [100]byte
	for i := 0; i < 100; i++ {
		b[i] = '0'
	}

	aamsg := MessageFromServer{
		RoomName:      "error",
		Turn:          0,
		Player1Online: false,
		Player2Online: false,
		Board:         b,
	}
	client.send <- aamsg.encode()

}
func (room *Room) resetBoard() {
	for i := 0; i < 100; i++ {
		room.Board[i] = '0'
	}
}
func (room *Room) unregisterClientInRoom(client *Client) {
	defer room.updateInfo()
	if client == room.Client1 {
		room.Turn = 0
		room.Player1Online = false
		room.resetBoard()
		room.checkEnd()
		return
	}
	if client == room.Client2 {
		room.Turn = 0
		room.Player2Online = false
		room.resetBoard()
		room.checkEnd()
		return
	}
}

func (room *Room) putPawn(x, y int32, c *Client) {
	log.Println("update pawn of room:", room.Name, x, y)
	defer room.updateInfo()
	if room.Player1Online && c == room.Client1 && room.Turn == 1 {
		room.Turn = 2
		if room.Board[10*x+y] == '0' {
			room.Board[10*x+y] = '1'
		}

		return
	}
	if room.Player2Online && c == room.Client2 && room.Turn == 2 {
		room.Turn = 1
		if room.Board[10*x+y] == '0' {
			room.Board[10*x+y] = '2'
		}
		return
	}
}
