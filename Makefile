build:
	go fmt *.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o fGetterSrv
	cp dbInit.sql ./bin
	cp lazygit.sh ./bin/HG
	mv fGetter ./bin
