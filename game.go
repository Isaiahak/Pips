package main 

import(
"fmt"
"math/rand"
)

type Piece struct {
	id int
	s1 int
	s2 int
}

type BoardSquare struct {
	pieceId int
	pieceValue int
	modifier string
	filled bool 
}

type CurrentPiece struct {
	x1 int
	x2 int 
	y1 int 
	y2 int 
	s1 int 
	s2 int 
	id int 
}


func createGame(size string) ([]Piece, [][]BoardSquare){
	var pieces []Piece
	switch size {
	// 5 to 7 pieces
	case "small":
		numOfPieces := 7 - rand.Intn(2)
		pieces = generatePieces(numOfPieces)
	// 7 to 9 pieces	
	case "medium":
		numOfPieces := 9 - rand.Intn(2)
		pieces = generatePieces(numOfPieces)
	// 10 to 13 pieces
	case "large":
		numOfPieces := 13 - rand.Intn(3)
		pieces = generatePieces(numOfPieces)
	default:
		fmt.Print("invalid size!")
	}

	var board [][]BoardSquare
	switch size {
	case "small":
		board = createBoard(6, pieces)
	case "medium":
		board = createBoard(8, pieces)
	case "large":
		board = createBoard(10, pieces)
	default:
		fmt.Print("invalid size!")
	}

	board = addConditions(board)

	return pieces, board

}

func generatePieces(size int) []Piece{
	const numberOfPieces int = 28
	pieces := make([]Piece, 0, numberOfPieces*(numberOfPieces+1)/2)
	for i := 0; i <= 6; i++ {
		for j := i; j <= 6; j++ {
			index := i*(numberOfPieces-1) - (i*(i-1))/2 + (j-i-1)
			if index == 0 {
				index = 420
			}
			newPiece := Piece{
				id: index,
				s1:i,
				s2:j,
			}
			pieces = append(pieces, newPiece)
			
		}
	}

	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	shuffledPieces := pieces[:size-1] 

	return shuffledPieces
}

