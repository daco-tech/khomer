all: deps build

deps:
	@go get github.com/labstack/echo/v4
	@go get github.com/common-nighthawk/go-figure
	@go get k8s.io/client-go@v0.25.0
	@go get k8s.io/client-go/discovery@v0.25.0
	@go get k8s.io/client-go/util/cert@v0.25.0
	@go get k8s.io/client-go/applyconfigurations/internal@v0.25.0
	@go get k8s.io/client-go/util/flowcontrol@v0.25.0
	@go get k8s.io/client-go/rest@v0.25.0

build:
	@go build

install: 
	@go install

run: 
	@go run server.go
