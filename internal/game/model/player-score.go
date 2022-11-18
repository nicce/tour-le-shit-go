package model

import "tour-le-shit-go/internal/players/model"

type PlayerScore struct {
	Player     model.Player
	Points     int
	LastPlayed string
}
