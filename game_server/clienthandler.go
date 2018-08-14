package game_server

import (
	"net"
	"github.com/FedericoPonzi/TicTacToe/gamelogic"
	"encoding/binary"
	"log"
	"encoding/gob"
	"time"
	"errors"
	"io"
)

var ErrTimeOut = errors.New("Timeout")
var ErrSocketClosed = errors.New("EOF - socket closed.")

const (
	InvalidMoveError = -1
	OpponentDisconnectedError = -2
)

type ClientHandler struct {
	conn net.Conn
	game *gamelogic.TicTacToeGame
}

func NewClientHandler(conn net.Conn)(*ClientHandler) {
	return &ClientHandler{conn, nil}
}

func (handler *ClientHandler) binaryRead(dest * int8) (err error){
	handler.conn.SetReadDeadline(time.Now().Add(5*time.Second))

	if err := binary.Read(handler.conn, binary.BigEndian, dest); err != nil {
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
		log.Println("Error with receive move: " + err.Error())
		return err
	}

	return nil
}
func (handler *ClientHandler) binaryWrite(source int8) (err error ){
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
		log.Println("Error with receive move: " + err.Error())
		return err
	}
	return nil
}

func (handler *ClientHandler) receiveMove() (toRet int8, err error){
	if err = handler.binaryRead(&toRet); err != nil {
		return -1, err;
	}
	return toRet, nil
}

// Send gamelogic board to the client
func (handler *ClientHandler) sendStatus() {
	encoder := gob.NewEncoder(handler.conn)
	encoder.Encode(handler.game)

}
// Send Mark to the client - either 0 or 1.
func (handler *ClientHandler) sendMark(mark int8) error{
	return handler.binaryWrite(mark)
}

func (handler *ClientHandler) sendMove(move int8) error {
	log.Println("Going to send move '" + gamelogic.ConvInt8(move) + "' to client")
	err := binary.Write(handler.conn, binary.BigEndian, move)
	if err != nil {
		log.Println("Error while sending mark" + string(move))
	}
	return err
}
func (handler *ClientHandler) sendError(err int8) error{
	return handler.binaryWrite(err)
}
// Check if the socket is still open.
func (handler *ClientHandler) IsConnectionOpen() bool{

	one := []byte{}
	handler.conn.SetReadDeadline(time.Now())
	if _, err := handler.conn.Read(one); err == io.EOF {
		return false
	}
	return true
}


