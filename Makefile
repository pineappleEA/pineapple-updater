GOFILES := src/main.go src/data.go src/ui.go

build:
	go build -ldflags "-s -w" -o bin/pineapple-updater ${GOFILES}

debug:
	go build -o bin/pineapple-updater ${GOFILES}

run:
	go run ${GOFILES}

clean:
	rm -rf bin

generate-icon:
	ifeq (, $(shell which fyne 2>/dev/null))
		$(error "fyne is not available in $\PATH, consider running `go get -u fyne.io/fyne`")
	endif
	fyne bundle -o src/data.go data/icon.png
