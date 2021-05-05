package main

import (
	"log"
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
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	mongoClientId := os.Getenv("MONGO_CLIENT_ID")
	mongo.ConfigureMongoClient(mongoClientId)

	kafkaTopic = os.Getenv("KAFKA_TOPIC")
	kafkaBrokerId = os.Getenv("KAFKA_BROKER_ID")
	kafka.LoadPublisher(kafkaBrokerId, kafkaTopic, "")

	routes := routes.Routes()
	handler := http.HandlerFunc(routes.ServeHTTP)
	http.ListenAndServe(":5000", handler)
}
