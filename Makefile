
default: build-app
build: build-app

build-app:
	go build -o bin/gocomments -v .

run:
	./bin/gocomments .

generate-dataset:
	python ./dataset/generate_func_comments_from_local_repo.py

generate-model:
	sudo docker-compose -f ./model/docker-compose.yml down
	sudo docker-compose -f ./model/docker-compose.yml build # --no-cache
	sudo docker-compose -f ./model/docker-compose.yml up -d # --build
	sudo docker-compose -f ./model/docker-compose.yml logs tensorflow-container

test-model:
	pip install tf2onnx
	pip install onnxruntime
	chmod -R 755 runs/saved_model_*
	chmod -R 755 onnx
	go run ./model/prompt.go

install:
	sudo go install ./
