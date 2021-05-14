package objects

import "time"

type (
	User struct {
		Id       string `json:"id" bson:"_id"`
		Username string `json:"username" bson:"username"`
		Name     string `json:"name" bson:"name"`
		Rating   int    `json:"rating" bson:"rating"`
	}

	Config struct {
		Row   int    `json:"row" bson:"row, omitempty"`
		Col   int    `json:"col" bson:"col, omitempty"`
		Mines int    `json:"mines" bson:"mines, omitempty"`
		Type  string `json:"name" bson:"name, omitempty"`
	}

	Range struct {
		Start time.Time `json:"start" bson:"start, omitempty"`
		End   time.Time `json:"end" bson:"end, omitempty"`
	}

	Game struct {
		Player   User    `json:"player" bson:"player, omitempty"`
		Conf     Config  `json:"config" bson:"config, omitempty"`
		Times    []Range `json:"times" bson:"times, omitempty"`
		Score    int     `json:"score" bson:"score, omitempty"`
		Won      bool    `json:"won" bson:"won, omitempty"`
		Finished bool    `json:"finished" bson:"finished, omitempty"`
	}
)
