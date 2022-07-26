package repository

import (
	"database/sql"
	"fmt"

	"github.com/sylph4/pipedrive-challenge/storage/postgres"
)

type IUserRepository interface {
	SelectAllUsers() ([]*postgres.User, error)
	SelectUserByUserName(userName string) (*postgres.User, error)
	InsertUser(user postgres.User) error
	DeleteUserByUserName(userName string) error
}

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn: conn}
}

func (ur *UserRepository) SelectAllUsers() ([]*postgres.User, error) {
	rows, err := ur.conn.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]*postgres.User, 0)

	for rows.Next() {
		user := &postgres.User{}
		err := rows.Scan(&user.UserName, &user.GithubAPIKey, &user.PipedriveAPIKey, &user.PipedriveUserID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) DeleteUserByUserName(userName string) error {
	_, err := ur.conn.Exec("DELETE FROM users WHERE user_name=$1", userName)

	if err != nil {
		return err
	}

	fmt.Println("User deleted: ", userName)

	return nil
}

func (ur *UserRepository) SelectUserByUserName(userName string) (*postgres.User, error) {
	row := ur.conn.QueryRow("SELECT * FROM users WHERE user_name=$1", userName)

	user := &postgres.User{}
	err := row.Scan(&user.UserName, &user.GithubAPIKey, &user.PipedriveAPIKey, &user.PipedriveUserID)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) InsertUser(user postgres.User) error {
	_, err := ur.conn.Exec("INSERT INTO users(user_name, github_api_key, pipedrive_api_key, pipedrive_user_id)  VALUES ($1, $2, $3, $4)",
		user.UserName, user.GithubAPIKey, user.PipedriveAPIKey, user.PipedriveUserID)

	if err != nil {
		return err
	}

	fmt.Println("New user added: ", user)

	return nil
}
