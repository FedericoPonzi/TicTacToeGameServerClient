package main

import (
	"flag"
	"github.com/FedericoPonzi/TicTacToe/game_client"
	"regexp"
	"fmt"
)
func validatedIp(ip string) bool {
	r := regexp.MustCompile(`^\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b$`)
	return r.MatchString(ip)
}

func main() {

	var server string
	flag.StringVar(&server, "server", "127.0.0.1", "Address of the server in the form of ip:port.")
	flag.Parse()
	if validatedIp(server) {
		client := game_client.NewClientGameHandler(server)
		client.RunGame()
	} else {
		fmt.Println("It seems like you have entered a wrong ip address for the server: " + server)
	}
}
