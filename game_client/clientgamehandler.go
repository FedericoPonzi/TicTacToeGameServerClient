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
// Client game handler implements both the protocol and the interconnection between the logic and the game.
// This is because the client library it's easier wrt the server counterpart.
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

func (handler *ClientGameHandler) readMove(scanner *bufio.Scanner) (string, error) {
	var move string
	var err error
	if scanner.Scan() {
		move = scanner.Text()
		if len(move) > 2 {
			return "", errors.New("Please provide a valid move, in the form: [a-c][0-2]")
		}
		if move[0] < 97 || move[0] > 99 {
			return "", errors.New("There are only three rows: a,b,c.")
		}
		if move[1] - '0' > 2 {
			return "", errors.New("There are only three columns: 0,1,2.")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return move, err

}

func (handler *ClientGameHandler) RunGame() {
	defer handler.conn.Close()
	var move string
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	marks := []string {"x", "o"}

	// Receive mark:
	binary.Read(handler.conn, binary.BigEndian, &handler.myMark)
	fmt.Println("Your mark: " + marks[handler.myMark])
	for !handler.game.IsOver() {
		if handler.isMyTurn() { // TODO: and it's still connected - maybe using the heartbeating?
			fmt.Println("It's your turn baby :D")
			handler.game.PrintBoard()
			move, err = handler.readMove(scanner)
			if err != nil {
				fmt.Println("Errore: " + err.Error())
				continue
			}
			err = handler.sendMove(moveToInt(move))

			if err != nil {
				fmt.Println("Error with the move: ", err)
			}
			if err == nil {
				handler.game.DoMove(moveToInt(move))
			}
		} else {
			time.Sleep(time.Second)
			fmt.Println("Waiting for opponent...")
			var nextMove int8
			binary.Read(handler.conn, binary.BigEndian, &nextMove)
			handler.game.DoMove(nextMove)
		}


	}
	handler.game.PrintBoard()
	if handler.game.HasWon(handler.myMark) {
		fmt.Println("Yayy! You won! :D")
	}else if handler.game.HasWon((handler.myMark + 1) % 2 ) {
		fmt.Println("Aww :( You lost!")
	} else {
		fmt.Println("It's a tie!")
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