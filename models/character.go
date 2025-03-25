package models

type Character struct {
	Name  string `json:"name"`
	Race  string `json:"race"`
	Stats Stats  `json:"stats"`
	Image string `json:"image"`
}

type Stats struct {
	Vitality     int `json:"vitality"`
	Endurance    int `json:"endurance"`
	Mind         int `json:"mind"`
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
	Faith        int `json:"faith"`
}
