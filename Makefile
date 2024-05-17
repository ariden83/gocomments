
default: build-app
build: build-app

build-app:
	go build -o bin/gocomments -v main.go

run:
	./bin/gocomments .
