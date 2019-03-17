build:
	go fmt *.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o fGetter
	mv fGetter ./bin
