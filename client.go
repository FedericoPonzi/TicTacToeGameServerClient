package main

import (
	"flag"
	"github.com/FedericoPonzi/TicTacToe/game_client"
)


func main() {

	var server string

	flag.StringVar(&server, "server", "127.0.0.1:10100", "Address of the server in the form of ip:port.")
	flag.Parse()

	client := game_client.NewClientGameHandler(server)
	client.RunGame()
}
