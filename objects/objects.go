package objects

import "time"

type (
	Config struct {
		Row   int    `json:"row"`
		Col   int    `json:"col"`
		Mines int    `json:"mines"`
		Type  string `json:"name"`
	}

	Range struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	Game struct {
		Conf     Config  `json:"config"`
		Times    []Range `json:"times"`
		Score    int     `json:"score"`
		Won      bool    `json:"won"`
		Finished bool    `json:"finished"`
	}
)
