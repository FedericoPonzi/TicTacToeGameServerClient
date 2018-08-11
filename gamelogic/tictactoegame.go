package gamelogic

import (
	"errors"
	"fmt"
	"strconv"
	"log"
)

func ConvInt8(i int8) string{
	return strconv.Itoa(int(i))
}


type TicTacToeGame struct {
	board [9] int8
	turn int8
}

func (game *TicTacToeGame) Turn () int8 {
	return game.turn
}
func (game *TicTacToeGame) Board() [9] int8{
	return game.board
}

func (game *TicTacToeGame) IsOver () (bool){
	toRet := true
	for _, i := range game.board {
		toRet = toRet && i != -1
	}
	return game.HasWon(0) || game.HasWon(1) || toRet
}
// check if mark has won. mark can be 0 or 1.
func (game *TicTacToeGame) HasWon(mark int8) (toRet bool) {
	toRet = false

	if game.board[4] == mark {
		toRet = game.board[0] == mark && game.board[8] == mark
		toRet = toRet || (game.board[1] == mark && game.board[7] == mark )
		toRet = toRet || (game.board[2] == mark && game.board[6] == mark )
		toRet = toRet || (game.board[3] == mark && game.board[5] == mark )
	}
	if game.board[0] == mark {
		toRet = toRet || (game.board[3] == mark && game.board[6] == mark)
		toRet = toRet || (game.board[1] == mark && game.board[2] == mark)
	}
	if game.board[8] == mark {
		toRet = toRet || (game.board[7] == mark && game.board[6] == mark)
		toRet = toRet || (game.board[5] == mark && game.board[2] == mark)

	}
	return toRet

}
func (game *TicTacToeGame) DoMove(move int8) error {

	if move > 9 || move < 0 || game.board[move] != -1 {
		log.Println("Illegal move.")
		return errors.New("Invalid move.")
	}

	game.board[move] = game.turn
	game.turn = (game.turn + 1 ) % 2

	return nil
}
func (game *TicTacToeGame) PrintBoard() {
	var mark string
	for i, el := range game.board {
		if el == 0 {
			mark = " x "
		} else if el == 1 {
			mark = " o "
		} else {
			mark = " - "
		}
		if i % 3 == 2 {
			fmt.Println(mark)
			if i != 8 {
				fmt.Println("-----------")
			}
		} else{
			fmt.Print(mark + "|")
		}
	}
}

func NewTicTacToeGame() (TicTacToeGame) {
	return TicTacToeGame {turn : 0 , board : [9]int8 { -1, -1, -1, -1,-1,-1,-1,-1,- 1,}}
}
