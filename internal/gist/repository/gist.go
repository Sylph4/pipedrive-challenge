package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/sylph4/pipedrive-challenge/storage/postgres"
)

type IGistRepository interface {
	SelectLastGistByUserName(userName string) (*postgres.Gist, error)
	InsertGist(gist *github.Gist) error
	SelectNewGistsByUserName(userName string) ([]*postgres.Gist, error)
	MarkGistsAsChecked(userName string) error
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

func (gr *GistRepository) SelectNewGistsByUserName(userName string) ([]*postgres.Gist, error) {
	rows, err := gr.conn.Query(
		"SELECT * FROM gists WHERE user_name=$1 AND is_checked IS FALSE ORDER BY created_at DESC", userName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	gists := make([]*postgres.Gist, 0)

	for rows.Next() {
		gist := &postgres.Gist{}
		err := rows.Scan(&gist.ID, &gist.UserName, &gist.CreatedAt)
		if err != nil {
			return nil, err
		}

		gists = append(gists, gist)
	}

	return gists, nil
}

func (gr *GistRepository) MarkGistsAsChecked(userName string) error {
	_, err := gr.conn.Exec(
		"UPDATE Gists SET is_checked=TRUE WHERE user_name=$1", userName)

	if err != nil {
		return err
	}

	return nil
}
