BUILD_DIR=build
ENV_FILE=.env
BIN_NAME=go-radius
DOCKER_TAG=dcaponi1/${BIN_NAME}

build:
	mkdir -p ${BUILD_DIR} && \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./${BUILD_DIR}/${BIN_NAME}

build_pi:
	mkdir -p ${BUILD_DIR} && \
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./${BUILD_DIR}/${BIN_NAME}

clean:
	rm -rf ${BUILD_DIR}

run:
	docker run --env-file=${ENV_FILE} -p 1812:1812/udp ${BIN_NAME}

image: build
	docker build -t ${BIN_NAME} .

tag:
	docker tag ${DOCKER_TAG}:latest ${DOCKER_TAG}:latest

deploy: build tag
	docker push ${DOCKER_TAG}:latest

rebuild: clean build
