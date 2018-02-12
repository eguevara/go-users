package users

import (
	"fmt"
	"log"
	"sync"
)

var (
	providersMutex sync.Mutex

	// registration is an in-memory map of implemented User services.
	registration = make(map[string]UserConstructor)
)

// User stores user fields exposed to the caller.
type User struct {
	// Name is the given name (fullname) of any user.
	Name string

	// UserName is the unique user identifier.
	UserName string

	// Status is the user status for a given user service.
	Status string

	// Status is the unique numeric ID for a user.
	ID int
}

// ListOptions is a wrapper interface for list options.
type ListOptions interface{}

// Users is a list of User types
type Users []User

// UserService is an abstract interface to be implemented by any User type.
type UserService interface {
	// List will return a list of users from the user service.
	List(ListOptions) (Users, error)

	// Get will return a User type from the user service.
	Get(string) (*User, error)

	// Disable will disable the user (not delete) from the user service.
	Disable(int) error
}

// Opts holds configuration options for the user backend.
// It is meant to be used by implementations of UserService.
type Opts struct {
	Opts map[string]string // key-value pair
}

// DialOpts is a daisy-chaining mechanism for setting options to a backend during Dial.
type DialOpts func(*Opts) error

// UserConstructor is a function that initializes and returns a UserService
// implementation with the given options.
type UserConstructor func(*Opts) (UserService, error)

// Register registers a user provider
func Register(name string, fn UserConstructor) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := registration[name]; found {
		log.Fatalf("User provider %q was registered twice", name)
	}
	log.Printf("Registered user provider %q", name)
	registration[name] = fn
}

// UserProviders returns the name of all registered user providers in a
// string slice
func UserProviders() []string {
	names := []string{}
	providersMutex.Lock()
	defer providersMutex.Unlock()
	for name := range registration {
		names = append(names, name)
	}
	return names
}

// WithKeyValue sets a key-value pair as option. If called multiple times with the same key, the last one wins.
func WithKeyValue(key, value string) DialOpts {
	return func(o *Opts) error {
		o.Opts[key] = value
		return nil
	}
}

// Dial dials the named user backend using the dial options opts.
func Dial(name string, opts ...DialOpts) (UserService, error) {
	fn, found := registration[name]
	if !found {
		return nil, fmt.Errorf("could not find user provider: %s", name)
	}
	dOpts := &Opts{Opts: make(map[string]string)}
	var err error
	for _, o := range opts {
		if o != nil {
			err = o(dOpts)
			if err != nil {
				return nil, err
			}
		}
	}
	return fn(dOpts)
}

// String is a helper function that allocates a new string value
func String(v string) *string { return &v }

// Bool is a helper function that allocates a new bool value
func Bool(v bool) *bool { return &v }

// Int is a helper function that alloates a new int value
func Int(v int) *int { return &v }
