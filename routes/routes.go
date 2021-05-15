package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	net "github.com/atulanand206/go-network"
	"github.com/atulanand206/minesweeper/mines"
	"github.com/atulanand206/minesweeper/objects"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	defaultRows    = 15
	defaultColumns = 10
	defaultMines   = 30
)

var database string
var collection string

func Routes() *http.ServeMux {
	database = os.Getenv("DATABASE")
	collection = os.Getenv("MONGO_COLLECTION")

	chain := net.MiddlewareChain{
		net.ApplicationJsonInterceptor(),
		// net.AuthenticationInterceptor(),
	}

	getChain := chain.Add(net.CorsInterceptor(http.MethodGet))
	postChain := chain.Add(net.CorsInterceptor(http.MethodPost))

	router := http.NewServeMux()
	router.HandleFunc("/game/new", getChain.Handler(newGameHandler))
	router.HandleFunc("/game/save", postChain.Handler(saveGameHandler))
	router.HandleFunc("/games", getChain.Handler(getUsersHandler))
	return router
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var games []objects.Game
	config := net.ParamString(values, "config", "Expert")
	cursor, err := mongo.Find(database, collection, bson.M{"config.name": config})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for cursor.Next(context.Background()) {
		var game objects.Game
		cursor.Decode(&game)
		games = append(games, game)
	}
	usernames := getUsernamesFromGames(games)
	response, _ := GetUsers(usernames, r.Header)
	json.NewEncoder(w).Encode(response)
}

func getUsernamesFromGames(games []objects.Game) []string {
	usersMap := make(map[string]bool)
	for _, v := range games {
		if v.Player.Username != "" {
			usersMap[v.Player.Username] = true
		}
	}
	usernames := make([]string, 0, len(usersMap))
	for k := range usersMap {
		usernames = append(usernames, k)
	}
	return usernames
}

func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var game objects.Game
	err := decoder.Decode(&game)
	if err == nil {
		fmt.Println(game)
	}
	token := r.Header.Get(net.Authorization)
	var match objects.Request
	match.Match = game
	match.Token = token
	formInBytes, _ := json.Marshal(match)
	if err := kafka.Push(nil, formInBytes); err != nil {
		log.Panic(err)
	}
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	rows := net.ParamInt(values, "rows", defaultRows)
	cols := net.ParamInt(values, "columns", defaultColumns)
	mins := net.ParamInt(values, "mines", defaultMines)
	fmt.Println(values)
	board := mines.GenerateBoard(rows, cols, mins)
	print2D(board)
	json.NewEncoder(w).Encode(board)
}

func print2D(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}
