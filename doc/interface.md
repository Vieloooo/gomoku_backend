# interface doc

interface doc introduce the baisc msg transfered between user client and serve

## actions under connnection

user websocket protocol for c/s communication

1. init upgrader (upgrade the http connection to websocket connection)
2. regular pingpong (use websocket default pingpong handler)
3. user action
   1. join room
   2. leave room
   3. play the pawn
4. server feedback
   1. return meta_info
5. unexpected situation:
   1. a client close connection(reset the )

## msg json design

### userMsg

```go
type MessageFromUser struct {
	Action   string  `json:"action"`
	RoomName string  `json:"message"`
	Target   *Room   `json:"target"`
	Sender   *Client `json:"sender"`
	X        int32   `json:"x"`
	Y        int32   `json:"y"`
}
```

1. actions:
   1. init ws connection
   2. join room (Roomname,)
   3. play the pawn(msg:"x,y")

### serverMsg

```go
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
```
