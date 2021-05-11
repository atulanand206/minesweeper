package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	"github.com/atulanand206/minesweeper/mines"
	"github.com/atulanand206/minesweeper/objects"
	"go.mongodb.org/mongo-driver/bson"
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
)

var database string
var collection string

func Routes() *http.ServeMux {
	database = os.Getenv("DATABASE")
	collection = os.Getenv("MONGO_COLLECTION")

	router := http.NewServeMux()
	router.HandleFunc("/game/new", http.HandlerFunc(newGameHandler))
	router.HandleFunc("/game/save", http.HandlerFunc(saveGameHandler))
	router.HandleFunc("/games", http.HandlerFunc(getGamesHandler))
	return router
}

func getGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(cors, "*")
	w.Header().Set(contentTypeKey, contentTypeApplicationJson)
	values := r.URL.Query()
	var response []objects.Game
	config := paramString(values, "config", "Expert")
	cursor, err := mongo.Find(database, collection, bson.M{"config.name": config})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for cursor.Next(context.Background()) {
		var game objects.Game
		cursor.Decode(&game)
		response = append(response, game)
	}
	json.NewEncoder(w).Encode(response)
}

func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(cors, "*")
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
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(cors, "*")
	w.Header().Set(contentTypeKey, contentTypeApplicationJson)
	values := r.URL.Query()
	rows := paramInt(values, "rows", defaultRows)
	cols := paramInt(values, "columns", defaultColumns)
	mins := paramInt(values, "mines", defaultMines)
	fmt.Println(values)
	board := mines.GenerateBoard(rows, cols, mins)
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

func print2D(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}
