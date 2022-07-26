package model

type CreateUserRequest struct {
	UserName        string `json:"userName" validate:"required,min=1,max=100"`
	GithubAPIKey    string `json:"githubAPIKey" validate:"required,max=40"`
	PipedriveAPIKey string `json:"pipedriveAPIKey" validate:"required,max=40"`
	PipedriveUserID uint   `json:"pipedriveUserID" validate:"required,numeric"`
}

type DeleteUserNameRequest struct {
	UserName string `json:"userName" validate:"required,min=1,max=100"`
}
