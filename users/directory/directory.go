package directory

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/eguevara/go-directory/directory"
	"github.com/eguevara/go-users/users"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// userService is the implementation of a User type.
type userService struct {
	client *directory.Client
}

// Keys used for storing dial options.
const (
	privateKeyFile  = "privateKeyFile"
	baseURL         = "baseURL"
	serviceEmail    = "serviceEmail"
	impersonateUser = "impersonateUser"
)

// NewUserService returns a userService implementation for directory api.
func NewUserService(opts *users.Opts) (users.UserService, error) {

	baseURL := opts.Opts[baseURL]

	var client *http.Client
	if _, ok := opts.Opts[privateKeyFile]; ok {
		client = getOAuthClient(opts.Opts[serviceEmail], opts.Opts[impersonateUser], opts.Opts[privateKeyFile])
	}

	directoryClient, err := directory.New(directory.SetHTTPClient(client), directory.SetBaseURL(baseURL))
	if err != nil {
		log.Fatalf("error in creating user client %v", err)
	}

	return &userService{client: directoryClient}, nil
}

func (u *userService) List(opt users.ListOptions) (users.Users, error) {
	return nil, errors.New("not implemented")
}

func (u *userService) Disable(id int) error {
	return errors.New("not implemented")
}

// Active will return false if the status is D for delete, true otherwise.
func (u *userService) Get(id string) (*users.User, error) {
	opt := &directory.UsersOptions{Fields: users.String("coreId,fullName,id,status")}
	resp, _, err := u.client.Users.Get(context.TODO(), id, opt)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Name:     resp.FullName,
		UserName: resp.ID,
		Status:   resp.Status,
	}

	return user, nil
}

func init() {
	users.Register("directory", NewUserService)
}

func getOAuthClient(serviceEmail, impersonateUser, keyFile string) *http.Client {

	if serviceEmail == "" {
		log.Fatal("service email option is required to create oauth httpclient")
	}

	if impersonateUser == "" {
		log.Fatal("impersonate user option is required to create oauth httpclient")
	}

	if keyFile == "" {
		log.Fatal("the file location of the pem key is required to create oauth httpclient")
	}

	configKeyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	conf := &jwt.Config{
		Email:      serviceEmail,
		PrivateKey: []byte(configKeyBytes),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		TokenURL: google.JWTTokenURL,
		Subject:  impersonateUser,
	}
	client := conf.Client(oauth2.NoContext)
	return client
}
