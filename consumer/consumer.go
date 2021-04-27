package main

import (
	"github.com/atulanand206/acquisition/cmd/kafka"
	mg "github.com/atulanand206/acquisition/cmd/mongo"
	"github.com/atulanand206/minesweeper/objects"
	"github.com/atulanand206/minesweeper/mongo"
	"encoding/json"
	"fmt"
)

const (
	database              = "minesweeper"
	collectionInformation = "games"
	kafkaTopic            = "games"
	kafkaBrokerId         = "localhost:29092"
)

func main() {
	kafka.LoadConsumer(kafkaBrokerId, kafkaTopic)
	kafka.Read(func(val string) {
		game := objects.Game{}
		json.Unmarshal([]byte(val), &game)
		fmt.Println(game)
		document, _ := mongo.Document(&game)
		response := mg.Write(database, collectionInformation, *document)
		fmt.Println(response)
	})
}
