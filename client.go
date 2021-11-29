package main

import (
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	rooms    map[*Room]bool
}
