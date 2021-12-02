# the implementation of GOMOKU

## system structure

### frontend:

### backend:

backend is written in golang using multi-threads and channels to

## c/s communication

we use websocket protocol to build steady channel between clients and server.

### message from user

1. There are 3 different actions that a client can send.

```go
const PlayPawn = "play-pawn"
const JoinRoom = "join-room"
const LeaveRoom = "leave-room"

type MessageFromUser struct {
	Action   string `json:"action"`
	RoomName string `json:"message"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
}
```

2. Besides actions above, user can also disconnect the websocket connection without notifying server.

### message from server

user can parse the json file in websocket connection to sync with other players
this structure bellow maintains the metadata of a Room
the 10-by-10 board is encoded in a single row 100 bytes array

```go
type MessageFromServer struct {
	RoomName      string    `json:"roomName"`
	Player        int32     `json:"player"`
	Player1Online bool      `json:"player1Online"`
	Player2Online bool      `json:"player2Online"`
	Turn          int32     `json:"turn"`
	Board         [100]byte `json:"board"`
}
```

## how to support multiple players concurrently
