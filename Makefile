all: deps build

deps:
	@go get github.com/labstack/echo/v4
	@go get github.com/common-nighthawk/go-figure

build:
	@go build

install: 
	@go install

run: 
	@go run server.go
