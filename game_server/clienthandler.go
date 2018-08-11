package game_server

import (
	"net"
	"github.com/FedericoPonzi/TicTacToe/gamelogic"
	"encoding/binary"
	"log"
	"encoding/gob"
)

type ClientHandler struct {
	conn net.Conn
	game *gamelogic.TicTacToeGame
}

func NewClientHandler(conn net.Conn)(ClientHandler) {
	return ClientHandler{conn, nil}
}

func (handler *ClientHandler) receiveMove() (toRet int8, err error){
	err = binary.Read(handler.conn, binary.BigEndian, &toRet)

	if err != nil {
		log.Println("Error on reading:", err)
	}


	return toRet, err
}

// Send gamelogic board to the client
func (handler *ClientHandler) sendStatus() {
	encoder := gob.NewEncoder(handler.conn)
	encoder.Encode(handler.game)

}
func (handler *ClientHandler) sendMark(mark int8) error{
	err := binary.Write(handler.conn, binary.BigEndian, mark)
	if err != nil {
		log.Println("Error while sending mark" + string(mark) + ", error: " + err.Error())
	}
	return err
}

func (handler *ClientHandler) sendMove(move int8) error {
	log.Println("Going to send move '" + gamelogic.ConvInt8(move) + "to client")
	err := binary.Write(handler.conn, binary.BigEndian, move)
	if err != nil {
		log.Println("Error while sending mark" + string(move))
	}
	return err
}
func (handler *ClientHandler) sendError(e error) {
	i := int8(-1)
	binary.Write(handler.conn, binary.BigEndian, &i)
}


