# Setup 
## Precommit Hooks
- `go install github.com/automation-co/husky@latest`
- `husky install`

## Golangci-lint
- [install golangci-lint locally](https://golangci-lint.run/welcome/install/#local-installation)


# Authorization Service

This is the **structure** for an Authorization/Authentification service.
It is intended to provide APIs for verifying user permissions so that applications can integrate secure role and permission-based access.

> ⚠️ This only contains the  **project skeleton**

## 📁 auth-service Structure

    auth-service
        ├── internal
        │   ├── auth
        │   │   ├── jwt
        │   │   │   ├── jwt_test.go
        │   │   │   └── jwt.go
        │   │   ├── types
        │   │   │   └── types.go
        │   │   └── services.go
        │   ├── server
        │   │   └── server.go
        ├── cmd
        │   └── main.go
        ├── .artifactignore
        ├── .gitignore
        ├── Dockerfile
        ├── go.mod
        ├── go.sum
        ├── odj-azure-pipeline....
        └── sonar-project.properties

* internal: This directory is for code that is private to your application. Other projects cannot import the code inside this directory.

* internal/auth: This package contains authentication-related logic, such as a service for user authentication, authorization or encryption.

* internal/server: This package contains the code that sets up and runs your HTTP server.

* cmd/main.go: The entry point of your application

## ⚙️ To run the service use the command (in progress)

    run go ./cmd/main.go

## ⚙️ To test the jwt token use the command in jwt folder

    go test -v