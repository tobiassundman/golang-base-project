# demo-app

## demo-app

The demo app is a project I have created to try out various things in go web development.

It is a service with a crud API /v1/users/* to manage users stored in a postgres database.

```go
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}
```

The controller layer is located in internal/app/controller, the business layer in internal/app/service and the repository layer in internal/app/repository.

Run `make` or `make help` to view information about available commands

### Setup

Run `make tools` to install necessary tools to use the Makefile

### Building

Run `make build` to check, build and test

## db-migration

The db-migration application is used to run sql migrations against a postgres database without any more manual steps.