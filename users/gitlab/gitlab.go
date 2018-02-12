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
func (u *userService) List(opts users.ListOptions) (users.Users, error) {
	var options gitlab.ListUsersOptions
	options, ok := opts.(gitlab.ListUsersOptions)
	if ok != true {
		return nil, errors.New("type assertion on opts")
	}

	var allUsers []*gitlab.User
	for {
		req, resp, err := u.client.Users.ListUsers(&options)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, req...)
		if resp.NextPage == 0 {
			break
		}

		options.ListOptions.Page = resp.NextPage
	}

	list := make(users.Users, len(allUsers))

	for i, u := range allUsers {
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
