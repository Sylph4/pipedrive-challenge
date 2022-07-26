package postgres

import "time"

type Gist struct {
	ID        string
	UserName  string
	CreatedAt time.Time
}
