GOFILES := src/main.go src/data.go

build:
	go build -o bin/pineapple-updater ${GOFILES}

run:
	go run ${GOFILES}
