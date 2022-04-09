/*
ai.go
By: Matthew Yungbluth
CSc 372
Final Project
Due: December 7th, 2020 at the beginning of class

This file controls the ai that complements chess.go
*/

package main

/*
values for each piece:
empty square: 0
pawn: 100
knight and bishop: 320/330
rook: 500
queen: 900
king: 20000
*/
var values = [...]int{0, 100, 330, 320, 500, 900, 20000}

/*
positioning values taken from internet
*/
var pawnPos = [8][8]int{{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 10, 10, -20, -20, 10, 10, 5},
	{5, -5, -10, 0, 0, -10, -5, 5},
	{0, 0, 0, 20, 20, 0, 0, 0},
	{5, 5, 10, 25, 25, 10, 5, 5},
	{10, 10, 20, 30, 30, 20, 10, 10},
	{50, 50, 50, 50, 50, 50, 50, 50},
	{0, 0, 0, 0, 0, 0, 0, 0}}

var knightPos = [8][8]int{{-50, -40, -30, -30, -30, -30, -40, -50},
	{-40, -20, 0, 5, 5, 0, -20, -40},
	{-30, 5, 10, 15, 15, 10, 5, -30},
	{-30, 0, 15, 20, 20, 15, 0, -30},
	{-30, 5, 15, 20, 20, 15, 5, -30},
	{-30, 0, 10, 15, 15, 10, 0, -30},
	{-40, -20, 0, 0, 0, 0, -20, -40},
	{-50, -40, -30, -30, -30, -30, -40, -50}}

var bishopPos = [8][8]int{{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 5, 0, 0, 0, 0, 5, -10},
	{-10, 10, 10, 10, 10, 10, 10, -10},
	{-10, 0, 10, 10, 10, 10, 0, -10},
	{-10, 5, 5, 10, 10, 5, 5, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20}}

var rookPos = [8][8]int{{0, 0, 0, 5, 5, 0, 0, 0},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0}}

var queenPos = [8][8]int{{-20, -10, -10, -5, -5, -10, -10, -20},
	{-10, 0, 5, 0, 0, 0, 0, -10},
	{-10, 5, 5, 5, 5, 5, 0, -10},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-10, 0, 5, 5, 5, 5, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -5, -5, -10, -10, -20}}

var kingPos = [8][8]int{{20, 30, 10, 0, 0, 10, 30, 20},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{-10, -20, -20, -20, -20, -20, -20, -10},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30}}

type moveGen struct {
	thisState [8][8]int
	aiTurn    bool
	children  []*moveGen
}

//points are calculated by adding up all current pieces on the board minus the king (game is over when king is captured so doesn't matter)
var aiPoints = 30
var playerPoints = 30
var aiMoves *moveGen
var isAiTurn bool
var aiThought bool
var depthCalc int

//AI needs to move, move with the current calculated best move
func aiMove(g *Game) {
	if g.gameOver == 2 {
		aiThink(g)
		curMax := -9999
		var bestI int
		for i, child := range aiMoves.children {
			thisVal := decideWhichMove(child, g, -9999, 9999)
			if thisVal > curMax {
				curMax = thisVal
				bestI = i
			}
		}
		//child at bestI is best move to make, make it
		g.boardState = aiMoves.children[bestI].thisState
		blackKing := false
		whiteKing := false
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if g.boardState[i][j] == 16 {
					blackKing = true
				}
				if g.boardState[i][j] == 6 {
					whiteKing = true
				}
			}
		}
		if blackKing == false {
			//white won
			if g.gameOver == 2 {
				g.gameOver = 0
			}
		}
		if whiteKing == false {
			//black won
			if g.gameOver == 2 {
				g.gameOver = 1
			}
		}
		g.playerTurn = g.playerColor
	}
}

