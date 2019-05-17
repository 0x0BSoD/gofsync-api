build:
	go fmt *.go
	GOOS=linux GOARCH=amd64 go build -race -ldflags="-w -s" -o fGetter
	mv fGetter ./dist
