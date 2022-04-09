/*
chess.go
By: Matthew Yungbluth
CSc 372
Final Project
Due: December 7th, 2020 at the beginning of class

This program contains the logic and display of the actual chess g.boardState and game
*/

package main

import (
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

var whiteKingCastle, whiteQueenCastle, blackKingCastle, blackQueenCastle, playerInCheck bool

/*
Game ... Controls the current state of the game
*/
type Game struct {
	//g.boardStateState is the current state of the chess g.boardState in a 2d slice
	boardState [8][8]int
	//playerColor: 0 is white at bottom, 1 is black at bottom
	playerColor int
	//playerTurn: 0 means white turn, 1 means black turn
	playerTurn int
	//gameOver is 2 if game still running, is 0 if white won and 1 if black won
	gameOver int
}

//returns a bool if the player is in a check mate
func isInMatePlayer(g *Game, board [8][8]int) bool {
	if g.playerColor == 0 {
		//now see if any white has any pieces that can move
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 0 && board[i][j] < 10 {
					for a := 0; a < 8; a++ {
						for b := 0; b < 8; b++ {
							if canMoveHere(g, j, i, b, a, board, true) {
								if !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, a, b)) {
									return false
								}
							}
						}
					}
				}
			}
		}
		g.gameOver = 1
		return true
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] > 10 {
				for a := 0; a < 8; a++ {
					for b := 0; b < 8; b++ {
						if canMoveHere(g, j, i, b, a, board, true) {
							if !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, a, b)) {
								return false
							}
						}
					}
				}
			}
		}
	}
	g.gameOver = 0
	return true
}

//returns a bool if the AI is in a check mate
func isInMateAi(g *Game, board [8][8]int) bool {
	if g.playerColor == 1 {
		//now see if any white has any pieces that can move
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 0 && board[i][j] < 10 {
					for a := 0; a < 8; a++ {
						for b := 0; b < 8; b++ {
							if canMoveHere(g, j, i, b, a, board, true) {
								if !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, a, b)) {
									return false
								}
							}
						}
					}
				}
			}
		}
		g.gameOver = 1
		return true
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] > 10 {
				for a := 0; a < 8; a++ {
					for b := 0; b < 8; b++ {
						if canMoveHere(g, j, i, b, a, board, true) {
							if !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, a, b)) {
								return false
							}
						}
					}
				}
			}
		}
	}
	g.gameOver = 0
	return true
}

//returns a bool if the player's king is in check
func isInCheckPlayer(g *Game, board [8][8]int) bool {
	if g.playerColor == 0 {
		//check if a black piece can move to white king
		//first find white king
		kingI, kingJ := -1, -1
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] == 6 || board[i][j] == 106 {
					kingI, kingJ = i, j
					break
				}
			}
			if kingI != -1 {
				//break out if king has been found
				break
			}
		}
		//now see if any black piece can move to white king's square
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 10 && board[i][j] < 100 {
					if canMoveHere(g, j, i, kingJ, kingI, board, true) {
						return true
					}
				}
			}
		}
	} else {
		//check if a white piece can move to black king
		//first find black king
		kingI, kingJ := -1, -1
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] == 16 || board[i][j] == 116 {
					kingI, kingJ = i, j
					break
				}
			}
			if kingI != -1 {
				//break out if king has been found
				break
			}
		}
		//now see if any white piece can move to black king's square
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] < 10 && board[i][j] > 0 {
					if canMoveHere(g, j, i, kingJ, kingI, board, true) {
						return true
					}
				}
			}
		}
	}
	return false
}

//returns a bool if the ai's king is in check
func isInCheckAi(g *Game, board [8][8]int) bool {
	if g.playerColor == 1 {
		//check if a black piece can move to white king
		//first find white king
		kingI, kingJ := -1, -1
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] == 6 || board[i][j] == 106 {
					kingI, kingJ = i, j
					break
				}
			}
			if kingI != -1 {
				//break out if king has been found
				break
			}
		}
		//now see if any black piece can move to white king's square
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 10 && board[i][j] < 100 {
					if canMoveHere(g, j, i, kingJ, kingI, board, true) {
						return true
					}
				}
			}
		}
	} else {
		//check if a white piece can move to black king
		//first find black king
		kingI, kingJ := -1, -1
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] == 16 || board[i][j] == 116 {
					kingI, kingJ = i, j
					break
				}
			}
			if kingI != -1 {
				//break out if king has been found
				break
			}
		}
		//now see if any white piece can move to black king's square
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] < 10 && board[i][j] > 0 {
					if canMoveHere(g, j, i, kingJ, kingI, board, true) {
						return true
					}
				}
			}
		}
	}
	return false
}

//moves a piece at [i1][j1] to [i2][j2]
func createTempBoardMove(board [8][8]int, i1, j1, i2, j2 int) [8][8]int {
	newBoard := board
	newBoard[i2][j2] = board[i1][j1]
	newBoard[i1][j1] = 0
	return newBoard
}