//Gets the value of the piece at position i,j for the ai
func getPosValue(board [8][8]int, playerColor, i, j int) int {
	if playerColor == 0 {
		//black pieces are ai pieces
		if board[i][j] < 10 {
			return values[board[i][j]]
		}
		return -values[board[i][j]-10]
	}
	//white pieces are ai pieces
	if board[i][j] < 10 {
		return -values[board[i][j]]
	}
	return values[board[i][j]-10]
}

//returns the value of the board with positive for the ai
func getCurrentBoardValueAi(board [8][8]int, g *Game) int {
	value := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] != 0 {
				if board[i][j] > 10 {
					//black piece
					if g.playerColor == 0 {
						//player is white, ai is black
						value += values[board[i][j]-10]
						switch board[i][j] - 10 {
						case 1:
							value += pawnPos[i][j]
						case 2:
							value += bishopPos[i][j]
						case 3:
							value += knightPos[i][j]
						case 4:
							value += rookPos[i][j]
						case 5:
							value += queenPos[i][j]
						case 6:
							value += kingPos[i][j]
						}
					} else {
						//player is black, ai is white
						value -= values[board[i][j]-10]
						switch board[i][j] - 10 {
						case 1:
							value -= pawnPos[7-i][j]
						case 2:
							value -= bishopPos[7-i][j]
						case 3:
							value -= knightPos[7-i][j]
						case 4:
							value -= rookPos[7-i][j]
						case 5:
							value -= queenPos[7-i][j]
						case 6:
							value -= kingPos[7-i][j]
						}
					}
				} else {
					//white piece
					if g.playerColor == 0 {
						//player is white, ai is black
						value -= values[board[i][j]]
						switch board[i][j] {
						case 1:
							value -= pawnPos[7-i][j]
						case 2:
							value -= bishopPos[7-i][j]
						case 3:
							value -= knightPos[7-i][j]
						case 4:
							value -= rookPos[7-i][j]
						case 5:
							value -= queenPos[7-i][j]
						case 6:
							value -= kingPos[7-i][j]
						}
					} else {
						//player is black, ai is white
						value += values[board[i][j]]
						switch board[i][j] {
						case 1:
							value += pawnPos[i][j]
						case 2:
							value += bishopPos[i][j]
						case 3:
							value += knightPos[i][j]
						case 4:
							value += rookPos[i][j]
						case 5:
							value += queenPos[i][j]
						case 6:
							value += kingPos[i][j]
						}
					}
				}
			}
		}
	}
	return value
}

//returns the value of the board with positive for the ai
func getCurrentBoardValuePlayer(board [8][8]int, g *Game) int {
	value := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] != 0 {
				if board[i][j] > 10 {
					//black piece
					if g.playerColor == 1 {
						//player is black, ai is white
						value += values[board[i][j]-10]
						switch board[i][j] - 10 {
						case 1:
							value += pawnPos[i][j]
						case 2:
							value += bishopPos[i][j]
						case 3:
							value += knightPos[i][j]
						case 4:
							value += rookPos[i][j]
						case 5:
							value += queenPos[i][j]
						case 6:
							value += kingPos[i][j]
						}
					} else {
						//player is white, ai is black
						value -= values[board[i][j]-10]
						switch board[i][j] - 10 {
						case 1:
							value -= pawnPos[7-i][j]
						case 2:
							value -= bishopPos[7-i][j]
						case 3:
							value -= knightPos[7-i][j]
						case 4:
							value -= rookPos[7-i][j]
						case 5:
							value -= queenPos[7-i][j]
						case 6:
							value -= kingPos[7-i][j]
						}
					}
				} else {
					//white piece
					if g.playerColor == 1 {
						//player is black, ai is white
						value -= values[board[i][j]]
						switch board[i][j] {
						case 1:
							value -= pawnPos[7-i][j]
						case 2:
							value -= bishopPos[7-i][j]
						case 3:
							value -= knightPos[7-i][j]
						case 4:
							value -= rookPos[7-i][j]
						case 5:
							value -= queenPos[7-i][j]
						case 6:
							value -= kingPos[7-i][j]
						}
					} else {
						//player is white, ai is black
						value += values[board[i][j]]
						switch board[i][j] {
						case 1:
							value += pawnPos[i][j]
						case 2:
							value += bishopPos[i][j]
						case 3:
							value += knightPos[i][j]
						case 4:
							value += rookPos[i][j]
						case 5:
							value += queenPos[i][j]
						case 6:
							value += kingPos[i][j]
						}
					}
				}
			}
		}
	}
	return value
}

