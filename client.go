package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    *Room
}

func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    nil,
	}

}
func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			log.Println("connection over")
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (client *Client) disconnect() {
	if client.rooms != nil {
		client.rooms.unregister <- client
	}
	client.wsServer.unregister <- client
	close(client.send)
	client.rooms = nil
	client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	fmt.Println("get a ws con")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message MessageFromUser
	//parse user msg
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	switch message.Action {
	case JoinRoom:

		client.handleJoinRoom(message)
	case LeaveRoom:
		client.handleLeaveRoom(message)
	case PlayPawn:
		client.handlePlayPawn(message)

	}
}

func (client *Client) handleJoinRoom(msg MessageFromUser) {
	roomName := msg.RoomName
	log.Println("join room:", roomName)
	client.joinRoom(roomName)

}
func (client *Client) handleLeaveRoom(msg MessageFromUser) {
	roomName := msg.RoomName
	log.Println("leave room in:", roomName)
	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		return
	}
	room.unregister <- client

	var amsg MessageFromServer
	amsg.RoomName = ""
	if client == client.rooms.Client1 {
		amsg.RoomName = ""
		amsg.Player = 0
		amsg.Player1Online = false
		amsg.Player2Online = false
		amsg.Turn = 0
	}

	client.send <- amsg.encode()
	client.rooms = nil
}
func (client *Client) joinRoom(roomName string) {
	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}
	//client is not in any room
	if client.rooms == nil {
		//check if there are some space for new player
		if !room.Player1Online || !room.Player2Online {
			client.rooms = room
		}
		room.register <- client
	}
}
func (client *Client) handlePlayPawn(msg MessageFromUser) {
	roomName := msg.RoomName
	log.Println("get a pawn", msg.X, msg.Y, roomName)
	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		log.Println("play a pawn on nil room")
	}
	p := Pawn{
		X:      msg.X,
		Y:      msg.Y,
		client: client,
	}
	log.Println("send pawn to room", p)
	room.pawn <- p

}