//Checks if piece at coordinates (x1,y1) can move to coordinates (x2,y2)
func canMoveHere(g *Game, x1, y1, x2, y2 int, board [8][8]int, check bool) bool {
	//playerColor: 0 is white at bottom, 1 is black at bottom

	//white: 1 is pawn, 2 is bishop, 3 is knight, 4 is rook, 5 is queen, 6 is king
	//black: 11 is pawn, 12 is bishop, 13 is knight, 14 is rook, 15 is queen, 16 is king
	piece := board[y1][x1]
	if piece > 100 {
		piece -= 100
	}
	if g.playerTurn == 0 && piece > 10 && !check {
		//black player trying to move on white turn
		return false
	}
	if g.playerTurn == 1 && piece < 10 && piece > 0 && !check {
		//white player tring to move on black turn
		return false
	}
	switch piece {
	case 1:
		//white pawn
		if g.playerColor == 0 {
			if y1-y2 == 1 && x1 == x2 && board[y2][x2] == 0 {
				//can move forward to an empty square
				return true
			}
			if y1-y2 == 2 && x1 == x2 && board[y2][x2] == 0 && board[y2+1][x2] == 0 && y1 == 6 {
				//can move twice forward
				return true
			}
			if y1-y2 == 1 && (x2-x1 == 1 || x1-x2 == 1) && board[y2][x2] > 10 && (board[y2][x2] < 100 || board[y2][x2] > 110) {
				//can take diagonally
				return true
			}
		} else {
			if y2-y1 == 1 && x1 == x2 && board[y2][x2] == 0 {
				//can move forward to an empty square
				return true
			}
			if y2-y1 == 2 && x1 == x2 && board[y2][x2] == 0 && board[y2-1][x2] == 0 && y1 == 1 {
				//can move twice forward
				return true
			}
			if y2-y1 == 1 && (x2-x1 == 1 || x1-x2 == 1) && board[y2][x2] > 10 && (board[y2][x2] < 100 || board[y2][x2] > 110) {
				//can take diagonally
				return true
			}
		}
		return false
	case 11:
		//black pawn
		//very similar to white pawn but since can only move forward direction is changed, making small changes in the programming forcing me to create a new case
		if g.playerColor == 1 {
			if y1-y2 == 1 && x1 == x2 && board[y2][x2] == 0 {
				//can move forward to an empty square
				return true
			}
			if y1-y2 == 2 && x1 == x2 && board[y2][x2] == 0 && board[y2+1][x2] == 0 && y1 == 6 {
				//can move twice forward
				return true
			}
			if y1-y2 == 1 && (x2-x1 == 1 || x1-x2 == 1) && (board[y2][x2] < 10 || board[y2][x2] > 100) && board[y2][x2] > 0 {
				//can take diagonally
				return true
			}
		} else {
			if y2-y1 == 1 && x1 == x2 && board[y2][x2] == 0 {
				//can move forward to an empty square
				return true
			}
			if y2-y1 == 2 && x1 == x2 && board[y2][x2] == 0 && board[y2-1][x2] == 0 && y1 == 1 {
				//can move twice forward
				return true
			}
			if y2-y1 == 1 && (x2-x1 == 1 || x1-x2 == 1) && (board[y2][x2] < 10 || board[y2][x2] > 100) && board[y2][x2] > 0 {
				//can take diagonally
				return true
			}
		}
	case 2, 12:
		//bishops
		//check between the click and the position diagonally, there must be nothing in between and must not land on same color, doesn't matter if it lands on enemy or empty square
		if math.Abs(float64(x1-x2)) == math.Abs(float64(y1-y2)) && x1 != x2 {
			difference := math.Abs(float64(x1 - x2))
			//click is indeed diagonal, now check for obstacles in the way
			for offset := 1; offset < int(difference); offset++ {
				if x1 < x2 {
					if y1 < y2 {
						//moving to the bottom right
						if board[y1+offset][x1+offset] != 0 {
							return false
						}
					} else {
						//moving to the top right
						if board[y1-offset][x1+offset] != 0 {
							return false
						}
					}
				} else {
					if y1 < y2 {
						//moving to the bottom left
						if board[y1+offset][x1-offset] != 0 {
							return false
						}
					} else {
						//moving to the top left
						if board[y1-offset][x1-offset] != 0 {
							return false
						}
					}
				}
			}
			//finally check to make sure the final position isn't a friendly piece
			if piece == 2 {
				//white bishop
				if board[y2][x2] < 10 && board[y2][x2] > 0 {
					return false
				}
			} else {
				//black bishop
				if board[y2][x2] > 10 && board[y2][x2] < 100 {
					return false
				}
			}
			return true
		}
	case 3, 13:
		//knights
		//either difference of x is 2 and y is 1, or difference of x is 1 and y is 2
		xDif := math.Abs(float64(x1 - x2))
		yDif := math.Abs(float64(y1 - y2))
		if (xDif == 2 && yDif == 1) || (xDif == 1 && yDif == 2) {
			//now just check to make sure final position isn't a friendly piece
			if piece == 3 {
				//white knight
				if board[y2][x2] < 10 && board[y2][x2] > 0 {
					return false
				}
			} else {
				//black knight
				if board[y2][x2] > 10 && board[y2][x2] < 100 {
					return false
				}
			}
			return true
		}
	case 4, 14:
		//rooks
		//either difference of x is 0 or difference of y is 0
		xDif := math.Abs(float64(x1 - x2))
		yDif := math.Abs(float64(y1 - y2))
		if xDif == 0 && yDif == 0 {
			//not moving
			return false
		}
		if xDif > 0 && yDif > 0 {
			//not moving straight
			return false
		}
		if xDif == 0 || yDif == 0 {
			//now check that there are no obstacles in between the piece and destination
			if xDif == 0 {
				//moving up or down
				if y1-y2 < 0 {
					//moving down
					for i := 1; i < int(yDif); i++ {
						if board[y1+i][x1] != 0 {
							return false
						}
					}
				} else {
					//moving up
					for i := 1; i < int(yDif); i++ {
						if board[y1-i][x1] != 0 {
							return false
						}
					}
				}
			} else if yDif == 0 {
				//moving left or right
				if x1-x2 < 0 {
					//moving right
					for i := 1; i < int(xDif); i++ {
						if board[y1][x1+i] != 0 {
							return false
						}
					}
				} else {
					//moving left
					for i := 1; i < int(xDif); i++ {
						if board[y1][x1-i] != 0 {
							return false
						}
					}
				}
			}
		}
		//lastly check that the destination isn't a friendly piece
		if piece == 4 {
			//white rook
			if board[y2][x2] < 10 && board[y2][x2] > 0 {
				return false
			}
		} else {
			//black rook
			if board[y2][x2] > 10 && board[y2][x2] < 100 {
				return false
			}
		}
		return true
	case 5, 15:
		//queens
		//first check if moving diagonally or in a straight line
		xDif := math.Abs(float64(x1 - x2))
		yDif := math.Abs(float64(y1 - y2))
		if xDif == 0 && yDif == 0 {
			//didn't move
			return false
		}
		if xDif == 0 || yDif == 0 {
			//moving straight line
			//now basically copy rook checks
			if xDif == 0 {
				//moving up or down
				if y1-y2 < 0 {
					//moving down
					for i := 1; i < int(yDif); i++ {
						if board[y1+i][x1] != 0 {
							return false
						}
					}
				} else {
					//moving up
					for i := 1; i < int(yDif); i++ {
						if board[y1-i][x1] != 0 {
							return false
						}
					}
				}
			} else if yDif == 0 {
				//moving left or right
				if x1-x2 < 0 {
					//moving right
					for i := 1; i < int(xDif); i++ {
						if board[y1][x1+i] != 0 {
							return false
						}
					}
				} else {
					//moving left
					for i := 1; i < int(xDif); i++ {
						if board[y1][x1-i] != 0 {
							return false
						}
					}
				}
			}
		} else if xDif == yDif {
			//moving diagonally
			//now basically copy bishop checks
			difference := math.Abs(float64(x1 - x2))
			//click is indeed diagonal, now check for obstacles in the way
			for offset := 1; offset < int(difference); offset++ {
				if x1 < x2 {
					if y1 < y2 {
						//moving to the bottom right
						if board[y1+offset][x1+offset] != 0 {
							return false
						}
					} else {
						//moving to the top right
						if board[y1-offset][x1+offset] != 0 {
							return false
						}
					}
				} else {
					if y1 < y2 {
						//moving to the bottom left
						if board[y1+offset][x1-offset] != 0 {
							return false
						}
					} else {
						//moving to the top left
						if board[y1-offset][x1-offset] != 0 {
							return false
						}
					}
				}
			}
		} else {
			//moving neither straight or diagonal, bad move
			return false
		}
		//finally check if the destination isn't a friendly piece
		if piece == 5 {
			//white queen
			if board[y2][x2] < 10 && board[y2][x2] > 0 {
				return false
			}
		} else {
			//black queen
			if board[y2][x2] > 10 && board[y2][x2] < 100 {
				return false
			}
		}
		return true
	case 6, 16:
		//kings
		xDif := math.Abs(float64(x1 - x2))
		yDif := math.Abs(float64(y1 - y2))
		//king can only move 1 space, first check to make sure he isn't trying to sprint away
		if yDif > 1 {
			return false
		}
		if xDif > 1 {
			//check for castle
			//must be true in order to castle:
			//there must be no pieces between king and rook
			//king must not currently be in check
			//king must not pass through a square that is under attack
			//king must not end up in check
			if g.playerColor == 0 {
				if x2-x1 == 2 {
					if whiteKingCastle {
						//white king side castle
						//first check no pieces between king and rook
						if board[y1][x1+1] != 0 || board[y1][x1+2] != 0 {
							//piece in between them
							return false
						}
						//next check king not in check
						if isInCheckPlayer(g, g.boardState) {
							return false
						}
						//next check no square kings passes is under attack
						//create a temp board state with king there, check to make sure he isn't in check
						newBoardOne := createTempBoardMove(g.boardState, y1, x1, y1, x1+1)
						if isInCheckPlayer(g, newBoardOne) {
							return false
						}
						newBoardTwo := createTempBoardMove(g.boardState, y1, x1, y1, x1+2)
						if isInCheckPlayer(g, newBoardTwo) {
							return false
						}
						return true
					}
				} else if x2-x1 == -2 {
					if whiteQueenCastle {
						//white queen side castle
						if board[y1][x1-1] != 0 || board[y1][x1-2] != 0 || board[y1][x1-3] != 0 {
							//piece in between them
							return false
						}
						//next check king not in check
						if isInCheckPlayer(g, g.boardState) {
							return false
						}
						//next check no square kings passes is under attack
						//create a temp board state with king there, check to make sure he isn't in check
						newBoardOne := createTempBoardMove(g.boardState, y1, x1, y1, x1-1)
						if isInCheckPlayer(g, newBoardOne) {
							return false
						}
						newBoardTwo := createTempBoardMove(g.boardState, y1, x1, y1, x1-2)
						if isInCheckPlayer(g, newBoardTwo) {
							return false
						}
						return true
					}
				}
			} else {
				if x2-x1 == -2 {
					if blackKingCastle {
						//black king side castle
						//first check no pieces between king and rook
						if board[y1][x1-1] != 0 || board[y1][x1-2] != 0 {
							//piece in between them
							return false
						}
						//next check king not in check
						if isInCheckPlayer(g, g.boardState) {
							return false
						}
						//next check no square kings passes is under attack
						//create a temp board state with king there, check to make sure he isn't in check
						newBoardOne := createTempBoardMove(g.boardState, y1, x1, y1, x1-1)
						if isInCheckPlayer(g, newBoardOne) {
							return false
						}
						newBoardTwo := createTempBoardMove(g.boardState, y1, x1, y1, x1-2)
						if isInCheckPlayer(g, newBoardTwo) {
							return false
						}
						return true
					}
				} else if x2-x1 == 2 {
					if blackQueenCastle {
						//black queen side castle
						if board[y1][x1+1] != 0 || board[y1][x1+2] != 0 || board[y1][x1+3] != 0 {
							//piece in between them
							return false
						}
						//next check king not in check
						if isInCheckPlayer(g, g.boardState) {
							return false
						}
						//next check no square kings passes is under attack
						//create a temp board state with king there, check to make sure he isn't in check
						newBoardOne := createTempBoardMove(g.boardState, y1, x1, y1, x1+1)
						if isInCheckPlayer(g, newBoardOne) {
							return false
						}
						newBoardTwo := createTempBoardMove(g.boardState, y1, x1, y1, x1+2)
						if isInCheckPlayer(g, newBoardTwo) {
							return false
						}
						return true
					}
				}
			}
			return false
		}
		//next check to make sure hes moving
		if xDif == 0 && yDif == 0 {
			return false
		}
		//lastly check to make sure destination isn't a friendly piece
		if piece == 6 {
			//white king
			if board[y2][x2] < 10 && board[y2][x2] > 0 {
				return false
			}
		} else {
			//black king
			if board[y2][x2] > 10 && board[y2][x2] < 100 {
				return false
			}
		}
		return true
	}
	return false
}

