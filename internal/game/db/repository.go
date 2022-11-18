package db

import (
	"database/sql"
	"tour-le-shit-go/internal/game/model"
	"tour-le-shit-go/internal/ierrors"
	playersModel "tour-le-shit-go/internal/players/model"
)

const GET_SCOREBOARD_QUERY = `
	SELECT s.points, p.name, p.id, s.last_played 
	from scoreboard s inner join player p ON(s.player_id = p.id) 
	WHERE season = $1;`

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) PostgresRepository {
	return PostgresRepository{db: db}
}

func (r PostgresRepository) GetScore(season int) ([]model.PlayerScore, error) {
	stmt, err := r.db.Prepare(GET_SCOREBOARD_QUERY)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error preparing statement from db " + err.Error(),
		}
	}

	rows, err := stmt.Query(season)

	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error fetching from db: " + err.Error(),
		}
	}

	p, err := getPlayers(rows)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error parsing result from db " + err.Error(),
		}
	}

	return p, nil
}

func getPlayers(rows *sql.Rows) ([]model.PlayerScore, error) {
	players := make([]model.PlayerScore, 0)

	for rows.Next() {
		var points int

		var name string

		var id string

		var lastPlayed string

		err := rows.Scan(&points, &name, &id, &lastPlayed)

		if err != nil {
			return nil, ierrors.DbError{
				Message: "error trying to scan rows " + err.Error(),
			}
		}

		players = append(players, model.PlayerScore{
			Player: playersModel.Player{
				Id:   id,
				Name: name,
			},
			Points:     points,
			LastPlayed: lastPlayed,
		})
	}

	return players, nil
}
