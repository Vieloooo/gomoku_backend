package main

import "log"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	delRoom    chan *Room
	rooms      map[*Room]bool
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		delRoom:    make(chan *Room),
		rooms:      make(map[*Room]bool),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)
		case r := <-server.delRoom:
			server.deleteRoom(r)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
	//send *Client back to user
	var msg MessageFromServer
	msg.RoomName = ""
	log.Println("register client ")
	client.send <- msg.encode()
}

func (server *WsServer) unregisterClient(client *Client) {
	delete(server.clients, client)
}

func (server *WsServer) deleteRoom(room *Room) {
	log.Println("del room", room.Name)
	delete(server.rooms, room)
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.Name == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string) *Room {
	room := NewRoom(name, server)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}