/*
Update ... updates the game board
*/
func (g *Game) Update() error {
	//update the logical state
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		flag := true
		cursX, cursY := ebiten.CursorPosition()
		x := cursX / 125
		y := cursY / 125
		if g.gameOver == 2 {
			for i := 0; i < 8; i++ {
				for j := 0; j < 8; j++ {
					if g.boardState[i][j] > 100 {
						if x <= 7 && y <= 7 {
							if canMoveHere(g, j, i, x, y, g.boardState, false) {
								if !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, y, x)) {
									//if the previously selected piece can move to this new spot, move the piece
									if g.boardState[y][x] == 6 || g.boardState[y][x] == 16 {
										//if the piece is capturing a king, end the game
										if g.gameOver == 2 {
											g.gameOver = g.playerTurn
										}
									}
									g.boardState[y][x] = g.boardState[i][j] - 100
									g.boardState[i][j] = 0
									if g.boardState[y][x] == 6 {
										if x-j == 2 {
											//king side castle success
											if g.playerColor == 0 {
												g.boardState[7][7] = 0
												g.boardState[7][5] = 4
											} else {
												g.boardState[0][0] = 0
												g.boardState[0][2] = 4
											}
										}
										if j-x == 2 {
											//queen side castle success
											if g.playerColor == 0 {
												g.boardState[7][0] = 0
												g.boardState[7][3] = 4
											} else {
												g.boardState[0][7] = 0
												g.boardState[0][4] = 4
											}
										}
										whiteKingCastle = false
										whiteQueenCastle = false
									}
									if g.boardState[y][x] == 16 {
										if x-j == 2 {
											//king side castle success
											if g.playerColor == 0 {
												g.boardState[0][7] = 0
												g.boardState[0][5] = 14
											} else {
												g.boardState[7][7] = 0
												g.boardState[7][4] = 14
											}
										}
										if j-x == 2 {
											//queen side castle success
											if g.playerColor == 0 {
												g.boardState[0][0] = 0
												g.boardState[0][2] = 14
											} else {
												g.boardState[7][0] = 0
												g.boardState[7][2] = 14
											}
										}
										blackKingCastle = false
										blackQueenCastle = false
									}
									if i == 7 {
										if g.playerColor == 0 {
											//white = king to the right, black = king to the left
											if g.boardState[y][x] == 4 {
												//rook has been moved, check if king side or queen side and disable castling option
												if j == 0 {
													whiteQueenCastle = false
												}
												if j == 7 {
													whiteKingCastle = false
												}
											}
										} else {
											if g.boardState[y][x] == 14 {
												if j == 0 {
													blackKingCastle = false
												}
												if j == 7 {
													blackQueenCastle = false
												}
											}
										}
									}
									if (y == 0 || y == 7) && (g.boardState[y][x] == 1 || g.boardState[y][x] == 11) {
										//pawn hitting the end of the g.boardState, promote to queen
										g.boardState[y][x] += 4
									}
									//change player turn
									if g.playerTurn == 0 {
										g.playerTurn = 1
									} else {
										g.playerTurn = 0
									}
									go isInMateAi(g, g.boardState)
									go aiMove(g)
									//go isInMatePlayer(g, g.boardState)
									flag = false
								} else {
									g.boardState[i][j] -= 100
								}
							} else {
								g.boardState[i][j] -= 100
							}
						}
					}
				}
			}
			if flag && g.boardState[y][x] != 0 {
				g.boardState[y][x] += 100
			}
		}
	}
	return nil
}