func createBoard(size int, pieces []Piece) [][]BoardSquare{
	corner := rand.Intn(3)

	var current CurrentPiece
	switch corner {
	case 0:
		current.x1 = 0
		current.y1 = 0
	case 1:
		current.x1 = size-1
		current.y1 = 0
	case 2:
		current.x1 = 0
		current.y1 = size-1
	case 3:	
		current.x1 = size-1
		current.y1 = size-1
	}

	board := make([][]BoardSquare, size)
	for column := range board {
		board[column] = make([]BoardSquare, size)
	}

	for j := 0; j < len(board); j ++ {
		for i := 0; i < len(board); i++ {
			board[j][i].pieceValue = -1
			board[j][i].modifier = "None"
		}
	}

	current.s1 = pieces[0].s1
	current.s2 = pieces[0].s2
	current.id  = pieces[0].id
	placedSide := rand.Intn(1)

	var firstValue, secondValue int
	if placedSide == 0{
		firstValue = current.s1
		secondValue = current.s2
	}else{
		firstValue = current.s2
		secondValue = current.s1
	}

	board[current.y1][current.x1].pieceId = current.id
	board[current.y1][current.x1].pieceValue = firstValue
	board[current.y1][current.x1].filled = true

	orientation := rand.Intn(1)
	// 0 == horizontal 1 == vertical
	var pX int
	var pY int

	if orientation == 0{
		if corner == 0{
			pX = 1
			pY = 0
		} else if corner == 1{
			pY = 0
			pX = size - 2
		} else if corner == 2{
			pX = 1
			pY = size - 1
		}else{
			pX = size - 2
			pY = size - 1
		}
	}else{
		if corner == 0{
			pX = 0
			pY = 1
		} else if corner == 1{
			pX = size - 1
			pY = 1
		} else if corner == 2{
			pX = 0
			pY = size - 2
		} else{
			pX = size - 1
			pY = size - 2
		}
	}

	board[pY][pX].pieceId = current.id 
	board[pY][pX].pieceValue = secondValue
	board[pY][pX].filled = true

	current.x2 = pX
	current.y2 = pY

	restOfPieces := pieces[1:]

	// put pieces into a map based on their values
	pieceMap := make(map[int][]Piece)
	for _, piece := range restOfPieces{
		if piece.s1 != piece.s2 {
			pieceMap[piece.s1] = append(pieceMap[piece.s1], piece)
			pieceMap[piece.s2] = append(pieceMap[piece.s2], piece)
		} else {
			pieceMap[piece.s1] = append(pieceMap[piece.s1], piece)
		}
	}

	for true{
		var placement []int

		placedSide := rand.Intn(1)
		var x, y, s1 int
		if placedSide == 0 {
			s1 = current.s1
			x = current.x1 
			y = current.y1
		} else {
			s1 = current.s2
			x = current.x2
			y = current.y2
		}

		potentialPlacements := findPotentialPlacements(x, y, board)

		// find a place for a piece
		if len(potentialPlacements) > 0 {
			rand.Shuffle(len(potentialPlacements), func(i, j int) {
				potentialPlacements[i], potentialPlacements[j] = potentialPlacements[j], potentialPlacements[i]
			})
			placement = potentialPlacements[0]
		} else {
			for i := 0; i < len(board); i ++ {
				for j := 0; j < len(board); j++ {
					if board[j][i].filled == true {
						potentialPlacements := findPotentialPlacements(i,j, board)
						if len(potentialPlacements) > 0 {
							rand.Shuffle(len(potentialPlacements), func(i, j int) {
								potentialPlacements[i], potentialPlacements[j] = potentialPlacements[j], potentialPlacements[i]
							})
							placement = potentialPlacements[0]
							break
						}
					}
				}
			}
		}

		// find a piece to put in the place
		var piece Piece
		if len(pieceMap[s1]) > 0  && rand.Float64() < 0.75{
			piece = pieceMap[s1][0]
			if len(pieceMap[s1]) == 1{
				pieceMap[s1] = nil
			} else {
				pieceMap[s1] = pieceMap[s1][1:]
			}

			var value int
			if piece.s1 != piece.s2{

				if piece.s1 == s1{
					value = piece.s2
				} else {
					value = piece.s1
				}

				pieces = pieceMap[value]

				var remove int
				for index, p := range pieces{
					if p.id == piece.id{
						remove = index
					}
				}
				if len(pieceMap[value]) == 1 {
					pieceMap[value] = nil
				} else if remove == len(pieces) - 1{
					pieceMap[value] = pieceMap[value][0:remove]
				} else if remove == 0 {
					pieceMap[value] = pieceMap[value][1:]	
				} else {
					pieceMap[value] = append(pieceMap[value][0:remove],pieceMap[value][remove+1:]...)
				}
			} 
		}else {
			for index, pieces := range pieceMap {
				if len(pieces) > 0 {
					piece = pieces[0]
					if len(pieceMap[index]) == 1 {
						pieceMap[index] = nil
					} else {
						pieceMap[index] = pieces[1:]
					}

					if piece.s1 != piece.s2 {
						var value int
						if piece.s1 == index{
							value = piece.s2
						} else {
							value = piece.s1
						}

						list := pieceMap[value]
						var remove int 
						for index, p := range list{
							if piece.id == p.id{
								remove = index
							}
						}
						if len(list) == 1 {
							pieceMap[value] = nil
						} else if remove == 0 {
							pieceMap[value] = pieceMap[value][1:]
						} else if remove == len(list) - 1 {
							pieceMap[value] = pieceMap[value][0:remove]
						} else {
							pieceMap[value] = append(pieceMap[value][0:remove],pieceMap[value][remove+1:]...)
						}
					}
					break
				}
			}
		}

		board[placement[1]][placement[0]].pieceId = piece.id
		board[placement[1]][placement[0]].pieceValue = piece.s1
		board[placement[1]][placement[0]].filled = true
		
		current.s1 = piece.s1
		current.x1 = placement[0]
		current.y1 = placement[1]

		var pv int
		if piece.s1 == s1 {
			pv = piece.s2
		}else{
			pv  = piece.s1
		}

		board[placement[3]][placement[2]].pieceId = piece.id
		board[placement[3]][placement[2]].pieceValue = pv
		board[placement[3]][placement[2]].filled = true

		current.s2 = piece.s2
		current.x2 = placement[2]
		current.y2 = placement[3]
		current.id = piece.id

		var count int
		for _, pieces := range pieceMap{
			if len(pieces) == 0 {
				count += 1 
			}
		}

		if count == len(pieceMap){
			break
		}
	}

	return board

}

