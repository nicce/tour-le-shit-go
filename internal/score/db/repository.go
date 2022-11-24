package db

import (
	"database/sql"
	"fmt"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/players"
	"tour-le-shit-go/internal/score/model"
	"tour-le-shit-go/internal/utils"

	"github.com/google/uuid"
)

const GetPlayerScoreBySeasonQuery = `
	SELECT s.id, s.player_id, p.name, s.points, s.birdies, s.eagles, s.muligans, s.season, s.day
	FROM score s INNER JOIN player p on (s.player_id = p.id) 
	WHERE s.player_id=$1 and season=$2;
`

const DeleteScoreById = `DELETE FROM score WHERE id=$1;`

const InsertScoreQuery = `
	INSERT INTO score (id, player_id, points, birdies, eagles, muligans, season, day) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`

const GetScoreboardQuery = `
	WITH player_points as (
        SELECT s.player_id, sum(s.points) + 2*SUM(s.birdies) + 3*SUM(s.eagles) - 3*SUM(s.muligans) AS points, MAX(s.day) AS day
        FROM score s
        WHERE season=$1
        GROUP BY s.player_id
	)
	SELECT COALESCE(pp.points, 0), p.id AS player_id, p.name, COALESCE(pp.day, '') AS last_played 
	FROM player_points pp 
	RIGHT JOIN player p ON (pp.player_id = p.id);
`

type PostgresRepository struct {
	db                *sql.DB
	playersRepository players.Repository
}

func NewRepository(db *sql.DB, repository players.Repository) *PostgresRepository {
	return &PostgresRepository{db: db, playersRepository: repository}
}

func (r *PostgresRepository) GetPlayerScore(id string, season int) ([]model.Score, error) {
	stmt, err := r.db.Prepare(GetPlayerScoreBySeasonQuery)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error preparing statement from db " + err.Error(),
		}
	}

	rows, err := stmt.Query(id, season)

	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error fetching from db: " + err.Error(),
		}
	}

	p, err := getPlayerScores(rows)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error parsing result from db " + err.Error(),
		}
	}

	return p, nil
}

func (r *PostgresRepository) DeleteScore(id string) error {
	stmt, err := r.db.Prepare(DeleteScoreById)
	if err != nil {
		return ierrors.DbError{
			Message: "Error preparing statement from db " + err.Error(),
		}
	}

	_, err = stmt.Exec(id)

	if err != nil {
		return ierrors.DbError{
			Message: "Error executing statement from db: " + err.Error(),
		}
	}

	return nil
}

func (r *PostgresRepository) AddScore(scoreInput model.ScoreInput) (*model.Score, error) {
	player, err := r.playersRepository.GetPlayerById(scoreInput.PlayerId)
	if err != nil {
		return nil, ierrors.DbError{
			Message: fmt.Sprintf("error fetching player %v", err),
		}
	}

	if player == nil {
		return nil, ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "player does not exists",
			InnerError: "",
		}
	}

	stmt, err := r.db.Prepare(InsertScoreQuery)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error preparing statement from db " + err.Error(),
		}
	}

	score := model.Score{
		Id:         uuid.New().String(),
		PlayerId:   player.Id,
		PlayerName: player.Name,
		Points:     scoreInput.Points,
		Birdies:    scoreInput.Birdies,
		Eagles:     scoreInput.Eagles,
		Muligans:   scoreInput.Muligans,
		Season:     scoreInput.Season,
		Day:        utils.GetToday(),
	}

	_, err = stmt.Exec(score.Id, score.PlayerId, score.Points, score.Birdies, score.Eagles, score.Muligans, score.Season, score.Day)
	if err != nil {
		return nil, ierrors.DbError{
			Message: "Error executing statement from db: " + err.Error(),
		}
	}

	return &score, nil
}

func (r *PostgresRepository) GetScoreboard(season int) (model.Scoreboard, error) {
	stmt, err := r.db.Prepare(GetScoreboardQuery)
	if err != nil {
		return model.Scoreboard{}, ierrors.DbError{
			Message: "Error preparing statement from db " + err.Error(),
		}
	}

	rows, err := stmt.Query(season)
	if err != nil {
		return model.Scoreboard{}, ierrors.DbError{
			Message: "Error querying db " + err.Error(),
		}
	}

	players := make([]model.ScoreboardPlayer, 0)

	for rows.Next() {
		var playerName string

		var playerId string

		var points int

		var lastPlayed string

		err = rows.Scan(&points, &playerId, &playerName, &lastPlayed)
		if err != nil {
			return model.Scoreboard{}, ierrors.DbError{
				Message: "Error scanning rows: " + err.Error(),
			}
		}

		players = append(players, model.ScoreboardPlayer{
			Id:         playerId,
			Name:       playerName,
			Points:     points,
			LastPlayed: lastPlayed,
		})
	}

	return model.Scoreboard{
		Players: players,
		Season:  season,
	}, nil
}

func getPlayerScores(rows *sql.Rows) ([]model.Score, error) {
	playerScores := make([]model.Score, 0)

	for rows.Next() {
		var id string

		var playerId string

		var playerName string

		var points int

		var birdies int

		var eagles int

		var muligans int

		var season int

		var day string

		err := rows.Scan(&id, &playerId, &playerName, &points, &birdies, &eagles, &muligans, &season, &day)

		if err != nil {
			return nil, ierrors.DbError{
				Message: "error trying to scan rows " + err.Error(),
			}
		}

		playerScores = append(playerScores, model.Score{
			Id:         id,
			PlayerId:   playerId,
			PlayerName: playerName,
			Points:     points,
			Birdies:    birdies,
			Eagles:     eagles,
			Muligans:   muligans,
			Season:     season,
			Day:        day,
		})
	}

	return playerScores, nil
}
