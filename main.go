package main

import (
	"net/http"
	"os"

	"github.com/atulanand206/go-kafka"
	"github.com/atulanand206/go-mongo"
	"github.com/atulanand206/minesweeper/routes"
	"github.com/joho/godotenv"
)

var kafkaTopic string
var kafkaBrokerId string

func main() {
	godotenv.Load()

	// Register the MongoDB cloud atlas database.
	mongoClientId := os.Getenv("MONGO_CLIENT_ID")
	mongo.ConfigureMongoClient(mongoClientId)

	// Register the Kafka cluster publisher.
	kafkaTopic = os.Getenv("KAFKA_TOPIC")
	kafkaBrokerId = os.Getenv("KAFKA_BROKER_ID")
	kafka.LoadPublisher(kafkaBrokerId, kafkaTopic)

	// Register the endpoints exposed from the service.
	routes := routes.Routes()
	handler := http.HandlerFunc(routes.ServeHTTP)
	http.ListenAndServe(":5000", handler)
}
