# Go-Users

Simple packages that implement a User interface that can be used to manage users consistently. 


```go

import (
 _ "github.com/eguevara/go-users/users/directory"
)

// Create an instance of the gitlab user service.
gitlabSvc, err := users.Dial("gitlab",
    users.WithKeyValue("token", "token"),
    users.WithKeyValue("baseURL", "https://gitlab.com/api/v3"),
)
if err != nil {
    logrus.Fatalf("error creating gitlab service: %s", err)
}

Options.GitlabUserService = gitlabSvc
```