package model

type Score struct {
	Id         string
	PlayerId   string
	PlayerName string
	Points     int
	Birdies    int
	Eagles     int
	Muligans   int
	Season     int
	Day        string
}

type ScoreInput struct {
	PlayerId string
	Points   int
	Birdies  int
	Eagles   int
	Muligans int
	Season   int
}

type Scoreboard struct {
	Players []ScoreboardPlayer
	Season  int
}

type ScoreboardPlayer struct {
	Id         string
	Name       string
	Points     int
	LastPlayed string
}
