postgres:
	docker run --name entain -p 5432:5432 -e POSTGRES_USER=entain -e POSTGRES_PASSWORD=password -e POSTGRES_DB=entain -d postgres
generate:
	go generate ./...
