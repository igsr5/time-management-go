GOOS ?= linux
GOARCH ?= amd64
BIN_OUTPUT ?= bin/api-server
MAIN_PATH ?= cmd/server.go
DSN ?= mysql://root:root@(mysqld)/api-server
SQLBOILER_OUTPUT ?= app/models

.PHONY: setup
setup:
	go mod tidy 
	go get github.com/google/wire/cmd/wire@v0.5.0
	go generate ./...
	docker-compose build

.PHONY: build
build: 
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_OUTPUT) $(MAIN_PATH)

.PHONY: migrate-create
migrate-create:
	docker-compose exec app migrate create -ext sql -dir ./migrate ${file}

.PHONY: migrate-up
migrate-up:
	docker-compose exec app migrate -database "$(DSN)" -path migrate/. up

.PHONY: migrate-down
migrate-down:
	docker-compose exec app migrate -database "$(DSN)" -path migrate/. down

.PHONY: gen
gen:
	go generate ./...

