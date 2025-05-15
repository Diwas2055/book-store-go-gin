# CRUD REST API using Postgresql DB in GO

> Simple CRUD application that uses `Go Gin framework` to route HTTP requests to handlers that correspond to the CRUD operations. We are going to use the `PostgreSQL` database as our datastore. Beside that, we are also going to incorporate a developer friendly ORM called `Gorm`.

## Application structure

> The application structure that we will use is quite simple. It is made of two modules namely:

- `database` module: For all database related operations such as connection to the database
- `controllers` module: For CRUD operations.

At the root level we will have the `main.go` file that will host code for initializing the database, and HTTP routing and running the server.

## Install dependencies

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/google/uuid
```

## Run the application
```bash
go run main.go
```

## Test the Application
```
bash test_api.sh
````
