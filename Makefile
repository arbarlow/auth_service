# vi: ft=make
.PHONY: proto test benchmark get ci docker

proto:
	protoc -I $$GOPATH/src/ -I . auth_service.proto --lile-server_out=. --go_out=plugins=grpc:$$GOPATH/src

test:
	SIGNING_KEY="sometestkey" go test -v ./...

benchmark:
	go test -bench=. -benchmem -benchtime 10s

get:
	go get -u -t ./...

ci: get test docker

docker:
	GOOS=linux GOARCH=amd64 go build -o build/auth_service ./auth_service
	docker build . -t lileio/auth_service:`git rev-parse --short HEAD`
