package repository

import (
	"database/sql"
	"fmt"
	"github.com/sylph4/pipedrive-challenge/storage/postgres"

	"github.com/google/go-github/github"
)

type IGistRepository interface {
	SelectLastGistByUserName(userName string) (*postgres.Gist, error)
	InsertGist(gist *github.Gist) error
}

type GistRepository struct {
	conn *sql.DB
}

func NewGistRepository(conn *sql.DB) *GistRepository {
	return &GistRepository{conn: conn}
}

func (gr *GistRepository) SelectLastGistByUserName(userName string) (*postgres.Gist, error) {
	row := gr.conn.QueryRow("SELECT * FROM gists WHERE user_name=$1 ORDER BY created_at DESC LIMIT 1", userName)

	gist := &postgres.Gist{}
	err := row.Scan(&gist.ID, &gist.UserName, &gist.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return gist, nil
}

func (gr *GistRepository) InsertGist(gist *github.Gist) error {
	_, err := gr.conn.Exec("INSERT INTO Gists(id, user_name, created_at)  VALUES ($1, $2, $3)",
		*gist.ID, *gist.Owner.Login, gist.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Println("New gist since last check added: ", gist)

	return nil
}
