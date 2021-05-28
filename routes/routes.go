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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Default rows used for creating a new game.
	DefaultRows = 15
	// Default columns used for creating a new game.
	DefaultColumns = 10
	// Default mines used for creating a new game.
	DefaultMines = 30
)

// Instance variable to store the Database name.
var Database string

// Instance variable to store the DB Collection name.
var Collection string

// Add handlers and interceptors to the endpoints.
func Routes() *http.ServeMux {
	Database = os.Getenv("DATABASE")
	Collection = os.Getenv("MONGO_COLLECTION")

	// Interceptor chain for attaching to the requests.
	chain := net.MiddlewareChain{
		net.ApplicationJsonInterceptor(),
		net.AuthenticationInterceptor(),
	}

	// Interceptor chain with only GET method.
	getChain := chain.Add(net.CorsInterceptor(http.MethodGet))
	// Interceptor chain with only POST method.
	postChain := chain.Add(net.CorsInterceptor(http.MethodPost))

	router := http.NewServeMux()
	// Endpoint for creating a new game.
	router.HandleFunc("/game/new", getChain.Handler(HandlerNewGame))
	// Endpoint for saving a game.
	router.HandleFunc("/game/save", postChain.Handler(HandlerSaveGame))
	// Endpoint for getting the leaderboard.
	router.HandleFunc("/games", getChain.Handler(HandlerGetUsers))
	return router
}

// Handler for getting the leaderboard based on configuration.
// Creates a list of users who has games saved for the selected configuration
// and returns the list in decreasing order of rating.
func HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var games []objects.Game
	// Extract the configuration from query parameter.
	config := net.ParamString(values, "config", "Expert")
	// Print the configuration string for the logs.
	fmt.Println(config)
	// Find the cursor for the games associated with the configuration.
	cursor, err := mongo.Find(Database, Collection, bson.M{"config.name": config}, &options.FindOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Decode the games from the returned cursor.
	for cursor.Next(context.Background()) {
		var game objects.Game
		cursor.Decode(&game)
		// Print the decoded game for the logs.
		// fmt.Println(game)
		games = append(games, game)
	}
	// Usernames associated with the games.
	usernames := GetUsernamesFromGames(games)
	// Print the usernames for logs.
	// fmt.Println(usernames)
	// Find the users from the usernames.
	response, _ := GetUsers(usernames, r.Header)
	// Print the users for logs.
	// fmt.Println(response)
	// Returns the users as a json response.
	json.NewEncoder(w).Encode(response)
}

// Creates a list of usernames from the list of games.
func GetUsernamesFromGames(games []objects.Game) []string {
	// Create a map of usernames from the games.
	usersMap := make(map[string]bool)
	for _, v := range games {
		if v.Player.Username != "" {
			usersMap[v.Player.Username] = true
		}
	}
	// Convert the map of usernames to a list.
	usernames := make([]string, 0, len(usersMap))
	for k := range usersMap {
		usernames = append(usernames, k)
	}
	return usernames
}

// Handler for saving a game to the kafka topic.
// The consumer of the topic saves the game to the DB.
// The consumer updates the rating of the user in the DB.
func HandlerSaveGame(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var game objects.Game
	// Decodes the game from the request body.
	err := decoder.Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Print the game for the logs.
	fmt.Println(game)
	// Authorization token to be used for passing to the Kafka Topic
	// as the consumer would be making requests to the users service.
	token := r.Header.Get(net.Authorization)
	var match objects.Request
	match.Match = game
	match.Token = token
	// Push the game along with token to the Kafka Topic.
	formInBytes, _ := json.Marshal(match)
	if err := kafka.Push(nil, formInBytes); err != nil {
		log.Panic(err)
	}
}

// Handler for creating a new game.
// Row, Column and Mine count can be passed as query parameters.
// If any of the count is not available, defaults will be used.
func HandlerNewGame(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// Print the query parameters for the logs.
	fmt.Println(values)
	// Extract the rows count from query parameter.
	rows := net.ParamInt(values, "rows", DefaultRows)
	// Extract the columns count from query parameter.
	cols := net.ParamInt(values, "columns", DefaultColumns)
	// Extract the mines count from query parameter.
	mins := net.ParamInt(values, "mines", DefaultMines)
	// Generate a new board with the row, column and mine count.
	board := mines.GenerateBoard(rows, cols, mins)
	// Print the board for the logs.
	Print2D(board)
	// Returns the board as a json response.
	json.NewEncoder(w).Encode(board)
}

// Prints a 2d array as a matrix.
func Print2D(arr [][]string) {
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}
