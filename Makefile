build:
	go fmt *.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o fGetter
	cp dbInit.sql ./bin
	cp docker-compose.yaml ./bin/
	cp lazygit.sh ./bin/HG
	mv fGetter ./bin
