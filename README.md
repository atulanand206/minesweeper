# Minesweeper

The famous game since the early days of Microsoft where you have to flag the cells with mine while opening the cells with value and get the highest score possible.

### Microservices Ecosystem
This repository is part of a larger ecosystem of microservices each serving it's purpose in bring Minesweeper to life.

- [Web Services](https://github.com/atulanand206/minesweeper)
- [Authorization](https://github.com/atulanand206/users)
- [Update Database](https://github.com/atulanand206/ms-db-publisher)
- [Web Interface](https://github.com/atulanand206/mines)

You can use the `docker-compose.yml` file to start all the services in separate docker containers in one go. This file is available in a different repository. Please let me know if you'd like to experiment with that.

### About
This repository contains the Web services APIs for the game in golang. 

#### Endpoints

- `/game/new` - Creates a board for rows, columns and mines count.
- `/game/save` - Publishes a game to the *Update Database* Service.
- `/games` - Returns the leaderboard for games linked to a configuration.

#### Environment Variables

- `MONGO_CLIENT_ID` - MongoDB Cloud Atlas client id
- `DATABASE` - MongoDB Database
- `MONGO_COLLECTION` - MongoDB Collection
- `KAFKA_TOPIC` - Kafka Cluster topic
- `KAFKA_BROKER_ID` - Kafka Cluster broker id
- `USERS_URL` - Users service URL
- `CLIENT_SECRET` - Client secret to verify access token.
- `CORS_ORIGIN` - CORS Origin to allow.

#### Run the service

The service can be run either directly as a go project or use a Docker container to build and run.

The projects uses go modules to manage the dependencies.

Run the following commands in order to build the project to build from the command line.

```
go get -d -v ./...
go mod download
go mod vendor
go mod verify
go build
go run main.go
```
Or the following to build the docker container.
```
docker build -t minesweeper .
docker run -it minesweeper
```

Ensure that the `go.sum` and `vendor/modules.txt` do not have it's content altered unless there is a dependency alteration.

### Author

- Atul Anand
