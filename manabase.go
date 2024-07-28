package manabase

type Game struct {
	UUID       string
	Players    []string
	DatePlayed string
	GameNumber int
}

type Player struct {
	UUID string
	Name string
}

type Deck struct {
	Name    string
	OwnerId string
}
