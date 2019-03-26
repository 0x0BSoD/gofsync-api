build:
	go fmt *.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o fGetter
	cp dbInit.sql ./bin
	cp lazygit.sh ./bin/HG
	mv fGetter ./bin