func checkPositions(x int, y int, board [][]BoardSquare) (validPositions [][]int){
	if x + 1 < len(board){
		if board[y][x+1].filled == false{
			position := []int{x+1,y}
			validPositions = append(validPositions, position)
		}
	}
	if y + 1 < len(board){
		if board[y+1][x].filled == false{
			position := []int{x,y+1}
			validPositions = append(validPositions, position)
		}

	}
	if y - 1 >= 0 {
		if board[y-1][x].filled == false{
			position := []int{x,y-1}
			validPositions = append(validPositions, position)
		}
	}
	if x - 1 >= 0 {
		if board[y][x-1].filled == false{
			position := []int{x-1,y}
			validPositions = append(validPositions, position)
		}
	}
	return validPositions

}

func findPotentialPlacements(x,y int, board [][]BoardSquare) (potentialPlacements [][]int){
	positions := checkPositions(x, y, board)
	for _, position := range positions{
		orientations := checkPositions(position[0],position[1], board)
		if len(orientations) > 0 {
			for _,orientation := range orientations{
				placement := []int{position[0],position[1],orientation[0],orientation[1]}
				potentialPlacements = append(potentialPlacements, placement)
			}
		}
	}
	return potentialPlacements

}

func addConditions(board [][]BoardSquare) ([][]BoardSquare){
	// first pass is to create equal relationships

	visited := make(map[int]bool)
	for i := 0; i < len(board); i ++ {
		for j := 0; j < len(board); j++ {
			if board[j][i].filled == true && visited[board[j][i].pieceId] == false{
				pieces := checkForEqual(j,i,board)
				mod := fmt.Sprintf("= %d", board[j][i].pieceValue )
				if len(pieces) > 0 {
					for _, piece := range pieces {
						if board[piece[0]][piece[1]].modifier == "" {
							board[piece[0]][piece[1]].modifier = mod
							visited[board[piece[0]][piece[1]].pieceId] = true
						}
					}
				}	
				board[j][i].modifier = mod
				visited[board[j][i].pieceId] = true

			}
		}
	}



	return board
}

func checkForEqual(x, y int, board[][]BoardSquare) (piecePositions [][]int){
	if x + 1 < len(board) - 1 {
		if board[y][x+1].filled == true && board[x+1][y].pieceValue == board[x][y].pieceValue{
			piece := []int{x+1, y}
			piecePositions = append(piecePositions, piece)
		} 
	}

	if x - 1 >= 0 {
		if board[y][x-1].filled == true && board[x-1][y].pieceValue == board[x][y].pieceValue{
			piece := []int{x-1, y}
			piecePositions = append(piecePositions, piece)
		}
	}

	if y + 1 < len(board) - 1 {
		if board[y+1][x].filled == true && board[x][y+1].pieceValue == board[x][y].pieceValue{
			piece := []int{x, y+1}
			piecePositions = append(piecePositions, piece)
		}
	}

	if y - 1 >= 0 {
		if board[y-1][x].filled == true && board[x][y-1].pieceValue == board[x][y].pieceValue{
			piece := []int{x, y-1}
			piecePositions = append(piecePositions, piece)
		}
	}
	return piecePositions

}
