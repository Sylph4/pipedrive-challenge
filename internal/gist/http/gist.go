package http

import (
	"encoding/json"
	"fmt"
	"github.com/sylph4/pipedrive-challenge/storage/postgres"
	"net/http"

	"github.com/sylph4/pipedrive-challenge/internal/gist/repository"
	"github.com/sylph4/pipedrive-challenge/internal/gist/service"
)

type GistHandler struct {
	gistService    service.IGistService
	userRepository repository.IUserRepository
}

func NewGistHandler(gistService service.IGistService, userRepository *repository.UserRepository) *GistHandler {
	return &GistHandler{gistService: gistService, userRepository: userRepository}
}

func (h *GistHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	user := &CreateUserRequest{}
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	existingUser, err := h.userRepository.SelectUserByUserName(user.UserName)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	if existingUser != nil {
		fmt.Println("CreateUserRequest already exists")
		http.Error(w, http.StatusText(400), http.StatusBadRequest)

		return
	}

	newUser := postgres.User{
		UserName:        user.UserName,
		GithubAPIKey:    user.GithubAPIKey,
		PipedriveAPIKey: user.PipedriveAPIKey,
		PipedriveUserID: user.PipedriveUserID,
	}

	err = h.userRepository.InsertUser(newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.StatusText(201)
}

func (h *GistHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userRepository.SelectAllUsers()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	//nolint
	w.Write(response)
}

func (h *GistHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)

	request := &DeleteUserNameRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	if request.UserName == "" {
		http.Error(w, http.StatusText(204), http.StatusNoContent)
	}

	existingUser, err := h.userRepository.SelectUserByUserName(request.UserName)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}

	if existingUser == nil {
		fmt.Println("No user with name: ", request.UserName)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)

		return
	}

	err = h.userRepository.DeleteUserByUserName(request.UserName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.StatusText(201)
}

type CreateUserRequest struct {
	UserName        string `json:"userName"`
	GithubAPIKey    string `json:"githubAPIKey"`
	PipedriveAPIKey string `json:"pipedriveAPIKey"`
	PipedriveUserID uint   `json:"pipedriveUserID"`
}

type DeleteUserNameRequest struct {
	UserName string `json:"userName"`
}
