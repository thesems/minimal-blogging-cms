build:
	go build -o bin/app

run: build
	./bin/app

run-dev: build
	./bin/app -storage memory

test:
	go test -v ./... -count