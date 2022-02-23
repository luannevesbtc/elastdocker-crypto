package model

type Passengers struct {
	TotalPassengers int `json:"totalPassengers"`
	TotalPages      int `json:"totalPages"`
	Data            []struct {
		ID      string `json:"_id"`
		Name    string `json:"name"`
		Trips   int    `json:"trips"`
		Airline []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Country     string `json:"country"`
			Logo        string `json:"logo"`
			Slogan      string `json:"slogan"`
			HeadQuaters string `json:"head_quaters"`
			Website     string `json:"website"`
			Established string `json:"established"`
		} `json:"airline"`
		V int `json:"__v"`
	} `json:"data"`
}
