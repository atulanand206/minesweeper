package main

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"encoding/json"
	"strconv"
)

const (
	defaultRows    = 15
	defaultColumns = 10
	defaultMines   = 30
)

const (
	contentTypeKey             = "content-type"
	cors                       = "Access-Control-Allow-Origin"
	contentTypeApplicationJson = "application/json"
	htmlTemplatePath           = "game.html"
)

func main() {
	routes := routes()
	handler := http.HandlerFunc(routes.ServeHTTP)
	http.ListenAndServe(":5000", handler)
}

func routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/game/new", http.HandlerFunc(newGameHandler))
	return router
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	rows := defaultRows
	cols := defaultColumns
	mins := defaultMines
	if values.Get("rows") != "" {
		rows, _ = strconv.Atoi(values.Get("rows"))
	}
	if values.Get("columns") != "" {
		cols, _ = strconv.Atoi(values.Get("columns"))
	}
	if values.Get("mines") != "" {
		mins, _ = strconv.Atoi(values.Get("mines"))
	}
	fmt.Println(values)
	w.Header().Set(cors, "*")
	w.Header().Set(contentTypeKey, contentTypeApplicationJson)
	board := generateBoard(rows, cols, mins)
	print2D(board)
	json.NewEncoder(w).Encode(board)
}

const (
	mine  = "*"
	empty = " "
)

func generateBoard(rows int, columns int, mines int) [][]string {
	board := emptyBoard(rows, columns)
	fmt.Println(len(board), len(board[0]))
	generateMines(board, rows, columns, mines)
	generateValues(board)
	return board
}

func generateValues(board [][]string) {
	directions := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if !isMine(board, i, j) {
				k := 0
				for _, d := range directions {
					k += checkMine(board, i+d[0], j+d[1])
				}
				if k != 0 {
					board[i][j] = string(k + 48)
				}
			}
		}
	}
}

func checkMine(board [][] string, x int, y int) int {
	if isMine(board, x, y) {
		return 1
	}
	return 0
}

func isMine(board [][] string, x int, y int) bool {
	if x < 0 || y < 0 || x >= len(board) || y >= len(board[x]) {
		return false
	}
	return board[x][y] == mine
}

func generateMines(mines [][]string, rows int, columns int, count int) {
	mp := minePositions(rows, columns, count)
	for k := range mp {
		mines[k/columns][k%columns] = mine
	}
}

func minePositions(rows int, columns int, count int) map[int]int {
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

func emptyBoard(rows int, columns int) [][]string {
	mines := make([][]string, rows)
	for i := 0; i < rows; i++ {
		mines[i] = make([]string, columns)
		for j := range mines[i] {
			mines[i][j] = empty
		}
	}
	return mines
}

func print2D(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}
