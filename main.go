package main

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/atulanand206/acquisition/cmd/kafka"
	"log"
	"github.com/atulanand206/minesweeper/objects"
	"github.com/atulanand206/acquisition/cmd/mongo"
	"net/url"
	"go.mongodb.org/mongo-driver/bson"
	"context"
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
	kafkaTopic                 = "games"
	kafkaBrokerId              = "localhost:19092"
	database                   = "minesweeper"
	collection                 = "games"
	mongoClientId              = "mongodb+srv://Anand:deathByChance@cluster0.x2q5g.mongodb.net/"
)

func main() {
	mongo.ConfigureMongoClient(mongoClientId)
	kafka.LoadPublisher(kafkaBrokerId, kafkaTopic, "0.0.0.0:9000")
	routes := routes()
	handler := http.HandlerFunc(routes.ServeHTTP)
	http.ListenAndServe(":5000", handler)
}

func routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/game/new", http.HandlerFunc(newGameHandler))
	router.HandleFunc("/game/save", http.HandlerFunc(saveGameHandler))
	router.HandleFunc("/games", http.HandlerFunc(getGamesHandler))
	return router
}

func getGamesHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var response []objects.Game
	config:= paramString(values, "config", "Expert")
	cursor := mongo.Find(database, collection, bson.M{ "config.name": config })
	for cursor.Next(context.Background()) {
		var game objects.Game
		cursor.Decode(&game)
		response = append(response, game)
	}
	w.Header().Set(cors, "*")
	w.Header().Set(contentTypeKey, contentTypeApplicationJson)
	json.NewEncoder(w).Encode(response)
}

func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ob objects.Game
	err := decoder.Decode(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	formInBytes, _ := json.Marshal(ob)
	if err := kafka.Push(nil, formInBytes); err != nil {
		log.Panic(err)
	}
	w.Header().Set(cors, "*")
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	rows := paramInt(values, "rows", defaultRows)
	cols := paramInt(values, "columns", defaultColumns)
	mins := paramInt(values, "mines", defaultMines)
	fmt.Println(values)
	w.Header().Set(cors, "*")
	w.Header().Set(contentTypeKey, contentTypeApplicationJson)
	board := generateBoard(rows, cols, mins)
	print2D(board)
	json.NewEncoder(w).Encode(board)
}

func paramInt(values url.Values, key string, def int) int {
	rows := def
	if values.Get(key) != "" {
		rows, _ = strconv.Atoi(values.Get(key))
	}
	return rows
}

func paramString(values url.Values, key string, def string) string {
	rows := def
	if values.Get(key) != "" {
		rows = values.Get(key)
	}
	return rows
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
