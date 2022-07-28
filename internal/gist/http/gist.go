package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sylph4/pipedrive-challenge/internal/gist/model"
	"github.com/sylph4/pipedrive-challenge/internal/gist/repository"
	"github.com/sylph4/pipedrive-challenge/internal/gist/service"
	"github.com/sylph4/pipedrive-challenge/storage/postgres"
)

type GistHandler struct {
	gistService    service.IGistService
	userRepository repository.IUserRepository
	gistRepository repository.IGistRepository
}

func NewGistHandler(gistService service.IGistService, userRepository *repository.UserRepository,
	gistRepository *repository.GistRepository) *GistHandler {
	return &GistHandler{
		gistService:    gistService,
		userRepository: userRepository,
		gistRepository: gistRepository,
	}
}

func (h *GistHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	user := &model.CreateUserRequest{}
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		fmt.Println("Request validation error: ", err)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)

		return
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

	request := &model.DeleteUserNameRequest{}
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		fmt.Println("Request validation error: ", err)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)

		return
	}

	if request.UserName == "" {
		http.Error(w, http.StatusText(204), http.StatusNoContent)

		return
	}

	existingUser, err := h.userRepository.SelectUserByUserName(request.UserName)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
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

func (h *GistHandler) GetNewUserGists(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)

		return
	}

	username := r.URL.Query().Get("username")
	if username != "" {
		fmt.Println("Request validation error: username param required")
		http.Error(w, http.StatusText(400), http.StatusBadRequest)

		return
	}

	gists, err := h.gistRepository.SelectNewGistsByUserName(username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	err = h.gistRepository.MarkGistsAsChecked(username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(gists)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)

		return
	}

	//nolint
	w.Write(response)
}

func (h *GistHandler) RunGistsCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)

		return
	}

	h.gistService.RunGistCheck()

	http.StatusText(200)
}
