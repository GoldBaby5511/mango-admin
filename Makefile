PROJECT:=mango-admin

.PHONY: build
build:
	CGO_ENABLED=0 go build -o mango-admin main.go
build-linux:
	env GOOS=linux GOARCH=amd64 go build
build-sqlite:
	go build -tags sqlite3 -o mango-admin main.go
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t mango-admin:latest
