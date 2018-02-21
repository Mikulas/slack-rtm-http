BINARY = ./slack-rtm-http

run: build
	${BINARY} -v

build:
	go fmt main/*.go
	go build -o ${BINARY} main/*.go

release:


deps:
	glide update
