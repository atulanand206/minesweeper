package mongo

import (
	"log"

	"github.com/atulanand206/minesweeper/objects"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Converts a struct object instance to a mongo document.
func Document(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// Converts a mongo single result to a game instance.
func Game(document *mongo.SingleResult) (v objects.Game, err error) {
	var game objects.Game
	if err = document.Decode(&game); err != nil {
		log.Panic(err)
	}
	return game, err
}
