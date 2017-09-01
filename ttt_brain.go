package ttt_brain

import (
	"fmt"

	"github.com/genzi/brain"
)

type TicTacToe struct {
	board   []int
	myBrain *brain.Brain
}

func New() *TicTacToe {
	game := new(TicTacToe)
	game.board = make([]int, 9)
	for i := 0; i < 9; i++ {
		game.board[i] = 0
	}
	game.myBrain = brain.New(9, 36, 9)
	return game
}

//not finished yet
func (ttt *TicTacToe) Train() {
	boardBeforeAImove := make([]int, 9)
	copy(boardBeforeAImove, ttt.board)
	ttt.AiMove()
	position := checkDiff(boardBeforeAImove, ttt.board, 9)

	fmt.Println("Position: ", position)
	result := ttt.myBrain.Process(convertSliceIntToFloat64(boardBeforeAImove))
	fmt.Println(result)

	var pattern [][][]float64

	pattern = append(pattern, [][]float64{convertSliceIntToFloat64([]int{1, -1, 0, 0, 0, 0, 0, 0, 0}), convertSliceIntToFloat64([]int{0, 0, 0, 1, 0, 0, 0, 0, 0})})
	pattern = append(pattern, [][]float64{convertSliceIntToFloat64([]int{1, -1, 1, -1, 0, 0, 0, 0, -1}), convertSliceIntToFloat64([]int{0, 0, 0, 0, 0, 0, 0, 1, 0})})
	pattern = append(pattern, [][]float64{convertSliceIntToFloat64([]int{1, -1, 1, -1, 0, 0, 0, 0, 0}), convertSliceIntToFloat64([]int{0, 0, 0, 0, 0, 0, 0, 0, 1})})

	fmt.Println("Tic tac toe Brain is trainning...")
	ttt.myBrain.Train(pattern, 10000, 0.6, 0.5)
	fmt.Println("Result after trainning:")
	result = ttt.myBrain.Process(convertSliceIntToFloat64(boardBeforeAImove))
	fmt.Println(result)
	result = ttt.myBrain.Process(convertSliceIntToFloat64([]int{1, -1, 1, -1, 0, 0, 0, 0, -1}))
	fmt.Println(result)
	result = ttt.myBrain.Process(convertSliceIntToFloat64([]int{1, -1, 1, -1, 0, 0, 0, 0, 0}))
	fmt.Println(result)
}

//for cli game
func (ttt *TicTacToe) RunGame() {
	fmt.Printf("Computer: O, You: X\nPlay (1)st or (2)nd? ")
	player := 0
	fmt.Scanln(&player)

	for turn := 0; turn < 9 && ttt.Win() == 0; turn++ {
		if (turn+player)%2 == 0 {
			ttt.AiMove()
		} else {
			ttt.draw()
			for {
				var move int
				fmt.Printf("\nInput move ([0..8]): ")
				fmt.Scanln(&move)
				if ttt.HumanMove(move) == true {
					break
				}
			}
		}
	}

	switch ttt.Win() {
	case 0:
		fmt.Println("A draw.")
	case 1:
		fmt.Println("You lose.")
	case -1:
		fmt.Println("You win.")
	}
}

func (ttt *TicTacToe) GetBoard() []int {
	return ttt.board
}

func (ttt *TicTacToe) RunAIvsAI() {
	fmt.Println("AI vs AI")
	player := -1

	for turn := 0; turn < 9 && ttt.Win() == 0; turn++ {
		if (turn+player)%2 == 0 {
			ttt.AiMove(1)
			ttt.draw()
		} else {
			ttt.AiMove(-1)
			ttt.draw()
		}
	}

	switch ttt.Win() {
	case 0:
		fmt.Println("A draw.")
	case 1:
		fmt.Println("You lose.")
	case -1:
		fmt.Println("You win.")
	}
}

func (ttt *TicTacToe) Win() int {
	wins := [8][3]uint{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}

	for i := 0; i < 8; i++ {
		if ttt.board[wins[i][0]] != 0 &&
			ttt.board[wins[i][0]] == ttt.board[wins[i][1]] &&
			ttt.board[wins[i][0]] == ttt.board[wins[i][2]] {
			return ttt.board[wins[i][2]]
		}
	}
	return 0
}

func (ttt *TicTacToe) AiMove(player ...int) {
	if len(player) == 0 {
		player = append(player, 1)
	}

	move := -1
	score := -2
	for i := 0; i < 9; i++ { //all moves
		if ttt.board[i] == 0 { //if legal
			ttt.board[i] = player[0] //Try the move // WAS 1
			thisScore := -ttt.minimax(player[0] * -1)
			ttt.board[i] = 0
			if thisScore > score {
				score = thisScore
				move = i
			}
		}
	}
	if move >= 0 && move < 9 {
		ttt.board[move] = player[0] // WAS 1
	}
}

func (ttt *TicTacToe) HumanMove(move int) bool {

	if move >= 9 || move < 0 || ttt.board[move] != 0 {
		return false
	}
	ttt.board[move] = -1
	return true
}

func (ttt *TicTacToe) minimax(player int) int {
	winner := ttt.Win()
	if winner != 0 {
		return winner * player
	}
	move := -1
	score := -2
	for i := 0; i < 9; i++ { //all moves
		if ttt.board[i] == 0 { //if legal
			ttt.board[i] = player //Try the move
			thisScore := -ttt.minimax(player * -1)
			if thisScore > score {
				score = thisScore
				move = i
			} //Pick the one that's worst for the opponent
			ttt.board[i] = 0
		}
	}
	if move == -1 {
		return 0
	}
	return score
}

func (ttt *TicTacToe) draw() {
	fmt.Printf(" %c | %c | %c\n",
		BoardStateToChar(ttt.board[0]), BoardStateToChar(ttt.board[1]), BoardStateToChar(ttt.board[2]))
	fmt.Printf("---+---+---\n")
	fmt.Printf(" %c | %c | %c\n",
		BoardStateToChar(ttt.board[3]), BoardStateToChar(ttt.board[4]), BoardStateToChar(ttt.board[5]))
	fmt.Printf("---+---+---\n")
	fmt.Printf(" %c | %c | %c\n\n",
		BoardStateToChar(ttt.board[6]), BoardStateToChar(ttt.board[7]), BoardStateToChar(ttt.board[8]))
}

func BoardStateToChar(i int) byte {
	char := byte(' ')

	switch i {
	case -1:
		char = 'X'
	case 0:
		char = ' '
	case 1:
		char = 'O'
	}
	return char
}

func convertSliceIntToFloat64(slice []int) []float64 {
	float64Slice := make([]float64, len(slice))

	for i := 0; i < len(slice); i++ {
		float64Slice[i] = float64(slice[i])
	}

	return float64Slice
}

func checkDiff(a []int, b []int, size int) int {
	i := 0
	for ; i < size; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

/*
func main() {
	ttt_game := New()
	//	ttt_game.RunGame()
	//	ttt_game.board[0] = -1
	ttt_game.board[1] = 1
	ttt_game.RunAIvsAI()
	ttt_game = New()
	ttt_game.Train()
}
*/
