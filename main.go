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

func main(){
	pieces, board := createGame("small")


	fmt.Println(pieces)

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			fmt.Print( board[i][j].pieceId, " ")
		}
		fmt.Println("")	
	}

	fmt.Println("-----------------------------------------")
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			fmt.Print( board[i][j].pieceValue, " ")
		}
		fmt.Println("")	
	}

	fmt.Println("-----------------------------------------")
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			fmt.Print( board[i][j].modifier, " ")
		}
		fmt.Println("")	
	}

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

	//board = addConditions(board)

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

	fmt.Println(pieceMap)

	fmt.Println("initial board")

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			fmt.Print( board[j][i].pieceId, " ")
		}
		fmt.Println("")	
	}

	//place pieces
	for true{
		fmt.Println("start of loop")
		fmt.Println("current: ",current.id, ", s1", current.s1, ",s2", current.s2, ", x1", current.x1,", y1", current.y1, ", x2", current.x2, ", y2", current.y2)
		var vp1, vp2 []int
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


		// we need to find a valid place for a piece then we pick, re-write check position
		vp1 = checkPositions(x, y, board)
		if result == false {
			vp1 = findNewSpot()
		}

		var piece Piece
		if len(pieceMap[s1]) > 0 {
			piece = pieceMap[s1][0]
			fmt.Println(piece, s1)
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
				fmt.Println(pieces, value, len(pieces))

				var remove int
				for index, p := range pieces{
					if p.id == piece.id{
						remove = index
					}
				}
				if len(pieces) == 1 {
					pieceMap[value] = nil
				} else if remove == len(pieces) - 1{
					pieceMap[value] = pieces[0:remove]
				} else if remove == 0 {
					pieceMap[value] = pieces[1:]	
				} else {
					pieceMap[value] = append(pieces[0:remove],pieces[remove+1:]...)
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
							pieceMap[value] = pieces[1:]
						} else if remove == len(list) - 1 {
							pieceMap[value] = pieces[0:remove]
						} else {
							pieceMap[value] = append(pieces[0:remove],pieces[remove+1:]...)
						}
					}
					break
				}
			}
		}

		board[vp1[1]][vp1[0]].pieceId = piece.id
		board[vp1[1]][vp1[0]].pieceValue = piece.s1
		board[vp1[1]][vp1[0]].filled = true
		
		current.s1 = piece.s1
		current.x1 = vp1[0]
		current.y1 = vp1[1]

		vp2, result = checkPositions(vp1[0],vp1[1], board)

		fmt.Println("after check2")

		var pv int
		if piece.s1 == s1 {
			pv = piece.s2
		}else{
			pv  = piece.s1
		}

		board[vp2[1]][vp2[0]].pieceId = piece.id
		board[vp2[1]][vp2[0]].pieceValue = pv
		board[vp2[1]][vp2[0]].filled = true

		current.s2 = piece.s2
		current.x2 = vp2[0]
		current.y2 = vp2[1]
		current.id = piece.id

		fmt.Println("current: ",current.id, ", s1", current.s1, ",s2", current.s2, ", x1", current.x1,", y1", current.y1, ", x2", current.x2, ", y2", current.y2)
		for i := 0; i < len(board); i++ {
			for j := 0; j < len(board); j++ {
				fmt.Print( board[i][j].pieceId, " ")
			}
			fmt.Println("")	
		}

		var count int
		for _, pieces := range pieceMap{
			if len(pieces) == 0 {
				count += 1 
			}
		}

		if count == len(pieceMap){
			break
		}
		fmt.Println("end of loop") 
	}

	return board

}

func checkPositions(x int, y int, board [][]BoardSquare) (validPosition []int, result bool){
	var validPositions [][]int
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

	rand.Shuffle(len(validPositions), func(i, j int){
		validPositions[i], validPositions[j] = validPositions[j], validPositions[i]
	})
	if len(validPositions) > 0 {
		validPosition = validPositions[0]
		result = true
	} else {
		validPosition = []
		result = false
	}
	return validPosition, result
}

func addConditions(board [][]BoardSquare) ([][]BoardSquare){
	var x, y int
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board); j++ {
			if board[i][j].filled == true{
				x = i
				y = j
				break
			}
		}
	}

	var modless [][]int

	// first pass is to create equal relationships
	var currentX, currentY = x, y
	initial := 0.9
	for true {
		if rand.Float64() < initial {
			pieces := checkForEqual(currentX, currentY, board, modless)

			mod := fmt.Sprintf("= %d", board[currentX][currentY].pieceValue)
			board[currentX][currentY].modifier = mod
			for _, piece := range pieces {
				if board[piece[0]][piece[1]].modifier == "" {
					board[piece[0]][piece[0]]. modifier = mod
				}
			}

			initial -= 0.2
		} 
		// we need a break or exit condition
	}



	return board
}

func checkForEqual(x, y int, board[][]BoardSquare, modless [][]int) (piecePositions [][]int){
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
