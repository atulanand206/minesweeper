package objects

import "time"

type (
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
		Conf     Config  `json:"config" bson:"config, omitempty"`
		Times    []Range `json:"times" bson:"times, omitempty"`
		Score    int     `json:"score" bson:"score, omitempty"`
		Won      bool    `json:"won" bson:"won, omitempty"`
		Finished bool    `json:"finished" bson:"finished, omitempty"`
	}
)
