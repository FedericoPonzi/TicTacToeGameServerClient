package game_server

import (
	"github.com/FedericoPonzi/TicTacToe/gamelogic"
	"log"
)

type GameHandler struct{
	first  *ClientHandler
	second *ClientHandler
	game   gamelogic.TicTacToeGame
}

func NewGameHandler (first *ClientHandler, second *ClientHandler) (GameHandler) {
	game := gamelogic.NewTicTacToeGame()

	return GameHandler{first, second, game}
}


func (handler *GameHandler) RunGame() {
	defer handler.first.conn.Close()
	defer handler.second.conn.Close()
	if err := handler.first.sendMark(0); err != nil {
		if err == ErrSocketClosed {
			log.Println("Client disconnected before sending mark")
		}
		handler.second.sendError(OpponentDisconnectedError)
	}

	if err := handler.second.sendMark(1); err != nil {
		if err == ErrSocketClosed {
			log.Println("Client disconnected before sending mark")
		}
		handler.second.sendError(OpponentDisconnectedError)
	}

	//1. send board status,
	//2. receive from first,
	//3. receive from second,
	//4. add liveness probe.
	log.Println("Succesfully sent marks! Game is set, let's start!")
	var(
		move int8
		err error
	)

	connLoop:for ! handler.game.IsOver() {

		currentHandler := handler.getCurrentTurnHandler()
		if move, err = currentHandler.receiveMove(); err != nil {
			switch err {
				case ErrSocketClosed:
					log.Println("Socket closed.")
					handler.getNextTurnHandler().sendError(OpponentDisconnectedError) //Ignore timeout/error
					break connLoop
				case ErrTimeOut:
					// Everybody need some time to think TicTacToe is a serious game.
					continue
				default:
					log.Println("Something very strange happend.")
					continue
			}
		}
		log.Println("Trying to register move " + gamelogic.ConvInt8(move) + "...")
		if err = handler.game.DoMove(move); err != nil {
			log.Println("Received an invalid move.")
			currentHandler.sendError(InvalidMoveError)
		} else {
			handler.first.sendMove(move)
			handler.second.sendMove(move)
			log.Println("Move sent to both clients!")
		}
	}

	log.Println("Game is over, Thanks for playing!")
}


func (handler *GameHandler) getCurrentTurnHandler() *ClientHandler {
	if handler.game.Turn() == 0 {
		return handler.first
	}
	return handler.second
}
func (handler *GameHandler) getNextTurnHandler() *ClientHandler {
	if handler.game.Turn() == 0 {
		return handler.second
	}
	return handler.first
}
