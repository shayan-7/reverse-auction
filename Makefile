build:
	@go build -o http ./main.go
.PHONY: http

run:
	@go run ./main.go
.PHONY: run

swagger:
	@swag init -g handlers/server.go
.PHONY: swagger

docker:
	@docker build -t uniproject:0.1 .
.PHONY: docker

db:
	@sqlite3 test.db
.PHONY:db

up: docker
	@docker run -p 8080:8080 uniproject:0.1