//Just matches if piece is same as game turn so only displays boxes when piece selected when its their turn
func canDisplayBox(g *Game, piece int) bool {
	if g.playerTurn == 0 && piece < 10 {
		return true
	}
	if g.playerTurn == 1 && piece > 10 {
		return true
	}
	return false
}

//appens image and imageoptions to slice and returns that slice
func appendImageToSlice(markers []*ebiten.Image, markerop []*ebiten.DrawImageOptions, i, j, joffset, ioffset int, jadd, iadd float64) ([]*ebiten.Image, []*ebiten.DrawImageOptions) {
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	sqr := ebiten.NewImage(24, 24)
	sqr.Fill(red)
	sqrOP := &ebiten.DrawImageOptions{}
	sqrOP.GeoM.Translate(float64(125*(j+joffset))+jadd, float64(125*(i+ioffset))+iadd)
	markers = append(markers, sqr)
	markerop = append(markerop, sqrOP)
	return markers, markerop
}

/*
Draw ... Draws the chess board
*/
func (g *Game) Draw(screen *ebiten.Image) {
	go isInMatePlayer(g, g.boardState)
	//Render the screen
	var boardImage, blackPawn, whitePawn, blackBishop, whiteBishop, blackKing, whiteKing, blackKnight, whiteKnight, blackQueen, whiteQueen, blackRook, whiteRook *ebiten.Image
	//draw the image of the g.boardState
	boardImage, _, _ = ebitenutil.NewImageFromFile("img\\board.png")
	screen.DrawImage(boardImage, nil)

	//setup each image
	blackPawn, _, _ = ebitenutil.NewImageFromFile("img\\black\\pawn.png")
	blackBishop, _, _ = ebitenutil.NewImageFromFile("img\\black\\bishop.png")
	blackKing, _, _ = ebitenutil.NewImageFromFile("img\\black\\king.png")
	blackKnight, _, _ = ebitenutil.NewImageFromFile("img\\black\\knight.png")
	blackQueen, _, _ = ebitenutil.NewImageFromFile("img\\black\\queen.png")
	blackRook, _, _ = ebitenutil.NewImageFromFile("img\\black\\rook.png")
	whitePawn, _, _ = ebitenutil.NewImageFromFile("img\\white\\pawn.png")
	whiteBishop, _, _ = ebitenutil.NewImageFromFile("img\\white\\bishop.png")
	whiteKing, _, _ = ebitenutil.NewImageFromFile("img\\white\\king.png")
	whiteKnight, _, _ = ebitenutil.NewImageFromFile("img\\white\\knight.png")
	whiteQueen, _, _ = ebitenutil.NewImageFromFile("img\\white\\queen.png")
	whiteRook, _, _ = ebitenutil.NewImageFromFile("img\\white\\rook.png")

	//want markers to appear above all other, so need to store them so I can draw them last
	var markers []*ebiten.Image
	var markerop []*ebiten.DrawImageOptions

	//now draw each piece, based on the g.boardStateState
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := g.boardState[i][j]
			curOP := &ebiten.DrawImageOptions{}
			curOP.GeoM.Translate(float64(125*j), float64(125*i))
			flag := false
			if j == 1 {

			}
			if piece > 100 {
				curOP.ColorM.Translate(100, 0, 0, 0)
				g.boardState[i][j] -= 100
				piece -= 100
				//flag == true means piece selected
				flag = true
			}
			switch piece {
			case 1, 11:
				//PAWNS
				if piece == 1 {
					//white pawn
					screen.DrawImage(whitePawn, curOP)
				} else {
					//black pawn
					screen.DrawImage(blackPawn, curOP)
				}
				if flag {
					var offset float64
					var ioffset, startOffset int
					canMoveOne := false
					if (g.playerColor == 0 && piece != 11) || (g.playerColor == 1 && piece != 1) {
						//moves -i
						ioffset = -1
						startOffset = -2
						offset = -74.5
					} else {
						//moves +i
						ioffset = 1
						startOffset = 2
						offset = 180
					}
					if g.boardState[i+ioffset][j] == 0 {
						//empty, can move
						if (g.playerColor == 0 && piece == 1 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j))) || (g.playerColor == 1 && piece == 11 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j))) {
							markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, 0, 50.5, offset)
							canMoveOne = true
						}
					}
					if (g.playerColor == 0 && piece != 11) || (g.playerColor == 1 && piece != 1) {
						if i == 6 {
							//pawn hasn't moved, can move to spaces if empty
							if canMoveOne {
								if g.boardState[i+startOffset][j] == 0 {
									if (g.playerColor == 0 && piece == 1 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+startOffset, j))) || (g.playerColor == 1 && piece == 11 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+startOffset, j))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -1, 50.5, offset)
									}
								}
							}
						}
					} else {
						if i == 1 {
							//pawn hasn't moved, can move to spaces if empty
							if canMoveOne {
								if g.boardState[i+startOffset][j] == 0 {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, 1, 50.5, offset)
								}
							}
						}
					}
					if j+1 < 8 {
						if (g.boardState[i+ioffset][j+1] > 10 && piece == 1) || (g.boardState[i+ioffset][j+1] < 10 && g.boardState[i+ioffset][j+1] > 0 && piece == 11) {
							//enemy piece to the right, can take
							if (g.playerColor == 0 && piece == 1 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j+1))) || (g.playerColor == 1 && piece == 11 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j+1))) {
								markers, markerop = appendImageToSlice(markers, markerop, i, j, 1, 0, 50.5, offset)
							}
						}
					}
					if j-1 >= 0 {
						if (g.boardState[i+ioffset][j-1] > 10 && piece == 1) || (g.boardState[i+ioffset][j-1] < 10 && g.boardState[i+ioffset][j-1] > 0 && piece == 11) {
							//enemy piece to the left, can take
							if (g.playerColor == 0 && piece == 1 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j-1))) || (g.playerColor == 1 && piece == 11 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+ioffset, j-1))) {
								markers, markerop = appendImageToSlice(markers, markerop, i, j, -1, 0, 50.5, offset)
							}
						}
					}
				}
			case 2, 12:
				//BISHOPS
				if piece == 2 {
					//white bishop
					screen.DrawImage(whiteBishop, curOP)
				} else {
					//black bishop
					screen.DrawImage(blackBishop, curOP)
				}
				if flag {
					topLeft, topRight, bottomLeft, bottomRight := true, true, true, true
					for offset := 1; offset < 8; offset++ {
						if topLeft {
							if i-offset >= 0 && j-offset >= 0 {
								if g.boardState[i-offset][j-offset] != 0 {
									if (g.boardState[i-offset][j-offset] > 10 && piece == 2) || (g.boardState[i-offset][j-offset] < 10 && piece == 12) {
										//can take piece
										if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, -offset, 50.5, 53)
										}
									}
									topLeft = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, -offset, 50.5, 53)
									}
								}
							}
						}
						if topRight {
							if i-offset >= 0 && j+offset < 8 {
								if g.boardState[i-offset][j+offset] != 0 {
									if (g.boardState[i-offset][j+offset] > 10 && piece == 2) || (g.boardState[i-offset][j+offset] < 10 && piece == 12) {
										//can take piece
										if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, -offset, 50.5, 53)
										}
									}
									topRight = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, -offset, 50.5, 53)
									}
								}
							}
						}
						if bottomLeft {
							if i+offset < 8 && j-offset >= 0 {
								if g.boardState[i+offset][j-offset] != 0 {
									if (g.boardState[i+offset][j-offset] > 10 && piece == 2) || (g.boardState[i+offset][j-offset] < 10 && piece == 12) {
										//can take piece
										if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, offset, 50.5, 53)
										}
									}
									bottomLeft = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) {
									}
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, offset, 50.5, 53)
								}
							}
						}
						if bottomRight {
							if i+offset < 8 && j+offset < 8 {
								if g.boardState[i+offset][j+offset] != 0 {
									if (g.boardState[i+offset][j+offset] > 10 && piece == 2) || (g.boardState[i+offset][j+offset] < 10 && piece == 12) {
										//can take piece
										if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, offset, 50.5, 53)
										}
									}
									bottomRight = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 2 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) || (g.playerColor == 1 && piece == 12 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, offset, 50.5, 53)
									}
								}
							}
						}
					}
				}
			case 3, 13:
				//KNIGHTS
				if piece == 3 {
					//white knight
					screen.DrawImage(whiteKnight, curOP)
				} else {
					//black knight
					screen.DrawImage(blackKnight, curOP)
				}
				//check all spaces with abs(x1-x2) and abs(y1-y2) where xs are either 1 or 2, and y is the other (x is 1 then y is 2, y is 1 then x is 2)
				//can't think of a better solution than just checking each space individually
				if flag {
					if piece == 3 {
						if i+2 < 8 && j+1 < 8 {
							if g.boardState[i+2][j+1] == 0 || g.boardState[i+2][j+1] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j+1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j+1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 1, 2, 54.5, 53)
								}
							}
						}
						if i+2 < 8 && j-1 >= 0 {
							if g.boardState[i+2][j-1] == 0 || g.boardState[i+2][j-1] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j-1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j-1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -1, 2, 54.5, 53)
								}
							}
						}
						if i+1 < 8 && j+2 < 8 {
							if g.boardState[i+1][j+2] == 0 || g.boardState[i+1][j+2] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j+2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j+2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 1, 54.5, 53)
								}
							}
						}
						if i+1 < 8 && j-2 >= 0 {
							if g.boardState[i+1][j-2] == 0 || g.boardState[i+1][j-2] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j-2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j-2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 1, 54.5, 53)
								}
							}
						}
						if i-1 >= 0 && j+2 < 8 {
							if g.boardState[i-1][j+2] == 0 || g.boardState[i-1][j+2] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j+2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j+2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, -1, 54.5, 53)
								}
							}
						}
						if i-1 >= 0 && j-2 >= 0 {
							if g.boardState[i-1][j-2] == 0 || g.boardState[i-1][j-2] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j-2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j-2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, -1, 54.5, 53)
								}
							}
						}
						if i-2 >= 0 && j+1 < 8 {
							if g.boardState[i-2][j+1] == 0 || g.boardState[i-2][j+1] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j+1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j+1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 1, -2, 54.5, 53)
								}
							}
						}
						if i-2 >= 0 && j-1 >= 0 {
							if g.boardState[i-2][j-1] == 0 || g.boardState[i-2][j-1] > 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j-1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j-1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -2, -73.5, 53)
								}
							}
						}
					} else {
						if i+2 < 8 && j+1 < 8 {
							if g.boardState[i+2][j+1] == 0 || g.boardState[i+2][j+1] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j+1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j+1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 1, 2, 53.5, 53)
								}
							}
						}
						if i+2 < 8 && j-1 >= 0 {
							if g.boardState[i+2][j-1] == 0 || g.boardState[i+2][j-1] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j-1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+2, j-1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -1, 2, 50.5, 53)
								}
							}
						}
						if i+1 < 8 && j+2 < 8 {
							if g.boardState[i+1][j+2] == 0 || g.boardState[i+1][j+2] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j+2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j+2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 1, 54.5, 53)
								}
							}
						}
						if i+1 < 8 && j-2 >= 0 {
							if g.boardState[i+1][j-2] == 0 || g.boardState[i+1][j-2] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j-2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+1, j-2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 1, 54.5, 53)
								}
							}
						}
						if i-1 >= 0 && j+2 < 8 {
							if g.boardState[i-1][j+2] == 0 || g.boardState[i-1][j+2] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j+2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j+2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, -1, 54.5, 53)
								}
							}
						}
						if i-1 >= 0 && j-2 >= 0 {
							if g.boardState[i-1][j-2] == 0 || g.boardState[i-1][j-2] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j-2))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-1, j-2))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, -1, 54.5, 53)
								}
							}
						}
						if i-2 >= 0 && j+1 < 8 {
							if g.boardState[i-2][j+1] == 0 || g.boardState[i-2][j+1] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j+1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j+1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 1, -2, 52.5, 53)
								}
							}
						}
						if i-2 >= 0 && j-1 >= 0 {
							if g.boardState[i-2][j-1] == 0 || g.boardState[i-2][j-1] < 10 {
								if (g.playerColor == 0 && piece == 3 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j-1))) || (g.playerColor == 1 && piece == 13 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-2, j-1))) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -2, -73.5, 53)
								}
							}
						}
					}
				}
			case 4, 14:
				//ROOKS
				if piece == 4 {
					//white rook
					screen.DrawImage(whiteRook, curOP)
				} else {
					//black rook
					screen.DrawImage(blackRook, curOP)
				}
				if flag {
					up, down, left, right := true, true, true, true
					for offset := 1; offset < 8; offset++ {
						if up {
							if i-offset >= 0 {
								if g.boardState[i-offset][j] != 0 {
									if (g.boardState[i-offset][j] < 10 && piece > 10) || (g.boardState[i-offset][j] > 10 && piece < 10) {
										//rook can take a piece
										if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -offset, 50.5, 53)
										}
									}
									up = false
								} else {
									if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -offset, 50.5, 53)
									}
								}
							}
						}
						if down {
							if i+offset < 8 {
								if g.boardState[i+offset][j] != 0 {
									if (g.boardState[i+offset][j] < 10 && piece > 10) || (g.boardState[i+offset][j] > 10 && piece < 10) {
										//rook can take a piece
										if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, offset, 50.5, 53)
										}
									}
									down = false
								} else {
									if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, offset, 50.5, 53)
									}
								}
							}
						}
						if left {
							if j-offset >= 0 {
								if g.boardState[i][j-offset] != 0 {
									if (g.boardState[i][j-offset] < 10 && piece > 10) || (g.boardState[i][j-offset] > 10 && piece < 10) {
										//rook can take a piece
										if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, 0, 50.5, 53)
										}
									}
									left = false
								} else {
									if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, 0, 50.5, 53)
									}
								}
							}
						}
						if right {
							if j+offset < 8 {
								if g.boardState[i][j+offset] != 0 {
									if (g.boardState[i][j+offset] < 10 && piece > 10) || (g.boardState[i][j+offset] > 10 && piece < 10) {
										//rook can take a piece
										if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, 0, 50.5, 53)
										}
									}
									right = false
								} else {
									if (g.playerColor == 0 && piece == 4 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) || (g.playerColor == 1 && piece == 14 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, 0, 50.5, 53)
									}
								}
							}
						}
					}
				}
			case 5, 15:
				//QUEENS
				if piece == 5 {
					//white queen
					screen.DrawImage(whiteQueen, curOP)
				} else {
					//black queen
					screen.DrawImage(blackQueen, curOP)
				}
				if flag {
					//just combine the bishop and rook logic
					up, down, left, right, topRight, topLeft, bottomRight, bottomLeft := true, true, true, true, true, true, true, true
					for offset := 1; offset < 8; offset++ {
						if up {
							if i-offset >= 0 {
								if g.boardState[i-offset][j] != 0 {
									if (g.boardState[i-offset][j] < 10 && piece > 10) || (g.boardState[i-offset][j] > 10 && piece < 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -offset, 50.5, 53)
										}
									}
									up = false
								} else {
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, -offset, 50.5, 53)
									}
								}
							}
						}
						if down {
							if i+offset < 8 {
								if g.boardState[i+offset][j] != 0 {
									if (g.boardState[i+offset][j] < 10 && piece > 10) || (g.boardState[i+offset][j] > 10 && piece < 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, offset, 50.5, 53)
										}
									}
									down = false
								} else {
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, 0, offset, 50.5, 53)
									}
								}
							}
						}
						if left {
							if j-offset >= 0 {
								if g.boardState[i][j-offset] != 0 {
									if (g.boardState[i][j-offset] < 10 && piece > 10) || (g.boardState[i][j-offset] > 10 && piece < 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, 0, 50.5, 53)
										}
									}
									left = false
								} else {
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j-offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, 0, 50.5, 53)
									}
								}
							}
						}
						if right {
							if j+offset < 8 {
								if g.boardState[i][j+offset] != 0 {
									if (g.boardState[i][j+offset] < 10 && piece > 10) || (g.boardState[i][j+offset] > 10 && piece < 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, 0, 50.5, 53)
										}
									}
									right = false
								} else {
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, 0, 50.5, 53)
									}
								}
							}
						}
						if topLeft {
							if i-offset >= 0 && j-offset >= 0 {
								if g.boardState[i-offset][j-offset] != 0 {
									if (g.boardState[i-offset][j-offset] > 10 && piece < 10) || (g.boardState[i-offset][j-offset] < 10 && piece > 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, -offset, 50.5, 53)
										}
									}
									topLeft = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j-offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, -offset, 50.5, 53)
									}
								}
							}
						}
						if topRight {
							if i-offset >= 0 && j+offset < 8 {
								if g.boardState[i-offset][j+offset] != 0 {
									if (g.boardState[i-offset][j+offset] > 10 && piece < 10) || (g.boardState[i-offset][j+offset] < 10 && piece > 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, -offset, 50.5, 53)
										}
									}
									topRight = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i-offset, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, -offset, 50.5, 53)
									}
								}
							}
						}
						if bottomLeft {
							if i+offset < 8 && j-offset >= 0 {
								if g.boardState[i+offset][j-offset] != 0 {
									if (g.boardState[i+offset][j-offset] > 10 && piece < 10) || (g.boardState[i+offset][j-offset] < 10 && piece > 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, offset, 50.5, 53)
										}
									}
									bottomLeft = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j-offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, -offset, offset, 50.5, 53)
									}
								}
							}
						}
						if bottomRight {
							if i+offset < 8 && j+offset < 8 {
								if g.boardState[i+offset][j+offset] != 0 {
									if (g.boardState[i+offset][j+offset] > 10 && piece < 10) || (g.boardState[i+offset][j+offset] < 10 && piece > 10) {
										//queen can take a piece
										if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) {
											markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, offset, 50.5, 53)
										}
									}
									bottomRight = false
								} else {
									//hit an empty spot
									if (g.playerColor == 0 && piece == 5 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) || (g.playerColor == 1 && piece == 15 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+offset, j+offset))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, offset, offset, 50.5, 53)
									}
								}
							}
						}
					}
				}
			case 6, 16:
				//KINGS
				if piece == 6 {
					//white king
					screen.DrawImage(whiteKing, curOP)
				} else {
					//black king
					screen.DrawImage(blackKing, curOP)
				}
				if flag {
					for a := -1; a <= 1; a++ {
						for b := -1; b <= 1; b++ {
							if i+a < 8 && i+a >= 0 && j+b < 8 && j+b >= 0 {
								if g.boardState[i+a][j+b] == 0 || (g.boardState[i+a][j+b] < 10 && piece > 10) || (g.boardState[i+a][j+b] > 10 && piece < 10) {
									//empty space or enemy
									if (g.playerColor == 0 && piece == 6 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+a, j+b))) || (g.playerColor == 1 && piece == 16 && !isInCheckPlayer(g, createTempBoardMove(g.boardState, i, j, i+a, j+b))) {
										markers, markerop = appendImageToSlice(markers, markerop, i, j, b, a, 50.5, 53)
									}
								}
							}
						}
					}
					if piece == 6 {
						//white king, check for castle
						if whiteKingCastle {
							if g.playerColor == 0 {
								if canMoveHere(g, 4, 7, 6, 7, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 0, 50.5, 53)
								}
							} else {
								if canMoveHere(g, 3, 0, 1, 0, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 0, 50.5, 53)
								}
							}
						}
						if whiteQueenCastle {
							if g.playerColor == 0 {
								if canMoveHere(g, 4, 7, 2, 7, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 0, 50.5, 53)
								}
							} else {
								if canMoveHere(g, 3, 0, 5, 0, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 0, 50.5, 53)
								}
							}
						}
					} else {
						//black king, check for castle
						if blackKingCastle {
							if g.playerColor == 0 {
								if canMoveHere(g, 4, 0, 6, 0, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 0, 50.5, 53)
								}
							} else {
								if canMoveHere(g, 3, 7, 1, 7, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 0, 50.5, 53)
								}
							}
						}
						if blackQueenCastle {
							if g.playerColor == 0 {
								if canMoveHere(g, 4, 0, 2, 0, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, -2, 0, 50.5, 53)
								}
							} else {
								if canMoveHere(g, 3, 7, 5, 7, g.boardState, false) {
									markers, markerop = appendImageToSlice(markers, markerop, i, j, 2, 0, 50.5, 53)
								}
							}
						}
					}
				}
			}
			if flag {
				g.boardState[i][j] += 100
				flag = false
			}
		}
	}
	for i := 0; i < len(markers); i++ {
		screen.DrawImage(markers[i], markerop[i])
	}
	if g.gameOver == 0 {
		whiteWin, _, _ := ebitenutil.NewImageFromFile("img\\whiteWin.png")
		winOP := &ebiten.DrawImageOptions{}
		winOP.GeoM.Translate(377, 437)
		screen.DrawImage(whiteWin, winOP)
	}
	if g.gameOver == 1 {
		blackWin, _, _ := ebitenutil.NewImageFromFile("img\\blackWin.png")
		winOP := &ebiten.DrawImageOptions{}
		winOP.GeoM.Translate(377, 437)
		screen.DrawImage(blackWin, winOP)
	}
}

