.PHONY: freeze

IMAGE := "111285186890.dkr.ecr.us-east-1.amazonaws.com/fastapi-btc-calc"
GIT_COMMIT_HASH := $(shell git rev-parse --short HEAD)
export GIT_COMMIT_HASH

build:
	docker build -t $(IMAGE):$(GIT_COMMIT_HASH) .
	docker push $(IMAGE):$(GIT_COMMIT_HASH)

freeze:
	pip freeze > requirements.txt

run:
	uvicorn main:app --reload

json_test:
	curl -X POST -H "Content-Type: application/json" -d @example_data.json http://localhost:8000

check:
	curl -X POST -H "Content-Type: application/json" -d @data.json https://btc-api.borisquiroz.dev
