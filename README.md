# gomoku_backend

a gomoku game backend written in golang

## connections

> you can test on : https://www.xhcode.com/websocket/

1. init connections

```
ws://localhost:8080/ws
```

2. join room

```
{
    "action":"join-room",
    "message":"firstroom"
}
```

3. leave room

```
{
    "action":"leave-room",
    "message":"firstroom"
}
```

4. play the pawn

```
{
    "action":"play-pawn",
    "message":"firstroom",
    "X":0,
    "Y":4
}
```

## build instructure

1. go installed
2. go get gorilla/websocket
3. go build ./
4. now server will host on localhost:8080
