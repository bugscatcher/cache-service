all: build

build:
	CGO_ENABLED=0 go build -o cache_service

generate:
	protoc -I=proto --gogofaster_out=plugins=grpc:pb `ls proto`