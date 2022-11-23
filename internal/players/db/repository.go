package db

import (
	"database/sql"
	"fmt"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/players/model"

	"github.com/google/uuid"
)

const GetPlayerByIdQuery = "SELECT id, name FROM player WHERE id = $1"
const GetCountPlayersByIdQuery = "SELECT count(*) FROM player WHERE id = $1"
const GetCountPlayersByNameQuery = "SELECT count(*) FROM player WHERE name = $1"
const GetPlayersQuery = "SELECT id, name from player ORDER BY name;"
const InsertPlayerQuery = "INSERT INTO player (id, name) VALUES ($1, $2);"
const UpdatePlayerQuery = "UPDATE player SET name = $2 WHERE id = $1;"
const DeletePlayerQuery = "DELETE FROM player WHERE id = $1;"

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetPlayerById(id string) (*model.Player, error) {
	stmt, err := r.db.Prepare(GetPlayerByIdQuery)
	if err != nil {
		return &model.Player{}, ierrors.DbError{
			Message: fmt.Sprintf("error preparing statement %v", err),
		}
	}

	row := stmt.QueryRow(id)

	if row == nil {
		return nil, nil
	}

	var playerId string

	var playerName string

	err = row.Scan(&playerId, &playerName)
	if err != nil {
		return &model.Player{}, ierrors.DbError{
			Message: fmt.Sprintf("error scanning result %v", err),
		}
	}

	return &model.Player{
		Id:   playerId,
		Name: playerName,
	}, nil
}

func (r *PostgresRepository) GetPlayers() ([]model.Player, error) {
	rows, err := r.db.Query(GetPlayersQuery)
	if err != nil {
		return nil, ierrors.DbError{Message: fmt.Sprintf("error fetching players %v", err)}
	}

	players := make([]model.Player, 0)

	for rows.Next() {
		var id string

		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, ierrors.DbError{Message: fmt.Sprintf("error scanning rows %v", err)}
		}

		players = append(players, model.Player{
			Id:   id,
			Name: name,
		})
	}

	return players, nil
}

func (r *PostgresRepository) CreatePlayer(name string) ([]model.Player, error) {
	count, err := r.countPlayersByQuery(GetCountPlayersByNameQuery, name)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    fmt.Sprintf("player with name %s already exists.", name),
			InnerError: "",
		}
	}

	stmt, err := r.db.Prepare(InsertPlayerQuery)
	if err != nil {
		return nil, ierrors.DbError{
			Message: fmt.Sprintf("error preparing insert player query %v", err),
		}
	}

	id := uuid.New().String()

	_, err = stmt.Exec(id, name)
	if err != nil {
		return nil, ierrors.DbError{
			Message: fmt.Sprintf("error executing statement insert player %v", err),
		}
	}

	return r.GetPlayers()
}

func (r *PostgresRepository) UpdatePlayer(id, name string) ([]model.Player, error) {
	count, err := r.countPlayersByQuery(GetCountPlayersByIdQuery, id)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    fmt.Sprintf("player with id %s does not exist", id),
			InnerError: "",
		}
	}

	stmt, err := r.db.Prepare(UpdatePlayerQuery)
	if err != nil {
		return nil, ierrors.DbError{
			Message: fmt.Sprintf("error preparing udpate player query %v", err),
		}
	}

	_, err = stmt.Exec(id, name)
	if err != nil {
		return nil, ierrors.DbError{Message: fmt.Sprintf("error executing update player query %v", err)}
	}

	return r.GetPlayers()
}

func (r *PostgresRepository) DeletePlayer(id string) ([]model.Player, error) {
	stmt, err := r.db.Prepare(DeletePlayerQuery)
	if err != nil {
		return nil, ierrors.DbError{Message: fmt.Sprintf("error preparing delete player query %v", err)}
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return nil, ierrors.DbError{
			Message: fmt.Sprintf("error executing delete player query %v", err),
		}
	}

	return r.GetPlayers()
}

func (r *PostgresRepository) countPlayersByQuery(query string, param string) (int, error) {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, ierrors.DbError{Message: fmt.Sprintf("error preparing player count query %v", err)}
	}

	rows, err := stmt.Query(param)
	if err != nil {
		return 0, ierrors.DbError{Message: fmt.Sprintf("error running player count query %v", err)}
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, ierrors.DbError{Message: fmt.Sprintf("error scanning result of player count query %v", err)}
		}
	}

	return count, nil
}
