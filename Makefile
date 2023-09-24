build:
	go build -o bin/simulationgame

run: build
	./bin/simulationgame

test:
	go test -v ./...
