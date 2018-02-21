BINARY = ./slack-rtm-http
VERSION = 1.0

run: build
	${BINARY} -v

build:
	go fmt main/*.go
	go build -o ${BINARY} main/*.go

release:
	./release.sh ${VERSION}

deps:
	glide update
