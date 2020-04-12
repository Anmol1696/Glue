include .env
export

CONTAINER_NAME=glue-container
TAG_NAME=glue
ENV_FILE=.env

TEST_COVERAGE_FILE=bin/coverage.out

.PHONY: go-build go-run test cover vet

all: go-build go-run


# All go commands
go-build:
	go build -a -o bin/app .

go-run:
	./bin/app --verbose

test:
	go test ./ -v coverprofile=${TEST_COVERAGE_FILE} -coverpkg=. -failfast

cover:
	go tool cover -func ${TEST_COVERAGE_FILE}

vet:
	go vet ./


# All Docker help commands
docker-build:
	docker build -t ${TAG_NAME} . --force-rm --no-cache

docker-run:
	docker run -it -p 5555:5555 --name ${CONTAINER_NAME} \
		--env-file ${ENV_FILENAME} \
		-d ${TAG_NAME}

docker-push: docker-build
	docker push ${TAG_NAME}

docker-stop:
	docker stop ${CONTAINER_NAME}

docker-start:
	docker start ${CONTAINER_NAME}

docker-bash:
	docker exec -it ${CONTAINER_NAME} /bin/sh

clear-container: docker-stop
	docker rm ${CONTAINER_NAME}

clear-image:
	docker rmi ${TAG_NAME}

clear: clear-container clear-image