/*
Layout ... Sets up the layout of the screen
*/
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return the game screen size
	return 1000, 1000
}

//Creates a new fresh chess game
func freshBoard() *Game {
	//boardState:
	//0 is empty
	//white: 1 is pawn, 2 is bishop, 3 is knight, 4 is rook, 5 is queen, 6 is king
	//black: 11 is pawn, 12 is bishop, 13 is knight, 14 is rook, 15 is queen, 16 is king

	//playerColor: 0 is white at bottom, 1 is black at bottom
	color := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(2)
	//color := 1
	var newBoard [8][8]int
	var curPiece int
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if i == 0 {
				if color == 0 {
					switch j {
					case 0:
						curPiece = 14
					case 1:
						curPiece = 13
					case 2:
						curPiece = 12
					case 3:
						curPiece = 15
					case 4:
						curPiece = 16
					case 5:
						curPiece = 12
					case 6:
						curPiece = 13
					case 7:
						curPiece = 14
					}
				} else {
					switch j {
					case 0:
						curPiece = 4
					case 1:
						curPiece = 3
					case 2:
						curPiece = 2
					case 3:
						curPiece = 6
					case 4:
						curPiece = 5
					case 5:
						curPiece = 2
					case 6:
						curPiece = 3
					case 7:
						curPiece = 4
					}
				}
			} else if i == 1 {
				if color == 0 {
					curPiece = 11
				} else {
					curPiece = 1
				}
			} else if i == 6 {
				if color == 0 {
					curPiece = 1
				} else {
					curPiece = 11
				}
			} else if i == 7 {
				if color == 1 {
					switch j {
					case 0:
						curPiece = 14
					case 1:
						curPiece = 13
					case 2:
						curPiece = 12
					case 3:
						curPiece = 16
					case 4:
						curPiece = 15
					case 5:
						curPiece = 12
					case 6:
						curPiece = 13
					case 7:
						curPiece = 14
					}
				} else {
					switch j {
					case 0:
						curPiece = 4
					case 1:
						curPiece = 3
					case 2:
						curPiece = 2
					case 3:
						curPiece = 5
					case 4:
						curPiece = 6
					case 5:
						curPiece = 2
					case 6:
						curPiece = 3
					case 7:
						curPiece = 4
					}
				}
			} else {
				curPiece = 0
			}
			newBoard[i][j] = curPiece
		}
	}
	g := &Game{
		playerColor: color,
		boardState:  newBoard,
		playerTurn:  0,
		gameOver:    2,
	}
	return g
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chess!")
	whiteKingCastle, whiteQueenCastle, blackKingCastle, blackQueenCastle = true, true, true, true
	g := freshBoard()
	aiInit(g)
	if g.playerColor == 1 {
		aiMove(g)
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
