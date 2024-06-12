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

## Run locally
Clone the project

`git clone https://github.com/ttagiyeva/entain.git`

Setup database

`make postgres`

Set configurations. Project contains `.env/dev` file to ease setup environment variables

`export $(cat .env/dev)`

Run service

`go run cmd/main.go`

Mock request 

`curl --location 'http://localhost:8080/api/v1/users/00000000-0000-0000-0000-000000000001/transactions' --header 'Content-Type: application/json' --header 'Content-Type: application/json' --header 'Source-Type: game' --data '{
    "state": "win",
    "amount": 10.15,
    "transactionId": "1"
}'`

## Run tests

1. Generate mocks
`make generate`
2. Run tests 
`make test`    
