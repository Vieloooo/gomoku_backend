package main

import (
	"encoding/json"
	"log"
)

const PlayPawn = "play-pawn"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"

type MessageFromUser struct {
	Action   string  `json:"action"`
	RoomName string  `json:"message"`
	Target   *Room   `json:"target"`
	Sender   *Client `json:"sender"`
	X        int32   `json:"x"`
	Y        int32   `json:"y"`
}
type MessageFromServer struct {
	RoomName      string    `json:"roomName"` //if roomname == null, no in room
	Target        *Room     `json:"target"`   //put in MessageFromUser.
	Sender        *Client   `json:"sender"`
	Player        int32     `json:"player"`
	Player1Online bool      `json:"player1Online"`
	Player2Online bool      `json:"player2Online"`
	Turn          int32     `json:"turn"` //0 not start, 1 player1's turn, 2 play2's turn
	Board         [100]byte `json:"board"`
}

func (message *MessageFromServer) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
