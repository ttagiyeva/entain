postgres:
	docker run --name entain -p 5432:5432 -e POSTGRES_USER=entain -e POSTGRES_PASSWORD=password -e POSTGRES_DB=entain -d postgres
generate:
	go generate ./...
test:
	go test -coverprofile=coverage.out ./... ;    go tool cover -html=coverage.out
