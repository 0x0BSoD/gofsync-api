cd ..amd64build:
	go fmt *.go
	GOOS=linux GOARCH=amd64 go build -race -ldflags="-w -s" -o gofsync
	mv gofsync ./dist

docker:
	sudo docker build --build-arg token=${TOKEN} -t goapi .