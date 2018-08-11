package game_server

import (
	"github.com/FedericoPonzi/TicTacToe/gamelogic"
	"log"
	"strconv"
	"reflect"
	"fmt"
	"io"
)

type GameHandler struct{
	first  ClientHandler
	second ClientHandler
	game   gamelogic.TicTacToeGame
}


type InvalidMoveError struct {
	arg  int
	prob string
}

func (e *InvalidMoveError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}


func NewGameHandler (first ClientHandler, second ClientHandler) (GameHandler) {
	game := gamelogic.NewTicTacToeGame()
	return GameHandler{first, second, game}
}
func (handler *GameHandler) RunGame() {
	//1. send board status,
	//2. receive from first,
	//3. receive from second,
	//4. add liveness probe.
	handler.first.sendMark(0) //TODO: handle error
	handler.second.sendMark(1) //TODO: handler error.
	log.Println("Succesfully sent marks! Game is set, let's start!")
	var(
		move int8
		err error
	)

	for !handler.game.IsOver() {

		currentHandler := handler.getCurrentTurnHandler()
		move, err = currentHandler.receiveMove() // todo handle error
		if err == io.EOF {
			// I can do better.
			break;
		}
		log.Println("Move received: " + strconv.Itoa(int((move))))
		if move == -1 { //something went very bad. Adieu
			log.Println("Move is -1, adieu!")
			break
		}

		err = handler.game.DoMove(move)

		log.Println("Trying to register move " + gamelogic.ConvInt8(move) + "...")
		if err == nil {
			// The move is valid, so send it to both clients.
			handler.first.sendMove(move)
			handler.second.sendMove(move)
			log.Println("Move sent to both clients!")
		}else {
			log.Println("Error sent to clients!")
			currentHandler.sendError(err)
		}
	}
	log.Println("Game is over, thanks for playing!")
}


func (handler *GameHandler) sendError(e error) {
	if reflect.TypeOf(e) == reflect.TypeOf(InvalidMoveError{}) {
		handler.first.sendError(e)
	}
}
func (handler *GameHandler) getCurrentTurnHandler() ClientHandler {
	if handler.game.Turn() == 0 {
		return handler.first
	}
	return handler.second
}
