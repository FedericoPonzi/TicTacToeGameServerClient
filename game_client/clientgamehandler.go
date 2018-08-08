package game_client

import (
	"net"
	"github.com/FedericoPonzi/TicTacToe/gamelogic"
	"strconv"
	"bufio"
	"os"
	"encoding/binary"
	"fmt"
	"log"
	"time"
	"errors"
)

type ClientGameHandler struct {
	conn   net.Conn
	game   gamelogic.TicTacToeGame
	myMark int8
}

func moveToInt(move string) int8 {
	toRet := int8(((move[0] - 97) * 3) + move[1] - '0')
	return toRet
}
func (handler *ClientGameHandler) isMyTurn() (toRet bool){
	return handler.myMark == handler.game.Turn()
}

func convInt8(i int8) string{
	return strconv.Itoa(int(i))
}

func (handler *ClientGameHandler) RunGame() {
	defer handler.conn.Close()
	var move string
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	// Receive mark:
	binary.Read(handler.conn, binary.BigEndian, &handler.myMark)
	fmt.Println("your mark: " + convInt8(handler.myMark))
	for !handler.game.IsOver() {
		if handler.isMyTurn() { // TODO: and it's still connected - maybe using the heartbeating?
			fmt.Println("It's your turn baby :D")
			handler.game.PrintBoard()
			if scanner.Scan() {
				move = scanner.Text()
				log.Println("Red: " + move)
				err = handler.sendMove(moveToInt(move))
				if err != nil {
					fmt.Println("Error with the move: ", err)
				} else {
					handler.game.DoMove(moveToInt(move))
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
		} else {
			time.Sleep(time.Second)
			fmt.Println("Waiting for opponent...")
			var nextMove int8
			binary.Read(handler.conn, binary.BigEndian, &nextMove)
			handler.game.DoMove(nextMove)
		}


	}
	if handler.game.HasWon(handler.myMark) {
		fmt.Println("Yayy! You won! :D")
	}else {
		fmt.Println("Aww :( You loose!")
	}
}

func (handler *ClientGameHandler) sendMove(move int8) error {
	var i int8
	err := binary.Write(handler.conn, binary.BigEndian, &move)

	if err != nil {
		log.Println("Error while econding/sending move: ", err)
	}

	err = binary.Read(handler.conn, binary.BigEndian, &i)
	log.Println("Received from server: " + convInt8(i) + " and I've sent:" + convInt8(move))
	if i == move {
		// Validated move.
		return nil
	}	else {
		if i == -1 {
			return errors.New("Illegal move error.")
		}
	}
	return errors.New("Error while sending your move, received: " + convInt8(i))
}

func NewClientGameHandler(server string) (ClientGameHandler) {
	toRet := ClientGameHandler{myMark : -1}
	toRet.game = gamelogic.NewTicTacToeGame()
	d := net.Dialer{Timeout: time.Second * 2}
	conn, err := d.Dial("tcp", server)

	if err != nil {
		log.Fatal("Error while connecting to " + server )
	}

	toRet.conn = conn
	log.Println("Succesfully connected to the server.")
	return toRet
}