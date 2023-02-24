# Entain


## Description
Service for transactions over user balance.

## Tech stack
Service has written using Go programming language.
* [Echo](https://echo.labstack.com/) as a HTTP framework
* [Viper](https://github.com/spf13/viper) for config management
* [Gomock](https://github.com/golang/mock) for mocking dependencies
* PostgreSQL
* Docker

## Installation
Clone the project

`git clone https://github.com/ttagiyeva/entain.git`

Setup database

`make postgres`

## Usage
Set configurations. Project contains `.env/dev` file to ease setup environment variables

`export $(cat .env/dev)`

Run service

`go run cmd/main.go`

Run tests

1. Generate mocks
`make generate`
2. Run tests 
`go test ./...`    
