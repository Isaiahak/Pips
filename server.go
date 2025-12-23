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

var GamesManager struct {
	games map[uuid.UUID]Game
	mu sync.Mutex
}

var templates = template.Must(template.ParseGlob("./templates/*.html"))

var pieceTemplate string =
`
<div class="w-20 h-40 bg-white rounded-xl border-2 border-purple-500 flex flex-col overflow-hidden shadow-lg">
  
  <div class="flex-1 grid grid-cols-3 grid-rows-3 place-items-center">
    %s
  </div>

  <div class="h-[2px] bg-purple-500"></div>

  <div class="flex-1 grid grid-cols-3 grid-rows-3 place-items-center">
    %s
  </div>

</div>
`

var whiteDot string = `<span class="w-3 h-3 bg-none rounded-full"></span>`

var blackDot string = `<span class="w-3 h-3 bg-none rounded-full"></span>`

var zeroValue = []int{0,0,0,0,0,0,0,0,0} 

var oneValue = []int{0,0,0,0,1,0,0,0,0}

var twoValue = []int{0,0,1,0,0,0,1,0,0}

var threeValue = []int{0,0,1,0,1,0,1,0,0}

var fourValue = []int{1,0,1,0,0,0,1,0,1} 

var fiveValue = []int{1,0,1,0,1,0,1,0,1}  

var sixValue = []int{1,0,1,1,0,1,1,0,1} 


func main(){

	createGameHandler := func(w http.ResponseWriter, r *http.Request){
		data := r.URL.Query()
		difficulty := data.Get("difficulty")
		if difficulty != "small" || difficulty != "medium" || difficulty != "large" {
			// return error

		} else {
			pieces, board := createGame(difficulty)
			id := uuid.New()
			GamesManager.mu.Lock()
			var newGame  Game
			newGame.pieces = pieces
			newGame.board = board
			newGame.gameId = id
			GamesManager.games[id] = newGame

			
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

				fmt.Fprintf(w, pieceTemplate, dots1Template, dots2Template)
			}


		}
		templates.ExecuteTemplate(w, "gamepage.html", nil)
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


