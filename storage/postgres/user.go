package postgres

type User struct {
	UserName        string
	GithubAPIKey    string
	PipedriveAPIKey string
	PipedriveUserID uint
}
