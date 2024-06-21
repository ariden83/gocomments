
default: build-app
build: build-app

build-app:
	go build -o bin/gocomments -v .

run:
	./bin/gocomments .

generate-dataset:
	python ./dataset/generate-from-github-v2.py

generate-model:
	pip install -q datasets
	pip install transformers
	pip install tf-keras
	pip install torch
	pip install python-dotenv
	pip install requests pygments
	python ./model/train.py

install:
	sudo go install ./
