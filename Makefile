.PHONY: default build build-app run generate-dataset generate-model generate-from-checkpoint generate-api install

DOCKER_COMPOSE = sudo docker-compose -f
DOCKER_DOWN = down
DOCKER_BUILD = build
DOCKER_UP = up -d
DOCKER_LOGS = logs -f

default: build

build: build-app

build-app:
	go build -o bin/gocomments -v .

run: build-app
	./bin/gocomments .

generate-dataset:
	python ./dataset/generate_func_comments_from_local_repo.py

generate-model:
	$(DOCKER_COMPOSE) ./model/docker-compose.yml -p model $(DOCKER_DOWN)
	$(DOCKER_COMPOSE) ./model/docker-compose.yml -p model $(DOCKER_BUILD) # --no-cache
	$(DOCKER_COMPOSE) ./model/docker-compose.yml -p model $(DOCKER_UP) # --build
	$(DOCKER_COMPOSE) ./model/docker-compose.yml -p model $(DOCKER_LOGS) create

# generate-from-checkpoint:
#	$(DOCKER_COMPOSE) ./model-from-checkpoint/docker-compose.yml -p model-from-checkpoint $(DOCKER_DOWN)
#	$(DOCKER_COMPOSE) ./model-from-checkpoint/docker-compose.yml -p model-from-checkpoint $(DOCKER_BUILD) # --no-cache
#	$(DOCKER_COMPOSE) ./model-from-checkpoint/docker-compose.yml -p model-from-checkpoint $(DOCKER_UP) # --build
#	$(DOCKER_COMPOSE) ./model-from-checkpoint/docker-compose.yml -p model-from-checkpoint $(DOCKER_LOGS) create

generate-test:
	cd ./test-models && go mod tidy && go mod vendor
	$(DOCKER_COMPOSE) ./test-models/docker-compose.yml $(DOCKER_DOWN)
	$(DOCKER_COMPOSE) ./test-models/docker-compose.yml $(DOCKER_BUILD) # --no-cache
	$(DOCKER_COMPOSE) ./test-models/docker-compose.yml $(DOCKER_UP) # --build
	$(DOCKER_COMPOSE) ./test-models/docker-compose.yml $(DOCKER_LOGS) go-tf-app

generate-api:

# convert-model:
#	pip install tf2onnx
#	pip install onnxruntime
#	chmod -R 755 runs/saved_model_*
#	chmod -R 755 onnx

install:
	sudo go install ./