# manabase

## develop

Prerequisites

- [go](https://go.dev/) 1.22+
- local postgres 16 database

Setup

1. open repository locally and run `go mod tidy`
1. copy `.env.example` to `.env`
1. update `.env` with your db user, password, and database name
1. run `goose -dir sql/schema postgres postgres://user:password@localhost:5432/database up`, replacing user, password, and database with your own
1. run `make dev` to start the app
1. access at `localhost:3000` 
