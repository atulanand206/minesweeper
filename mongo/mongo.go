package mongo

import (
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/atulanand206/minesweeper/objects"
)

func Document(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func Game(document *mongo.SingleResult) (v objects.Game, err error) {
	var game objects.Game
	if err = document.Decode(&game); err != nil {
		log.Panic(err)
	}
	return game, err
}