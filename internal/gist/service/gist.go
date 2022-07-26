package service

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"time"

	"github.com/genert/pipedrive-api/pipedrive"
	"github.com/google/go-github/github"
	"github.com/sylph4/pipedrive-challenge/internal/gist/repository"
	"github.com/sylph4/pipedrive-challenge/storage/postgres"
)

type IGistService interface {
	RunGistCheck()
}

type GistService struct {
	gistRepository repository.IGistRepository
	userRepository repository.IUserRepository
}

func NewGistService(gistRepository repository.IGistRepository, userRepository *repository.UserRepository) *GistService {
	gistService := &GistService{gistRepository: gistRepository, userRepository: userRepository}

	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		for range ticker.C {
			gistService.RunGistCheck()
		}
	}()

	return gistService
}

func (s *GistService) RunGistCheck() {
	fmt.Println("Gists check start at: ", time.Now().Format("2006-01-02T15:04:05"))
	ctx := context.Background()

	users, err := s.userRepository.SelectAllUsers()
	if err != nil {
		panic(err)
	}

	if len(users) == 0 {
		fmt.Println("check start at: No users found")
	}

	for i := range users {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: users[i].GithubAPIKey},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		responseGists, _, err := client.Gists.List(ctx, "", nil)
		if err != nil {
			panic(err)
		}

		gist, err := s.gistRepository.SelectLastGistByUserName(*responseGists[0].Owner.Login)
		if err != nil {
			panic(err)
		}

		pipedriveClient := pipedrive.NewClient(&pipedrive.Config{
			APIKey:        users[i].PipedriveAPIKey,
			CompanyDomain: "api.pipedrive.com/v1",
		})

		for n := 0; n < len(responseGists); n++ {
			if gist != nil && gist.CreatedAt.Before(*responseGists[n].CreatedAt) && !responseGists[n].CreatedAt.Equal(gist.CreatedAt) {
				err = createActivity(ctx, responseGists[n], users[i], pipedriveClient)
				if err != nil {
					fmt.Println("Could not create a gist: ", err)

					break
				}

				err = s.gistRepository.InsertGist(responseGists[i])
				if err != nil {

					panic(err)
				}
			} else {
				err = createActivity(ctx, responseGists[n], users[i], pipedriveClient)
				if err != nil {
					fmt.Println("Could not create a gist: ", err)

					break
				}

				err = s.gistRepository.InsertGist(responseGists[n])
				if err != nil {

					panic(err)
				}
			}
		}
	}

	fmt.Println("Gists check ended at: ", time.Now().Format("2006-01-02T15:04:05"))
}

func createActivity(ctx context.Context, gist *github.Gist, user *postgres.User, pipedriveClient *pipedrive.Client) error {
	_, res, err := pipedriveClient.Activities.Create(ctx, &pipedrive.ActivitiesCreateOptions{
		Subject:      gist.GetDescription(),
		Done:         1,
		Type:         "",
		DueDate:      gist.UpdatedAt.String(),
		DueTime:      "",
		Duration:     "",
		UserID:       user.PipedriveUserID,
		DealID:       0,
		PersonID:     0,
		Participants: nil,
		OrgID:        0,
	})
	if err != nil {
		return err
	}

	fmt.Println("Activity created: ", res)

	return nil
}
