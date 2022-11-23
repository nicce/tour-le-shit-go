# Tour le shit in GO

![workflow](https://github.com/nicce/tour-le-shit-go/actions/workflows/main.yaml/badge.svg)

Implementation of the infamous [Tour le shit](https://github.com/nicce/tour-le-shit) backend using GoLang.

## Development

### Start
You can run the app in many different ways:

`go run main.go`

For running the app with live reload support use [air](https://github.com/cosmtrek/air) and then run the command:

`air`

### Build

`go build`

## Environment variable

For convenient use `.env` file in root folder. Check `.env.default` for default values

| key               | description       |
|-------------------|-------------------|
| DATABASE_NAME     | Database name     |
| DATABASE_PASSWORD | Database password |
| DATABASE_USER     | Database user     |
| MEMBERS_MODE      | MOCK or PSQL      |
| PORT              | Server port       |
| SCORE_MODE        | MOCK or PSQL      |