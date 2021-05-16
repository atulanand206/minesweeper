package mines

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	// String to be marked for the mined cell.
	Mine = "*"
	// String to be marked for the Empty cell.
	Empty = " "
)

// Generate board based on the input count configuration.
func GenerateBoard(rows int, columns int, mines int) [][]string {
	board := EmptyBoard(rows, columns)
	fmt.Println(len(board), len(board[0]))
	GenerateMines(board, rows, columns, mines)
	GenerateValues(board)
	return board
}

// Generate values on the board.
func GenerateValues(board [][]string) {
	directions := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if !IsMine(board, i, j) {
				k := 0
				for _, d := range directions {
					k += CheckMine(board, i+d[0], j+d[1])
				}
				if k != 0 {
					board[i][j] = strconv.Itoa(k)
				}
			}
		}
	}
}

// Returns 1 if the coordinate has a mine.
func CheckMine(board [][]string, x int, y int) int {
	if IsMine(board, x, y) {
		return 1
	}
	return 0
}

// Returns true if the coordinate has a mine.
func IsMine(board [][]string, x int, y int) bool {
	if x < 0 || y < 0 || x >= len(board) || y >= len(board[x]) {
		return false
	}
	return board[x][y] == Mine
}

// Generates mines from the randomized mine positions.
func GenerateMines(mines [][]string, rows int, columns int, count int) {
	mp := MinePositions(rows, columns, count)
	for k := range mp {
		mines[k/columns][k%columns] = Mine
	}
}

// Generate mines after generating random positions.
func MinePositions(rows int, columns int, count int) map[int]int {
	total := rows * columns
	mp := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	for {
		if len(mp) < count {
			mp[rand.Intn(total)] = 1
		} else {
			break
		}
	}
	return mp
}

// Generate an empty board from rows and columns count.
func EmptyBoard(rows int, columns int) [][]string {
	mines := make([][]string, rows)
	for i := 0; i < rows; i++ {
		mines[i] = make([]string, columns)
		for j := range mines[i] {
			mines[i][j] = Empty
		}
	}
	return mines
}
