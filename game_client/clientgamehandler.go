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
	"io"
)
// Client game handler implements both the protocol and the interconnection between the logic and the game.
// This is because the client library it's easier wrt the server counterpart.
type ClientGameHandler struct {
	conn   net.Conn
	game   gamelogic.TicTacToeGame
	myMark int8
}

const (
	InvalidMoveError = -1
	OpponentDisconnectedError = -2
)

var ErrTimeOut = errors.New("Timeout")
var ErrSocketClosed = errors.New("EOF - socket closed.")

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
	handler.receiverMark()
	fmt.Println("Your mark: " + marks[handler.myMark])
	for !handler.game.IsOver() {
		if handler.isMyTurn() {
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
			} else {
				handler.game.DoMove(moveToInt(move))
			}

		} else {
			time.Sleep(time.Second)
			fmt.Println("Waiting for opponent...")

			var nextMove int8
			if err = handler.binaryRead(&nextMove); err != nil {
				if err == ErrSocketClosed {
					fmt.Println("Disconnected from server. Are you still connected to the internet?")
					return
				} else {
					//Timeout? In every other case, retry.
					continue
				}
			}

			if nextMove == OpponentDisconnectedError {
				fmt.Println("Ops! Your Opponent seems disconnected :( this probably makes you the winner :)")
				return
			}

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

// Send move is a two step process:
// 1. send the actual move.
// 2. receive answer from server. If it is the same as sent value, the move is validaded.
//    else it's an invalid move error.
func (handler *ClientGameHandler) sendMove(move int8) error {
	var i int8

	if err := handler.binaryWrite(move); err != nil {
		if err == ErrTimeOut{
			fmt.Println("Timeout sending the move, please retry.")
			return err
		}
	}

	if err := handler.binaryRead(&i); err != nil {
		return err
	}

	log.Println("Received from server: " + convInt8(i) + " and I've sent:" + convInt8(move))

	if i == move {
		// Validated move.
		return nil
	}else if i == InvalidMoveError {
			return errors.New("Illegal move.")
	}

	return errors.New("Error while sending your move. Received: " + convInt8(i) + ".")
}

// Write a binary number, with resiliency.
func (handler *ClientGameHandler) binaryWrite(source int8) (err error ){
	handler.conn.SetReadDeadline(time.Now().Add(5*time.Second))
	if err := binary.Write(handler.conn, binary.BigEndian, &source); err != nil {
		if opError, isOpError := err.(*net.OpError); isOpError {
			if opError.Timeout() {
				log.Println("Timeout.")
				return ErrTimeOut
			}
		}
		if err == io.EOF {
			return ErrSocketClosed
		}
		// What kind of error is this?
		log.Print("Error with receive move: " + err.Error())
		return err
	}
	return nil
}
// Read a binary number, with resiliency.
func (handler *ClientGameHandler) binaryRead(dest * int8) (err error){
	handler.conn.SetReadDeadline(time.Now().Add(5*time.Second))

	if err := binary.Read(handler.conn, binary.BigEndian, dest); err != nil {
		if opError, isOpError := err.(*net.OpError); isOpError {
			if opError.Timeout() {
				//log.Println("Timeout.")
				return ErrTimeOut
			}
		}
		if err == io.EOF {
			log.Println("EOF.")
			return ErrSocketClosed
		}

		// What kind of error is this?
		log.Println("Error with receive move: " + err.Error())
		return err
	}

	return nil
}

func (handler *ClientGameHandler) receiverMark() {
	waitTime := 24 // 2 mins
	for handler.myMark == -1  && waitTime >= 0{
		handler.binaryRead(&handler.myMark)
		fmt.Print(".")
		waitTime -= 1
	}
	if handler.myMark == -1 {
		fmt.Println("No opponent available at the moment! :(")
		os.Exit(0)
	}
}

func NewClientGameHandler(server string) (ClientGameHandler) {
	toRet := ClientGameHandler{myMark : -1}
	toRet.game = gamelogic.NewTicTacToeGame()
	d := net.Dialer{Timeout: time.Second * 2}
	conn, err := d.Dial("tcp", server + ":10119")

	if err != nil {
		log.Fatal("Error while connecting to " + server )
	}

	toRet.conn = conn
	fmt.Println("Succesfully connected to the server.")
	fmt.Print("Waiting for an opponent...")
	return toRet
}