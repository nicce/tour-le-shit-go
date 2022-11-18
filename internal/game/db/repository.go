package db

import (
	"database/sql"
	"tour-le-shit-go/internal/game/model"
	"tour-le-shit-go/internal/ierrors"
)

const GET_SCOREBOARD_QUERY = `
	SELECT s.points, p.name, s.last_played 
	from scoreboard s inner join player p ON(s.player_id = p.id) 
	WHERE season = $1;`

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) PostgresRepository {
	return PostgresRepository{db: db}
}

func (r PostgresRepository) GetScore(season int) ([]model.Player, error) {
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

func getPlayers(rows *sql.Rows) ([]model.Player, error) {
	players := make([]model.Player, 0)

	for rows.Next() {
		var points int

		var name string

		var lastPlayed string

		err := rows.Scan(&points, &name, &lastPlayed)

		if err != nil {
			return nil, ierrors.DbError{
				Message: "error trying to scan rows " + err.Error(),
			}
		}

		players = append(players, model.Player{
			Name:       name,
			Points:     points,
			LastPlayed: lastPlayed,
		})
	}

	return players, nil
}
