package main

import(
"net/http"
"log"
"html/template"
"github.com/google/uuid"
"sync"
"fmt"
)

type Game struct {
	pieces []Piece
	board [][]BoardSquare
	gameId uuid.UUID
}

type GamesManager struct {
	games map[uuid.UUID]Game
	mu sync.Mutex
}

var templates = template.Must(template.ParseGlob("./templates/*.html"))

var pieceTemplate string =
`
<div id="piece-%d"class="w-8 h-16 bg-white rounded-md 
transform transition-transform duration-300 ease-in-out p-1 rotate-0 flex flex-col overflow-hidden shadow-lg">
  
  <div class="flex-1 grid grid-cols-3 grid-rows-3 place-items-center">
    %s
  </div>

  <div class="h-[2px] bg-white/5 mx-auto"></div>

  <div class="flex-1 grid grid-cols-3 grid-rows-3 place-items-center">
    %s
  </div>

</div>
`

var puzzleBoardTemplate string = 
`
<div id="puzzleBoard" class="w-full h-full gap-2 bg-white/5 rounded-xl p-2 self-center grid grid-cols-%d grid-rows-%d">
%s
</div>
`

var boardSquareTemplate string = 
`
<div id="%d" class="p-4 w-full h-full place-self-center rounded-md bg-gray-300">
	<div name="modifier" value=%s></div>
</div>
`

var boardSquareEmptyTemplate string = 
`
<div id="%d" class="p-4 w-full h-full rounded-sm bg-none">
</div>
`

var gameTemplate string = 
`
<div id="game" class="min-h-screen bg-gradient-to-b from-slate-900 to-black flex flex-col items-center justify-center text-white">
	<div id="board" class="w-[%s] h-[%s] rounded-sm mx-auto flex mb-12">
		%s
	</div>
	<div id="pieces" class="w-90 h-full mx-auto gap-8 flex">
		%s
	</div>
</div>
`

var whiteDot string = `<span class="w-1 h-1 bg-none rounded-full"></span>`

var blackDot string = `<span class="w-1 h-1 bg-black rounded-full"></span>`

var zeroValue = []int{0,0,0,0,0,0,0,0,0} 

var oneValue = []int{0,0,0,0,1,0,0,0,0}

var twoValue = []int{0,0,1,0,0,0,1,0,0}

var threeValue = []int{0,0,1,0,1,0,1,0,0}

var fourValue = []int{1,0,1,0,0,0,1,0,1} 

var fiveValue = []int{1,0,1,0,1,0,1,0,1}  

var sixValue = []int{1,0,1,1,0,1,1,0,1} 


func main(){

	gm := &GamesManager{
		games: make(map[uuid.UUID]Game),
	}

	createGameHandler := func(w http.ResponseWriter, r *http.Request){
		data := r.URL.Query()
		difficulty := data.Get("difficulty")
		if difficulty == "small" || difficulty == "medium" || difficulty == "large" {
			pieces, board := createGame(difficulty)
			id := uuid.New()
			gm.mu.Lock()
			var newGame  Game
			newGame.pieces = pieces
			newGame.board = board
			newGame.gameId = id
			gm.games[id] = newGame
			gm.mu.Unlock()

			piecesTemplate := createPiecesTemplate(pieces)
			boardTemplate := createBoardTemplate(board)
			size := ""
			if difficulty == "small"{
				size = "20rem"
			} else if difficulty == "medium"{
				size = "25rem"
			} else {
				size = "30rem"
			}

			fmt.Fprintf(w,gameTemplate, size, size, boardTemplate, piecesTemplate)
		} else {
			// return error
			http.Error(w, "something went wrong", http.StatusNotFound)
		}
	}

	homePageHandler := func(w http.ResponseWriter, r *http.Request){
		templates.ExecuteTemplate(w, "homepage.html", nil)
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request){
		templates.ExecuteTemplate(w, "index.html", nil)
	}

	http.HandleFunc("/create-game", createGameHandler)
	http.HandleFunc("/home-page", homePageHandler)
	http.HandleFunc("/", indexHandler)

	
	fmt.Println("Server started listning to port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getDots(dots int)([]int){
	switch dots {
	case 0:
		return zeroValue
	case 1:
		return oneValue
	case 2:
		return twoValue
	case 3:
		return threeValue
	case 4:
		return fourValue
	case 5:
		return fiveValue
	case 6:
		return sixValue
	}
	return nil
}

func createPiecesTemplate(pieces []Piece) (pieceString string){
	for _,piece := range pieces{
		var dots1Template string
		dots1 := getDots(piece.s1)
		for i := 0; i < len(dots1); i ++ {
			if dots1[i] == 1 {
				dots1Template += blackDot
			} else {
				dots1Template += whiteDot
			}
		}

		var dots2Template string
		dots2 := getDots(piece.s2)
		for i := 0; i < len(dots2); i ++ {
			if dots2[i] == 1 {
				dots2Template += blackDot
			} else {
				dots2Template += whiteDot
			}
		}

		piece := fmt.Sprintf(pieceTemplate, piece.id, dots1Template, dots2Template)
		pieceString += piece
	}
	return pieceString
}


func createBoardTemplate(board [][]BoardSquare) (boardString string){
		var boardGridTemplate string
		for i := 0; i < len(board); i++ {
			for j := 0; j < len(board); j++ {
				boardSquare := ""
				numberOfPieces := 28
				index := i*(numberOfPieces-1) - (i*(i-1))/2 + (j-i-1)
				if board[j][i].filled == true {
					boardSquare = fmt.Sprintf(boardSquareTemplate, index, board[j][i].modifier)
				} else {
					boardSquare = fmt.Sprintf(boardSquareEmptyTemplate, index)
				}
				boardGridTemplate += boardSquare
			}
		}

		boardString = fmt.Sprintf(puzzleBoardTemplate, len(board), len(board), boardGridTemplate)
		return boardString
}
