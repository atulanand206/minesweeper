package main

import (
	"github.com/atulanand206/acquisition/cmd/kafka"
	"fmt"
	"github.com/atulanand206/minesweeper/objects"
	"encoding/json"
)

func main() {
	kafka.LoadConsumer("localhost:29092", "users")
	kafka.Read(func(val string) {
		game := objects.Game{}
		json.Unmarshal([]byte(val), &game)
		fmt.Println(game)
	})
}
