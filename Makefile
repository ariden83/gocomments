
default: build-app
build: build-app

build-app:
	go build -o bin/gocomments -v .

run:
	./bin/gocomments .

generate-model:
	pip install -q datasets
	pip install transformers
	pip install tf-keras
	pip install torch
	python ./model/train.py

install:
	sudo go install ./
