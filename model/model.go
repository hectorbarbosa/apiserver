package model

type RequestActor struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type Actor struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type RequestFilm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Release     string `json:"release"`
	Rating      byte   `json:"rating"`
}

type Film struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Release     string `json:"release"`
	Rating      byte   `json:"rating"`
}