//returns all moves possible for a current piece, basically ripped from the functionality used to create the red squares, maybe should find a way to adapt it so they can both use same code?
func getAllMovesPiece(board [8][8]int, playerColor, i, j int, g *Game) [][]int {
	//internal representation of moves is a slice of slices where each inner slice is a possible move and each inner slice is [i][j][value] pf the move,
	//value is taken from the value slice compared to the piece at the location
	var moves [][]int
	piece := board[i][j]
	switch piece {
	case 1, 11:
		//PAWNS
		var ioffset, startOffset int
		canMoveOne := false
		if (playerColor == 0 && piece != 11) || (playerColor == 1 && piece == 11) {
			//moves -i
			ioffset = -1
			startOffset = -2
		} else {
			//moves +i
			ioffset = 1
			startOffset = 2
		}
		if i+ioffset < 8 && i+ioffset >= 0 {
			if board[i+ioffset][j] == 0 {
				//empty, can move
				pos := []int{i + ioffset, j, 0}
				moves = append(moves, pos)
				canMoveOne = true
			}
		}
		if playerColor == 0 && piece != 11 {
			if i == 6 {
				//pawn hasn't moved, can move to spaces if empty
				if canMoveOne {
					if i+startOffset < 8 {
						if board[i+startOffset][j] == 0 {
							pos := []int{i + (ioffset * 2), j, 0}
							moves = append(moves, pos)
						}
					}
				}
			}
		} else {
			if i == 1 {
				//pawn hasn't moved, can move to spaces if empty
				if canMoveOne {
					if i+startOffset < 8 && i+startOffset >= 0 {
						if board[i+startOffset][j] == 0 {
							pos := []int{i + (ioffset * 2), j, 0}
							moves = append(moves, pos)
						}
					}
				}
			}
		}
		if j+1 < 8 && i+ioffset < 8 && i+ioffset >= 0 {
			if (board[i+ioffset][j+1] > 10 && piece == 1) || (board[i+ioffset][j+1] < 10 && board[i+ioffset][j+1] > 0 && piece == 11) {
				//enemy piece to the right, can take
				if board[i+ioffset][j+1] < 10 {
					pos := []int{i + ioffset, j + 1, values[board[i+ioffset][j+1]]}
					moves = append(moves, pos)
				} else {
					pos := []int{i + ioffset, j + 1, values[board[i+ioffset][j+1]-10]}
					moves = append(moves, pos)
				}
			}
		}
		if j-1 >= 0 && i+ioffset < 8 && i+ioffset >= 0 {
			if (board[i+ioffset][j-1] > 10 && piece == 1) || (board[i+ioffset][j-1] < 10 && board[i+ioffset][j-1] > 0 && piece == 11) {
				//enemy piece to the left, can take
				if board[i+ioffset][j-1] < 10 {
					pos := []int{i + ioffset, j - 1, values[board[i+ioffset][j-1]]}
					moves = append(moves, pos)
				} else {
					pos := []int{i + ioffset, j - 1, values[board[i+ioffset][j-1]-10]}
					moves = append(moves, pos)
				}
			}
		}
	case 2, 12:
		//BISHOPS
		topLeft, topRight, bottomLeft, bottomRight := true, true, true, true
		for offset := 1; offset < 8; offset++ {
			if topLeft {
				if i-offset >= 0 && j-offset >= 0 && i-offset < 8 && j-offset < 8 {
					if board[i-offset][j-offset] != 0 {
						if (board[i-offset][j-offset] > 10 && piece == 2) || (board[i-offset][j-offset] < 10 && piece == 12) {
							//can take piece
							if board[i-offset][j-offset] < 10 {
								pos := []int{i - offset, j - offset, values[board[i-offset][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j - offset, values[board[i-offset][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						topLeft = false
					} else {
						//hit an empty spot
						pos := []int{i - offset, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if topRight {
				if i-offset >= 0 && j+offset < 8 && i-offset < 8 && j+offset >= 0 {
					if board[i-offset][j+offset] != 0 {
						if (board[i-offset][j+offset] > 10 && piece == 2) || (board[i-offset][j+offset] < 10 && piece == 12) {
							//can take piece
							if board[i-offset][j+offset] < 10 {
								pos := []int{i - offset, j + offset, values[board[i-offset][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j + offset, values[board[i-offset][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						topRight = false
					} else {
						//hit an empty spot
						pos := []int{i - offset, j + offset, 0}
						moves = append(moves, pos)
					}
				} else {
					topRight = false
				}
			}
			if bottomLeft {
				if i+offset < 8 && j-offset >= 0 && i+offset >= 0 && j-offset < 8 {
					if board[i+offset][j-offset] != 0 {
						if (board[i+offset][j-offset] > 10 && piece == 2) || (board[i+offset][j-offset] < 10 && piece == 12) {
							//can take piece
							if board[i+offset][j-offset] < 10 {
								pos := []int{i + offset, j - offset, values[board[i+offset][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j - offset, values[board[i+offset][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						bottomLeft = false
					} else {
						//hit an empty spot
						pos := []int{i + offset, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if bottomRight {
				if i+offset < 8 && j+offset < 8 && i+offset >= 0 && j+offset >= 0 {
					if board[i+offset][j+offset] != 0 {
						if (board[i+offset][j+offset] > 10 && piece == 2) || (board[i+offset][j+offset] < 10 && piece == 12) {
							//white bishop can take black piece
							if board[i+offset][j+offset] < 10 {
								pos := []int{i + offset, j + offset, values[board[i+offset][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j + offset, values[board[i+offset][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						bottomRight = false
					} else {
						//hit an empty spot
						pos := []int{i + offset, j + offset, 0}
						moves = append(moves, pos)
					}
				}
			}
		}
	case 3, 13:
		//KNIGHTS
		//check all spaces with abs(x1-x2) and abs(y1-y2) where xs are either 1 or 2, and y is the other (x is 1 then y is 2, y is 1 then x is 2)
		//can't think of a better solution than just checking each space individually
		if piece == 3 {
			if i+2 < 8 && j+1 < 8 {
				if board[i+2][j+1] == 0 || board[i+2][j+1] > 10 {
					if board[i+2][j+1] > 10 {
						pos := []int{i + 2, j + 1, values[board[i+2][j+1]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i + 2, j + 1, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i+2 < 8 && j-1 >= 0 {
				if board[i+2][j-1] == 0 || board[i+2][j-1] > 10 {
					if board[i+2][j-1] > 10 {
						pos := []int{i + 2, j - 1, values[board[i+2][j-1]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i + 2, j - 1, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i+1 < 8 && j+2 < 8 {
				if board[i+1][j+2] == 0 || board[i+1][j+2] > 10 {
					if board[i+1][j+2] > 10 {
						pos := []int{i + 1, j + 2, values[board[i+1][j+2]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i + 1, j + 2, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i+1 < 8 && j-2 >= 0 {
				if board[i+1][j-2] == 0 || board[i+1][j-2] > 10 {
					if board[i+1][j-2] > 10 {
						pos := []int{i + 1, j - 2, values[board[i+1][j-2]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i + 1, j - 2, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i-1 >= 0 && j+2 < 8 {
				if board[i-1][j+2] == 0 || board[i-1][j+2] > 10 {
					if board[i-1][j+2] > 10 {
						pos := []int{i - 1, j + 2, values[board[i-1][j+2]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i - 1, j + 2, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i-1 >= 0 && j-2 >= 0 {
				if board[i-1][j-2] == 0 || board[i-1][j-2] > 10 {
					if board[i-1][j-2] > 10 {
						pos := []int{i - 1, j - 2, values[board[i-1][j-2]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i - 1, j - 2, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i-2 >= 0 && j+1 < 8 {
				if board[i-2][j+1] == 0 || board[i-2][j+1] > 10 {
					if board[i-2][j+1] > 10 {
						pos := []int{i - 2, j + 1, values[board[i-2][j+1]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i - 2, j + 1, 0}
						moves = append(moves, pos)
					}
				}
			}
			if i-2 >= 0 && j-1 >= 0 {
				if board[i-2][j-1] == 0 || board[i-2][j-1] > 10 {
					if board[i-2][j-1] > 10 {
						pos := []int{i - 2, j - 1, values[board[i-2][j-1]-10]}
						moves = append(moves, pos)
					} else {
						pos := []int{i - 2, j - 1, 0}
						moves = append(moves, pos)
					}
				}
			}
		} else {
			if i+2 < 8 && j+1 < 8 {
				if board[i+2][j+1] == 0 || board[i+2][j+1] < 10 {
					pos := []int{i + 2, j + 1, values[board[i+2][j+1]]}
					moves = append(moves, pos)
				}
			}
			if i+2 < 8 && j-1 >= 0 {
				if board[i+2][j-1] == 0 || board[i+2][j-1] < 10 {
					pos := []int{i + 2, j - 1, values[board[i+2][j-1]]}
					moves = append(moves, pos)
				}
			}
			if i+1 < 8 && j+2 < 8 {
				if board[i+1][j+2] == 0 || board[i+1][j+2] < 10 {
					pos := []int{i + 1, j + 2, values[board[i+1][j+2]]}
					moves = append(moves, pos)
				}
			}
			if i+1 < 8 && j-2 >= 0 {
				if board[i+1][j-2] == 0 || board[i+1][j-2] < 10 {
					pos := []int{i + 1, j - 2, values[board[i+1][j-2]]}
					moves = append(moves, pos)
				}
			}
			if i-1 >= 0 && j+2 < 8 {
				if board[i-1][j+2] == 0 || board[i-1][j+2] < 10 {
					pos := []int{i - 1, j + 2, values[board[i-1][j+2]]}
					moves = append(moves, pos)
				}
			}
			if i-1 >= 0 && j-2 >= 0 {
				if board[i-1][j-2] == 0 || board[i-1][j-2] < 10 {
					pos := []int{i - 1, j - 2, values[board[i-1][j-2]]}
					moves = append(moves, pos)
				}
			}
			if i-2 >= 0 && j+1 < 8 {
				if board[i-2][j+1] == 0 || board[i-2][j+1] < 10 {
					pos := []int{i - 2, j + 1, values[board[i-2][j+1]]}
					moves = append(moves, pos)
				}
			}
			if i-2 >= 0 && j-1 >= 0 {
				if board[i-2][j-1] == 0 || board[i-2][j-1] < 10 {
					pos := []int{i - 2, j - 1, values[board[i-2][j-1]]}
					moves = append(moves, pos)
				}
			}
		}
	case 4, 14:
		//ROOKS
		up, down, left, right := true, true, true, true
		for offset := 1; offset < 8; offset++ {
			if up {
				if i-offset >= 0 && i-offset < 8 {
					if board[i-offset][j] != 0 {
						if (board[i-offset][j] < 10 && piece > 10) || (board[i-offset][j] > 10 && piece < 10) {
							//rook can take a piece
							if board[i-offset][j] < 10 {
								pos := []int{i - offset, j, values[board[i-offset][j]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j, values[board[i-offset][j]-10]}
								moves = append(moves, pos)
							}
						}
						up = false
					} else {
						pos := []int{i - offset, j, 0}
						moves = append(moves, pos)
					}
				}
			}
			if down {
				if i+offset < 8 && i+offset >= 0 {
					if board[i+offset][j] != 0 {
						if (board[i+offset][j] < 10 && piece > 10) || (board[i+offset][j] > 10 && piece < 10) {
							//rook can take a piece
							if board[i+offset][j] < 10 {
								pos := []int{i + offset, j, values[board[i+offset][j]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j, values[board[i+offset][j]-10]}
								moves = append(moves, pos)
							}
						}
						down = false
					} else {
						pos := []int{i + offset, j, 0}
						moves = append(moves, pos)
					}
				}
			}
			if left {
				if j-offset >= 0 && j-offset < 8 {
					if board[i][j-offset] != 0 {
						if (board[i][j-offset] < 10 && piece > 10) || (board[i][j-offset] > 10 && piece < 10) {
							//rook can take a piece
							if board[i][j-offset] < 10 {
								pos := []int{i, j - offset, values[board[i][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i, j - offset, values[board[i][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						left = false
					} else {
						pos := []int{i, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if right {
				if j+offset < 8 && j+offset >= 0 {
					if board[i][j+offset] != 0 {
						if (board[i][j+offset] < 10 && piece > 10) || (board[i][j+offset] > 10 && piece < 10) {
							//rook can take a piece
							if board[i][j+offset] < 10 {
								pos := []int{i, j + offset, values[board[i][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i, j + offset, values[board[i][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						right = false
					} else {
						pos := []int{i, j + offset, 0}
						moves = append(moves, pos)
					}
				}
			}
		}
	case 5, 15:
		//QUEENS
		//just combine the bishop and rook logic
		up, down, left, right, topRight, topLeft, bottomRight, bottomLeft := true, true, true, true, true, true, true, true
		for offset := 1; offset < 8; offset++ {
			if up {
				if i-offset >= 0 && i-offset < 8 {
					if board[i-offset][j] != 0 {
						if (board[i-offset][j] < 10 && piece > 10) || (board[i-offset][j] > 10 && piece < 10) {
							//queen can take a piece
							if board[i-offset][j] < 10 {
								pos := []int{i - offset, j, values[board[i-offset][j]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j, values[board[i-offset][j]-10]}
								moves = append(moves, pos)
							}
						}
						up = false
					} else {
						pos := []int{i - offset, j, 0}
						moves = append(moves, pos)
					}
				}
			}
			if down {
				if i+offset < 8 && i+offset >= 0 {
					if board[i+offset][j] != 0 {
						if (board[i+offset][j] < 10 && piece > 10) || (board[i+offset][j] > 10 && piece < 10) {
							//queen can take a piece
							if board[i+offset][j] < 10 {
								pos := []int{i + offset, j, values[board[i+offset][j]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j, values[board[i+offset][j]-10]}
								moves = append(moves, pos)
							}
						}
						down = false
					} else {
						pos := []int{i + offset, j, 0}
						moves = append(moves, pos)
					}
				}
			}
			if left {
				if j-offset >= 0 && j-offset < 8 {
					if board[i][j-offset] != 0 {
						if (board[i][j-offset] < 10 && piece > 10) || (board[i][j-offset] > 10 && piece < 10) {
							//queen can take a piece
							if board[i][j-offset] < 10 {
								pos := []int{i, j - offset, values[board[i][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i, j - offset, values[board[i][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						left = false
					} else {
						pos := []int{i, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if right {
				if j+offset < 8 && j+offset >= 0 {
					if board[i][j+offset] != 0 {
						if (board[i][j+offset] < 10 && piece > 10) || (board[i][j+offset] > 10 && piece < 10) {
							//queen can take a piece
							if board[i][j+offset] < 10 {
								pos := []int{i, j + offset, values[board[i][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i, j + offset, values[board[i][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						right = false
					} else {
						pos := []int{i, j + offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if topLeft {
				if i-offset >= 0 && j-offset >= 0 && i-offset < 8 && j-offset < 8 {
					if board[i-offset][j-offset] != 0 {
						if (board[i-offset][j-offset] > 10 && piece < 10) || (board[i-offset][j-offset] < 10 && piece > 10) {
							//queen can take a piece
							if board[i-offset][j-offset] < 10 {
								pos := []int{i - offset, j - offset, values[board[i-offset][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j - offset, values[board[i-offset][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						topLeft = false
					} else {
						//hit an empty spot
						pos := []int{i - offset, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if topRight {
				if i-offset >= 0 && j+offset < 8 && i-offset < 8 && j+offset >= 0 {
					if board[i-offset][j+offset] != 0 {
						if (board[i-offset][j+offset] > 10 && piece < 10) || (board[i-offset][j+offset] < 10 && piece > 10) {
							//queen can take a piece
							if board[i-offset][j+offset] < 10 {
								pos := []int{i - offset, j + offset, values[board[i-offset][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i - offset, j + offset, values[board[i-offset][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						topRight = false
					} else {
						//hit an empty spot
						pos := []int{i - offset, j + offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if bottomLeft {
				if i+offset < 8 && j-offset >= 0 && i+offset >= 0 && j-offset < 8 {
					if board[i+offset][j-offset] != 0 {
						if (board[i+offset][j-offset] > 10 && piece < 10) || (board[i+offset][j-offset] < 10 && piece > 10) {
							//queen can take a piece
							if board[i+offset][j-offset] < 10 {
								pos := []int{i + offset, j - offset, values[board[i+offset][j-offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j - offset, values[board[i+offset][j-offset]-10]}
								moves = append(moves, pos)
							}
						}
						bottomLeft = false
					} else {
						//hit an empty spot
						pos := []int{i + offset, j - offset, 0}
						moves = append(moves, pos)
					}
				}
			}
			if bottomRight {
				if i+offset < 8 && j+offset < 8 && i+offset >= 0 && j+offset >= 0 {
					if board[i+offset][j+offset] != 0 {
						if (board[i+offset][j+offset] > 10 && piece < 10) || (board[i+offset][j+offset] < 10 && piece > 10) {
							//queen can take a piece
							if board[i+offset][j+offset] < 10 {
								pos := []int{i + offset, j + offset, values[board[i+offset][j+offset]]}
								moves = append(moves, pos)
							} else {
								pos := []int{i + offset, j + offset, values[board[i+offset][j+offset]-10]}
								moves = append(moves, pos)
							}
						}
						bottomRight = false
					} else {
						//hit an empty spot
						pos := []int{i + offset, j + offset, 0}
						moves = append(moves, pos)
					}
				}
			}
		}
	case 6, 16:
		//KINGS
		for a := -1; a <= 1; a++ {
			for b := -1; b <= 1; b++ {
				if i+a < 8 && i+a >= 0 && j+b < 8 && j+b >= 0 {
					if board[i+a][j+b] == 0 || (board[i+a][j+b] < 10 && piece > 10) || (board[i+a][j+b] > 10 && piece < 10) {
						//empty space or enemy
						if board[i+a][j+b] < 10 {
							pos := []int{i + a, j + b, values[board[i+a][j+b]]}
							moves = append(moves, pos)
						} else {
							pos := []int{i + a, j + b, values[board[i+a][j+b]-10]}
							moves = append(moves, pos)
						}
					}
				}
			}
		}
	}
	return moves
}

//returns the coordinates of all ai pieces in a 2d array in [[i,j],[i,j]] format
func getAllPieces(board [8][8]int, playerColor int) [][]int {
	var pieces [][]int
	if (playerColor == 0 && isAiTurn) || (playerColor == 1 && !isAiTurn) {
		//black pieces
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 10 && board[i][j] < 100 {
					pos := []int{i, j}
					pieces = append(pieces, pos)
				}
			}
		}
	} else {
		//white pieces
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] < 10 && board[i][j] > 0 {
					pos := []int{i, j}
					pieces = append(pieces, pos)
				}
			}
		}
	}
	return pieces
}

//returns the coordinates of all player pieces in a 2d array in [[i,j],[i,j]] format
func getAllPlayerPieces(board [8][8]int, playerColor int) [][]int {
	var pieces [][]int
	if playerColor == 1 {
		//player color is black
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] > 10 && board[i][j] < 100 {
					pos := []int{i, j}
					pieces = append(pieces, pos)
				}
			}
		}
	} else {
		//player color is white
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if board[i][j] < 10 && board[i][j] > 0 {
					pos := []int{i, j}
					pieces = append(pieces, pos)
				}
			}
		}
	}
	return pieces
}

//changes the board passed in, moves piece at [i1][j1] to [i2][j2]
func changeBoard(board [8][8]int, i1, j1, i2, j2 int) [8][8]int {
	newBoard := board
	newBoard[i2][j2] = newBoard[i1][j1]
	newBoard[i1][j1] = 0
	return newBoard
}

//Decides which move to play
func decideWhichMove(curNode *moveGen, g *Game, alpha, beta int) int {
	//returning 0 for some reason
	if len(curNode.children) == 0 {
		//return value of this option
		return (getCurrentBoardValueAi(curNode.thisState, g) - getCurrentBoardValuePlayer(curNode.thisState, g))
	}
	if curNode.aiTurn {
		//beta, minimizing player
		//whenever max score is assured of becomes less than minimum score of alpha player is assured of (beta < alpha), the alpha player need not consider further descendants of this node
		value := 9999
		for _, child := range curNode.children {
			otherVal := decideWhichMove(child, g, alpha, beta)
			if otherVal < value {
				value = otherVal
			}
			if otherVal < beta {
				beta = otherVal
			}
			if alpha >= beta {
				break
			}
		}
		return value
	}
	//alpha, maximizing player
	value := -9999
	for _, child := range curNode.children {
		otherVal := decideWhichMove(child, g, alpha, beta)
		if otherVal > value {
			value = otherVal
		}
		if otherVal > alpha {
			alpha = otherVal
		}
		if beta <= alpha {
			break
		}
	}
	return value

}

//calculates the best move based on the current board position
func aiThink(g *Game) {
	aiThought = false
	isAiTurn = false
	aiMoves = &moveGen{thisState: g.boardState, aiTurn: isAiTurn, children: make([]*moveGen, 0)}
	queue := make([]*moveGen, 0)
	queue = append(queue, aiMoves)
	newQueue := make([]*moveGen, 0)
	for i := 0; i < 4; i++ {
		//depth 4
		isAiTurn = !isAiTurn
		for _, eachChild := range queue {
			aiPieces := getAllPieces(eachChild.thisState, g.playerColor)
			for _, piece := range aiPieces {
				eachMove := getAllMovesPiece(eachChild.thisState, g.playerColor, piece[0], piece[1], g)
				for _, move := range eachMove {
					newMove := &moveGen{thisState: changeBoard(eachChild.thisState, piece[0], piece[1], move[0], move[1]), aiTurn: isAiTurn, children: make([]*moveGen, 0)}
					eachChild.children = append(eachChild.children, newMove)
					newQueue = append(newQueue, newMove)
				}
			}
			if i == 0 {
				aiMoves = eachChild
			}
		}
		queue = newQueue
		newQueue = make([]*moveGen, 0)
	}
	//when ai makes a move, choose highest option, when player makes a move assume worst option for ai
	aiThought = true
}

//initializes the bot
func aiInit(g *Game) {
	if g.playerColor == 1 {
		isAiTurn = true
	} else {
		isAiTurn = false
	}
	aiMoves = &moveGen{thisState: g.boardState, aiTurn: isAiTurn, children: make([]*moveGen, 0)}
	depthCalc = 0
}
