PROJECTNAME = auth-service

.PHONY: lint
lint:
	golangci-lint run  --config=.golangci.yaml --timeout=180s ./...


.PHONY: generate
generate:
	go generate ./..


.PHONY: run-migrate
run-migrate-local:
	sql-migrate up -env="local"

.PHONY: stop-migrate
stop-migrate-local:
	sql-migrate down -env="local"

.PHONY: build
build:
	go build -o ./build/${PROJECTNAME} ./cmd/${PROJECTNAME}/main.go || exit 1


.PHONY: run
run:
	go run ./cmd/${PROJECTNAME}/main.go
