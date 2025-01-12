#  Aqua
Assessment


### Requirements

the following are the software requirements for running this project

- Have Golang installed (this repo uses Go version **1.22.7** specifically)
- Postgres database instance

The project is divided into two services:

- Notification Service
- Core Service

The notification service is responsible for handling notifications and storing them.

The core service is responsible for handling customers and resources.



### Setting Up Development Environment

1. Clone the repository

```sh
 git git@github.com:papidb/aqua.git
 cd aqua
```

2. Setup your env variables: a copy of the required variables can be found in the `.env.example` file. copy and configure using your custom setup

```
 cp .env.example .env
```

3. Install Dependencies

```sh
 go mod tidy
```

4. Setup your database: the database specified in .env `POSTGRES_DATABASE` must be created before proceeding

5. Run in development mode (hot reload): running in development mode requires an installation of the [Air CLI tool](https://github.com/cosmtrek/air)

```sh
air -c .air.core.toml 
air -c .air.notification.toml
```

or Run 

```
docker-compose -f docker-compose.dev.yml up -d
```

To Run in Production mode

```
docker-compose -f docker-compose.yml up -d
```


### Notes on Key functionalities

- We use zerolog for logging
- We use [Bun SQL client](https://bun.uptrace.dev/), previously known as `go-pg` for Database access.
- DB migrations are very low level: **hand written atomic SQL commands** , managed using [golang migrate](https://github.com/golang-migrate/migrate) and are run automatically on server startup
- [Go Gin](https://github.com/gin-gonic/gin) (as well as it's middlewares) for routing
- [Ozzo](https://github.com/go-ozzo/ozzo-validation) is used as the validator for HTTP requests

### Deployments

The project is packaged using docker containers. a production grade `Dockerfile` already exists and be used to replicate remote environments locally



Generate the proto file:

```bash
protoc --go_out=. --go-grpc_out=. notification.proto
```