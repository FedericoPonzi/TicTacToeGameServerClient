package main

import (
	"net"
	"log"
	"github.com/FedericoPonzi/TicTacToe/game_server"
)

//Wait for two opponents, and start the gamelogic
func opponentSeeker(opponents chan game_server.ClientHandler){

	for {
		log.Println("Looking for opponents...")
		firstOpponent := <- opponents
		log.Println("First opponent in.")
		secondOpponent := <- opponents
		log.Println("Second opponent in, starting gamelogic")

		gameHandler := game_server.NewGameHandler(firstOpponent, secondOpponent)
		go gameHandler.RunGame()
	}
}


func main() {
	service := ":10100"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatal("Fatal error: %s", err.Error())
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("Fatal error: %s", err.Error())
	}

	opponents := make(chan game_server.ClientHandler)

	go opponentSeeker(opponents)

	for {
		//Accept connections
		conn, err := listener.Accept()
		log.Println("Received connection.")
		if err != nil {
			log.Println("Error on accept: ", err)
			continue
		}

		client := game_server.NewClientHandler(conn)
		opponents <- client

	}
}
