# TicTacToe game server

This is a simple client/server system to play TicTacToe.

The server waits for two opponents/clients to connect, and handle the communication between them.

There are still many TODO (it still need an appropriate error handling) but its seems usable already.

It's made in go and using goroutines it is able to handle multiple game simultaneously.

The protocol is very simple: 
 * when two opponents connect, it sends them their "mark", 
either a `o` or a `x`.

 * `x` moves first. The client sends the move. 
 * If the move is valid, the server will validate it 
by sending to both clients the next move. Both client have a copy of the board, and will update their 
local copy.

 * When the game is over both clients get disconnected.

There are still things to do, for example it needs a way to know if the connection 
is broken.


## Run it
It's made with golang 1.10.3 but it should work with other versions too.
Build both server and clients:
```golang
#Download dependencies first:
go get -d ./...
go build server.go
go build client.go
```
Run first the server:
```
./server
```
And then two clients:
```
./client
#In another terminal
./client
```
And have fun!

# TODO:
 * Add an AI: After a threshold waiting time, an opponent should start a game with an AI.
 * Add an heartbeating: to know if the other opponent is gone.
 * Add a UI: I'm planning to add a simple UI with QT library.
