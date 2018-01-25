package gitlab

import (
	"errors"

	"github.com/eguevara/go-users/users"
	"github.com/xanzy/go-gitlab"
)

// userService is the implementation of a User type.
type userService struct {
	client *gitlab.Client
}

const (
	token   = "token"
	baseURL = "baseURL"
)

// NewUserService returns a userService implementation for gitlab api.
func NewUserService(opts *users.Opts) (users.UserService, error) {
	token, ok := opts.Opts[token]
	if !ok {
		return nil, errors.New("token is not provided")
	}

	client := gitlab.NewClient(nil, token)

	baseURL, ok := opts.Opts[baseURL]
	if ok {
		client.SetBaseURL(baseURL)
	}

	return &userService{client: client}, nil
}

func (u *userService) Get(id string) (*users.User, error) {
	return nil, errors.New("not implemented")
}

func (u *userService) Disable(id int) error {
	return u.client.Users.BlockUser(id)
}

// List will return a full list of gitlab users.
func (u *userService) List() (users.Users, error) {
	opts := &gitlab.ListUsersOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 1000,
		},
		Active: gitlab.Bool(true),
	}

	resp, _, err := u.client.Users.ListUsers(opts)
	if err != nil {
		return nil, err
	}

	list := make(users.Users, len(resp))

	for i, u := range resp {
		list[i] = users.User{
			Name:     u.Name,
			UserName: u.Username,
			Status:   u.State,
			ID:       u.ID,
		}
	}

	return list, nil

}

func init() {
	users.Register("gitlab", NewUserService)
}
