
default: build-app
build: build-app

build-app:
	go build -o bin/gocomments -v .

run:
	./bin/gocomments .

generate-dataset:
	python ./dataset/generate_func_comments_from_local_repo.py

generate-model:
	pip install -q datasets
	pip install transformers
	pip install tf-keras
	pip install torch
	pip install python-dotenv
	pip install requests pygments
	python ./model/train.py

test-model:
	go run ./model/prompt.go

install:
	sudo go install ./
